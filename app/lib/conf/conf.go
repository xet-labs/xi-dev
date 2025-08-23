package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
	"xi/app/lib/cfg"
	"xi/app/lib/env"
	"xi/app/lib/hook"
	"xi/app/lib/util"
	"xi/app/model"

	"github.com/fsnotify/fsnotify"
	koanfJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

type ConfLib struct {
	Files        []string
	FilesLoaded  []string
	FilesDefault []string
	koanfCli     *koanf.Koanf
	Hook         hook.Hook

	IntermediateMap  map[string]any
	IntermediateJson []byte

	hasInitialized bool
	watch          *fsnotify.Watcher
	mu             sync.RWMutex
	once           sync.Once
}

var (
	Conf = &ConfLib{}
	reJsonEnv         = regexp.MustCompile(`\$\{([A-Z0-9_]+)(:-([^}]*))?\}`)
	reJsonEnvPost     = regexp.MustCompile(`(?m)(,\s*)?__REMOVE__(,\s*)?|^__REMOVE__(,\s*)?`)
	reJsonDoubleQuote = regexp.MustCompile(`""([^"\n\r]+?)""`)
	reJsonIntCast     = regexp.MustCompile(`:\s*"(-?\d+)\.int"`)
	reJsonBoolStr     = regexp.MustCompile(`:\s*"(true|false|1|0)"`)
	reJsonVar         = regexp.MustCompile(`\$\{([^}:]*)(:-([^}]*))?\}|\$\{\}`)
)

func init(){
	confFiles, err := util.File.GetWithExt(".json", "config", "app/data/config")
	if err != nil {
		log.Fatal().Err(err).Msg("Config")
	}
	Conf.FilesDefault = util.Str.UniqueSort(append([]string{
		"app/data/config/config.json",
	}, confFiles...))

}

func (c *ConfLib) Init(filePath ...string) {
	c.once.Do(func() {
		c.InitCore(filePath...)
		if err := Conf.Daemon(); err != nil {
			log.Warn().Msgf("Config Daemon: setup failed: %v", err)
		}
	})
}

func (c *ConfLib) InitCore(filePath ...string) error {
	env.Env.Init()
	c.mu.Lock()
	defer c.mu.Unlock()

	// Only clone defaults on first run if Files not set
	if !c.hasInitialized && c.Files == nil {
		c.Files = slices.Clone(c.FilesDefault)
	}
	if len(filePath) > 0 {
		c.Files = filePath
	}

	// Backup old states so if smthng goes wrong, restoration homecoming !!
	oldKoanfCli := c.koanfCli
	oldFilesLoaded := c.FilesLoaded
	oldIntermediateMap := c.IntermediateMap
	oldIntermediateJson := c.IntermediateJson

	okStatus := "loaded"
	errStatus := "No config loaded."
	if c.hasInitialized {
		okStatus = "reloaded"
		errStatus = "Config did not reload."
	}

	// Create a new Koanf instance for atomic reload
	newKoanf := koanf.New(".")

	var newFilesLoaded []string
	// Load all config files into newKoanf
	for i, path := range c.Files {

		noConfigKillSwitch := func() {
			// On initial run if no config file has loaded and this is the last config file with err, exit
			if !c.hasInitialized && len(newFilesLoaded) == 0 && i == len(c.Files)-1 {
				log.Fatal().Str("files", util.Str.QuoteSlice(c.Files)).
					Msg("Startup aborted: no valid configuration could be loaded from any source")
			}
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			log.Warn().Err(err).Str("file", path).
				Msg("Config Skipped: unable to read file")
			noConfigKillSwitch()
			continue
		}

		// Preproces and load env from the config
		// the env needs to be loaded and reprocess the config s it might be using ars from the env
		if _, tmp, err := c.preProcess(raw); err == nil {
			for _, file := range tmp.App.EnvFiles {
				if err := env.Env.Load(file); err != nil {
					log.Fatal().Err(err).Str("Env", file).Str("file", path).
						Msg("Config Preprocess failed to load Env")
					continue
				}
			}
		} else {
			log.Error().Err(err).Str("file", path).Msg("Config Preprocess failed for env")
			continue
		}

		// Fully preprocess data
		var resolved []byte
		if resolved, _, err = c.preProcess(raw); err != nil {
			log.Error().Err(err).Str("file", path).Msg("Config Preprocess failed")
		}
		// Sync/merge Json data
		if err := newKoanf.Load(rawbytes.Provider(resolved), koanfJson.Parser()); err != nil {
			log.Error().Str("file", path).Err(err).
				Msg("Config load failed: JSON is valid but merging into runtime configuration failed")
			noConfigKillSwitch()
			continue
		}
		newFilesLoaded = append(newFilesLoaded, path)
	}

	// Fail and rollback if nothing loaded
	if len(newFilesLoaded) == 0 {
		log.Warn().Msgf("Config %s", errStatus)
		if c.hasInitialized {
			// restore old state
			c.koanfCli = oldKoanfCli
			c.FilesLoaded = oldFilesLoaded
			c.IntermediateMap = oldIntermediateMap
			c.IntermediateJson = oldIntermediateJson
		}
		return fmt.Errorf("no config files could be loaded")
	}

	// Temporarily assign newKoanf to c so we can use c methods
	c.koanfCli = newKoanf
	c.FilesLoaded = newFilesLoaded

	// Preserve intermediate data
	// c.IntermediateMap = c.AllMap()
	// c.IntermediateJson = c.AllJson()

	c.ConfPostView()

	log.Info().Str("files", util.Str.QuoteSlice(newFilesLoaded)).Msgf("Config %s", okStatus)
	c.hasInitialized = true

	return nil
}

func (c *ConfLib) preProcess(rawJson []byte) ([]byte, model.Config, error) {
	// resolve Json varsand cleanups
	resolved, err := c.resolveJsonVars(c.cleanJson(c.resolveJsonEnv(string(rawJson))))
	if err != nil {
		log.Error().Err(fmt.Errorf("unable to resolve environment variables or sanitize JSON")).Msg("Config preprocess failed")
		return nil, model.Config{}, err
	}

	// Validate against config model
	structured := model.Config{}
	if err := json.Unmarshal([]byte(resolved), &structured); err != nil {
		log.Error().Err(fmt.Errorf("json does not match the expected Config structure")).
			Msg("Config preprocess failed, format invalid")
		return nil, model.Config{}, err
	}

	return []byte(resolved), structured, nil
}

// resolveJsonEnv replaces ${ENV} or ${ENV:-fallback} with actual values
func (c *ConfLib) resolveJsonEnv(input string) string {
	out := reJsonEnv.ReplaceAllStringFunc(input, func(match string) string {
		sub := reJsonEnv.FindStringSubmatch(match)
		key, def := sub[1], sub[3] // ENV, fallback

		if val, ok := os.LookupEnv(key); ok {
			return val
		}

		// Fallback
		if def != "" {
			return def
		}

		// No value, no fallback
		return "__REMOVE__"
	})
	// fmt.Printf("\n%s\n", out)
	return out

}

func (c *ConfLib) cleanJson(input string) string {
	// Remove "__REMOVE__"
	out := reJsonEnvPost.ReplaceAllString(input, "")

	// Optionally: fix trailing commas, multiple newlines, etc.
	out = strings.ReplaceAll(out, ",\n}", "\n}")
	out = strings.ReplaceAll(out, ",\n]", "\n]")

	// Fix ""value"" to "value", but skip empty ""
	out = reJsonDoubleQuote.ReplaceAllString(out, `"$1"`)

	out = reJsonIntCast.ReplaceAllString(out, ": $1")

	// Replaces string "true"/"false" -> true/false
	out = reJsonBoolStr.ReplaceAllStringFunc(out, func(match string) string {
		// Extract actual boolean value from the match using the submatch
		submatches := reJsonBoolStr.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		switch strings.ToLower(submatches[1]) {
		case "true":
			return ": true"
		case "false":
			return ": false"
		}
		return match // fallback
	})

	// fmt.Printf("\n%s\n", out)
	return out
}

// resolveJsonVars walks entire koanf data and resolves {{key.path}} expressions
func (c *ConfLib) resolveJsonVars(input string) (string, error) {
	var data any
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}

	var resolveValue func(any) any

	resolveValue = func(val any) any {
		switch v := val.(type) {
		case string:
			return reJsonVar.ReplaceAllStringFunc(v, func(match string) string {
				sub := reJsonVar.FindStringSubmatch(match)
				key := sub[1]
				def := sub[3]
				if val := c.koanfCli.String(key); val != "" {
					return val
				}
				return def
			})
		case map[string]any:
			for k, vv := range v {
				if str, ok := vv.(string); ok && reJsonVar.MatchString(str) {
					sub := reJsonVar.FindStringSubmatch(str)
					key := sub[1]
					def := sub[3]

					val := c.koanfCli.Get(key)
					if val != nil {
						switch typed := val.(type) {
						case map[string]any, []any:
							// If exactly like "${key}", replace whole field with object/array
							if str == "${"+key+"}" {
								v[k] = typed
							} else {
								v[k] = str
							}
						case string:
							if typed != "" {
								v[k] = typed
							} else {
								v[k] = def
							}
						default:
							v[k] = typed
						}
					} else {
						v[k] = def
					}
				} else {
					v[k] = resolveValue(vv)
				}
			}
			return v

		case []any:
			for i, vv := range v {
				v[i] = resolveValue(vv)
			}
			return v

		default:
			return v
		}
	}

	resolved := resolveValue(data)

	outBytes, err := json.Marshal(resolved)
	if err != nil {
		return "", err
	}
	return string(outBytes), nil
}

// sync json connfig with existing config
func (c *ConfLib) postProcess(jsonMap map[string]any) error {
	// Convert map[string]any to proper []byte(json) for further processing
	jsonConf, err := json.Marshal(jsonMap)
	if err != nil {
		log.Warn().Msgf("Config PostSetup: failed to marshal user Config: %v", err)
		return err
	}

	// Merge config
	if err := c.koanfCli.Load(rawbytes.Provider(jsonConf), koanfJson.Parser()); err != nil {
		log.Warn().Msgf("Config Post-Setup: Failed to load JSON Config into Koanf: %v", err)
		return err
	}

	// Store Config to global 'cfg'
	rawCfg := model.Config{}
	if err := json.Unmarshal(c.AllJson(), &rawCfg); err != nil {
		log.Warn().Msgf("Config Post-Setup: Failed to unmarshal into Config struct: %v", err)
		return err
	}
	cfg.Update(rawCfg)

	return nil
}

// Config Daemon to reload config file changes
func (c *ConfLib) Daemon() error {
	if c.watch != nil {
		return nil // already watching
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Warn().Msgf("Config Daemon failed to launch: %v", err)
		return err
	}
	c.watch = watcher

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					log.Info().Str("event", event.Op.String()).Str("file", event.Name).Msg("Config changed")

					if err := c.InitCore(); err != nil {
						log.Warn().Err(err).Msg("Config reload failed")
					}
					// Sleep briefly to avoid partial writes
					time.Sleep(100 * time.Millisecond)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Warn().Msgf("Config Daemon: %v", err)
			}
		}
	}()

	for _, path := range c.Files {
		// Ensure file exists before watching (else no event will be triggered)
		if _, err := os.Stat(path); err == nil {
			if err := watcher.Add(path); err != nil {
				log.Warn().Msgf("Config Daemon failed to watch %s: %v", path, err)
			}
		} else {
			log.Warn().Msgf("Config Daemon missing file: %s", path)
		}
	}

	return nil
}
