package lib

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
	"xi/app/schema"

	"github.com/fsnotify/fsnotify"
	koanfJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

type ConfLib struct {
	Global       schema.Config
	Files        []string
	FilesLoaded  []string
	FilesDefault []string
	koanfCli     *koanf.Koanf
	Hooks        Hook

	IntermediateMap  map[string]any
	IntermediateJson map[string]any

	DataMap  map[string]any
	DataJson map[string]any

	watch *fsnotify.Watcher
	mu    sync.RWMutex
	once  sync.Once
}

var (
	Cfg = &ConfLib{
		koanfCli: koanf.New("."),
		FilesDefault: []string{
			"app/data/config/config.json",
			"config/config.json",
		},
	}
	reJsonEnv     = regexp.MustCompile(`\$\{([A-Z0-9_]+)(:-([^}]*))?\}`)
	reJsonEnvPost = regexp.MustCompile(`(?m)(,\s*)?"__REMOVE__"(,\s*)?|^"__REMOVE__"(,\s*)?`)
	reJsonVar     = regexp.MustCompile(`\$\{([^}:]*)(:-([^}]*))?\}|\$\{\}`)
)

func init() {
	Cfg.Init()
	if err := Cfg.Daemon(); err != nil {
		log.Printf("⚠️ Config Daemon setup failed: %v", err)
	}

	// log.Printf("--> %s", Cfg.Get("app.domain"))
	// log.Printf("--> %s", Cfg.GetMap("app"))
	// log.Printf("--> %s", Cfg.koanfCli.Get(""))
}

func (c *ConfLib) Init(filePath ...string) { c.once.Do(func() { c.InitCore(filePath...) }) }

func (c *ConfLib) InitCore(filePath ...string) error {
	// Init Env
	Env.Init()

	c.Hooks.RunPre()

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
			log.Printf("⚠️  Config Skipped %s: %v\n", path, err)
			continue
		}

		// Resolve ${ENV} | ${ENV:-fallback}
		resolvedJsonEnv := []byte(c.cleanJson(c.resolveJsonEnv(string(raw))))

		if err := c.koanfCli.Load(rawbytes.Provider(resolvedJsonEnv), koanfJson.Parser()); err != nil {
			log.Printf("⚠️  Config parsing failed %s: %v\n", path, err)
			continue
		}

		c.FilesLoaded = append(c.FilesLoaded, path)
	}

	// Resolve internal refs ${ref} | ${ref:-fallback}
	c.resolveJsonVars()

	c.Hooks.RunPost()

	if len(c.FilesLoaded) > 0 {
		log.Printf("✅ Config loaded: %s\n", c.FilesLoaded)
	} else {
		log.Println("⚠️  No config loaded.")
	}
	return nil
}

// resolveJsonEnv replaces ${ENV} or ${ENV:-fallback} with actual values
func (c *ConfLib) resolveJsonEnv(input string) string {
	return reJsonEnv.ReplaceAllStringFunc(input, func(match string) string {
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
}

func (c *ConfLib) cleanJson(input string) string {
	return reJsonEnvPost.ReplaceAllString(input, "")
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

func (c *ConfLib) Daemon() error {
	if c.watch != nil {
		return nil // already watching
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("⚠️ Config Err: failed to create daemon: %v", err)
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
				log.Printf("⚠️ Config daemon failed to watch %s: %v", path, err)
			}
		} else {
			log.Printf("⚠️ Config missing file: %s", path)
		}
	}

	return nil
}
