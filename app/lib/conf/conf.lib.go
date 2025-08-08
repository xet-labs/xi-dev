package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"xi/app/lib/cfg"
	"xi/app/lib/env"
	"xi/app/lib/hook"
	"xi/app/model"

	"github.com/fsnotify/fsnotify"
	koanfJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
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
	Conf = &ConfLib{
		FilesDefault: []string{
			// "app/data/config/config.json",
			"config/config.json",
		},
	}

	reJsonEnv         = regexp.MustCompile(`\$\{([A-Z0-9_]+)(:-([^}]*))?\}`)
	reJsonEnvPost     = regexp.MustCompile(`(?m)(,\s*)?__REMOVE__(,\s*)?|^__REMOVE__(,\s*)?`)
	reJsonDoubleQuote = regexp.MustCompile(`""([^"\n\r]+?)""`)
	reJsonIntCast     = regexp.MustCompile(`:\s*"(-?\d+)\.int"`)
	reJsonBoolStr     = regexp.MustCompile(`:\s*"(true|false|1|0)"`)
	reJsonVar         = regexp.MustCompile(`\$\{([^}:]*)(:-([^}]*))?\}|\$\{\}`)
)

func init() {
	Conf.Init()
	if err := Conf.Daemon(); err != nil {
		log.Printf("⚠️ [Conf] Daemon WRN: setup failed: %v", err)
	}
}

func (c *ConfLib) Init(filePath ...string) { c.once.Do(func() { c.InitCore(filePath...) }) }

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

	okStatus := "Loaded"
	errStatus := "No config loaded."
	if c.hasInitialized {
		okStatus = "Reloaded"
		errStatus = "Config did not reload."
	}

	// Create a new Koanf instance for atomic reload
	newKoanf := koanf.New(".")

	var newFilesLoaded []string

	// Load all config files into newKoanf
	for _, path := range c.Files {
		raw, err := os.ReadFile(path)
		if err != nil {
			log.Printf("⚠️  [Conf] Init WRN: Skipped '%s': %v", path, err)
			continue
		}

		resolved, err := c.resolveJsonVars(c.cleanJson(c.resolveJsonEnv(string(raw))))
		if err != nil {
			log.Printf("⚠️  [Conf] Init WRN: Resolving failed '%s': %v", path, err)
			continue
		}
		if err := newKoanf.Load(rawbytes.Provider([]byte(resolved)), koanfJson.Parser()); err != nil {
			log.Printf("⚠️  [Conf] Init WRN: Parsing failed '%s': %v", path, err)
			continue
		}
		newFilesLoaded = append(newFilesLoaded, path)
	}

	// Fail and rollback if nothing loaded
	if len(newFilesLoaded) == 0 {
		log.Printf("⚠️  [Conf] WRN: %s", errStatus)
		if c.hasInitialized {
			// restore old state
			c.koanfCli = oldKoanfCli
			c.FilesLoaded = oldFilesLoaded
			c.IntermediateMap = oldIntermediateMap
			c.IntermediateJson = oldIntermediateJson
		}
		return fmt.Errorf("no config files loaded successfully")
	}

	// Temporarily assign newKoanf to c so we can use c methods
	c.koanfCli = newKoanf
	c.FilesLoaded = newFilesLoaded

	// Preserve intermediate data
	c.IntermediateMap = c.AllMap()
	c.IntermediateJson = c.AllJson()

	c.ConfPostView()

	log.Printf("✅ [Conf] %s: %s", okStatus, newFilesLoaded)
	c.hasInitialized = true

	return nil
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
func (c *ConfLib) postSetup(jsonMap map[string]any) error {
	// Convert map[string]any to proper []byte(json) for further processing
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Printf("⚠️  [Conf] PostSetup WRN: failed to marshal DataJson: %v", err)
		return err
	}

	// store
	rawCfg := model.Config{}
	if err := json.Unmarshal(jsonBytes, &rawCfg); err != nil {
		log.Printf("⚠️  [Conf] PostSetup WRN: failed to unmarshal into Config struct: %v", err)
		return err
	}
	cfg.Update(rawCfg)

	if err := c.koanfCli.Load(rawbytes.Provider(jsonBytes), koanfJson.Parser()); err != nil {
		log.Printf("⚠️  [Conf] PostSetup WRN: Failed to load JSON config into Koanf: %v", err)
		return err
	}
	return nil
}

func (c *ConfLib) Daemon() error {
	if c.watch != nil {
		return nil // already watching
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("⚠️  [Conf] Daemon WRN: failed to create: %v", err)
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
					log.Printf("🔄 [Conf] Changed: %s (%s)", event.Name, event.Op)

					if err := c.InitCore(); err != nil {
						log.Printf("⚠️  [Conf] Reload failed: %v", err)
					}
					// Sleep briefly to avoid partial writes
					time.Sleep(100 * time.Millisecond)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("⚠️  [Conf] daemon Err: %v", err)
			}
		}
	}()

	for _, path := range c.Files {
		// Ensure file exists before watching (else no event will be triggered)
		if _, err := os.Stat(path); err == nil {
			if err := watcher.Add(path); err != nil {
				log.Printf("⚠️  [Conf] Daemon WRN: failed to watch %s: %v", path, err)
			}
		} else {
			log.Printf("⚠️  [Conf] Daemon WRN: missing file: %s", path)
		}
	}

	return nil
}
