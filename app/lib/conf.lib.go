package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	// "strconv"
	"strings"
	"sync"
	"time"

	"xi/app/cfg"

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
	Hook         Hook

	IntermediateMap  map[string]any
	IntermediateJson []byte

	watch *fsnotify.Watcher
	mu    sync.RWMutex
	once  sync.Once
}

var (
	Cfg = &ConfLib{
		koanfCli: koanf.New("."),
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
	// reJsonBoolStr = regexp.MustCompile(`:\s*"(?i:true|false|1|0)"`) // detects "true", "false", "1", "0" as string values
	reJsonVar = regexp.MustCompile(`\$\{([^}:]*)(:-([^}]*))?\}|\$\{\}`)
)

func init() {
	Cfg.Init()
	if err := Cfg.Daemon(); err != nil {
		log.Printf("⚠️ [Conf] Daemon WRN: setup failed: %v", err)
	}
}

func (c *ConfLib) Init(filePath ...string) { c.once.Do(func() { c.InitCore(filePath...) }) }

func (c *ConfLib) InitCore(filePath ...string) error {
	// Init Env and pre funcs
	Env.Init()
	// c.Hook.RunPre()

	// Assign Config Files
	c.FilesLoaded = nil
	c.Files = c.FilesDefault
	if len(filePath) > 0 {
		c.Files = filePath
	}

	// Overlay Resolved Config Files
	for _, path := range c.Files {
		raw, err := os.ReadFile(path)
		if err != nil {
			log.Printf("⚠️  [Conf] Init WRN: Skipped '%s': %v\n", path, err)
			continue
		}

		// Resolve ${ENV} | ${ENV:-fallback}
		resolvedJsonEnv := []byte(c.cleanJson(c.resolveJsonEnv(string(raw))))

		if err := c.koanfCli.Load(rawbytes.Provider(resolvedJsonEnv), koanfJson.Parser()); err != nil {
			log.Printf("⚠️  [Conf] Init WRN: Parsing failed '%s': %v\n", path, err)
			continue
		}

		c.FilesLoaded = append(c.FilesLoaded, path)
	}

	// fmt.Printf("-->\n%s\n%s\n", c.AllMap(), c.AllJson())
	// Preserve intermediate data, to be used in future features
	c.IntermediateMap = c.AllMap()
	// c.IntermediateJson = c.AllJson()

	c.ConfPostView()

	// Resolve internal refs ${ref} | ${ref:-fallback}
	c.resolveJsonVars()

	// c.Hook.RunPost()

	if len(c.FilesLoaded) > 0 {
		log.Printf("✅ [Conf] loaded: %s\n", c.FilesLoaded)
	} else {
		log.Printf("⚠️  [Conf] Init WRN: No config loaded.")
	}
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
func (c *ConfLib) resolveJsonVars() {
	var nested = c.AllMap()

	var resolveValue func(val any) any

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
							// If the string is exactly like "{{key}}", we replace whole field with object/array
							if str == "{{"+key+"}}" {
								v[k] = typed
							} else {
								// Otherwise just resolve it to string
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

	resolved := resolveValue(nested).(map[string]any)

	newData, _ := json.Marshal(resolved)
	k := koanf.New(".")
	_ = k.Load(rawbytes.Provider(newData), koanfJson.Parser())
	c.koanfCli = k
}

// sync json connfig with existing config
func (c *ConfLib) postSetup(jsonMap map[string]any) error {
	// Convert map[string]any to proper []byte(json) for further processing
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("⚠️  [Conf] PostSetup WRN: failed to marshal DataJson: %w", err)
	}

	// store
	if err := json.Unmarshal(jsonBytes, cfg.Get()); err != nil {
		return fmt.Errorf("⚠️  [Conf] PostSetup WRN: failed to unmarshal into Config struct: %w", err)
		// fmt.Printf("--> Default:\n%v\nMap:\n%v\nJson:\n%s\n", pageDefault, c.AllMapStruct(), c.AllJsonStruct())
	}

	if err := c.koanfCli.Load(rawbytes.Provider(jsonBytes), koanfJson.Parser()); err != nil {
		return fmt.Errorf("⚠️  [Conf] PostSetup WRN: Failed to load JSON config into Koanf: %w", err)
		// fmt.Printf("--> Default:\n%v\nMap:\n%v\nJson:\n%s\n", pageDefault, c.AllMapStruct(), c.AllJsonStruct())
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
					log.Printf("🔄 Config changed: %s (%s)", event.Name, event.Op)
					// Sleep briefly to avoid partial writes
					time.Sleep(100 * time.Millisecond)
					err := c.InitCore(c.Files...)
					if err != nil {
						log.Printf("⚠️ Config reload failed: %v", err)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("⚠️ Config daemon Err: %v", err)
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
