package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf/v2"
	koanfJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
)

type ConfLib struct {
	Files         []string
	FilesLoaded   []string
	FilesDefault  []string
	koanfInstance *koanf.Koanf

	watch *fsnotify.Watcher
	rw   sync.RWMutex
	once sync.Once
}

var Cfg = &ConfLib{
	koanfInstance: koanf.New("."),
	FilesDefault: []string{
		"app/data/config/config.json",
		"config/config.json",
	},
}

func init() { Cfg.Init() }

func (c *ConfLib) Init(filePath ...string) {
	c.once.Do(func() {
		c.InitCore(filePath...)
	})
}

func (c *ConfLib) InitCore(filePath ...string) error {
	Env.Init()
	c.Files = c.FilesDefault
	if len(filePath) > 0 {
		c.Files = filePath
	}

	for _, path := range c.Files {
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Skipping config %s: %v\n", path, err)
			continue
		}

		// Resolve env vars like {{VAR:-fallback}} or {{VAR}}
		expanded := []byte(c.expandEnvVars(string(raw)))

		err = c.koanfInstance.Load(rawbytes.Provider(expanded), koanfJson.Parser())
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Failed parsing %s: %v\n", path, err)
			continue
		}

		c.FilesLoaded = append(c.FilesLoaded, path)
	}

	if len(c.FilesLoaded) > 0 {
		log.Printf("✅ Config loaded: %s\n", c.FilesLoaded)
		return nil
	}

	log.Println("⚠️  No config loaded.")
	return nil
}

var envPattern = regexp.MustCompile(`\{\{([A-Z0-9_]+)(:-([^}]*))?\}\}`)


// expandEnvVars replaces {{ENV}} or {{ENV:-fallback}} with actual values
func (c *ConfLib) expandEnvVars(s string) string {
	return envPattern.ReplaceAllStringFunc(s, func(match string) string {
		sub := envPattern.FindStringSubmatch(match)
		key := sub[1]
		def := sub[3]
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return def
	})
}

func (c *ConfLib) Get(path string) any {
	return c.koanfInstance.Get(path)
}

func (c *ConfLib) GetMap(path string) map[string]any {
	if val, ok := c.koanfInstance.Get(path).(map[string]any); ok {
		return val
	}
	return map[string]any{}
}

func (c *ConfLib) GetArray(path string) []any {
	if val, ok := c.koanfInstance.Get(path).([]any); ok {
		return val
	}
	return []any{}
}

func (c *ConfLib) JSON() string {
	out, err := json.MarshalIndent(c.All(), "", "  ")
	if err != nil {
		return "{}"
	}
	return string(out)
}

func (c *ConfLib) All() map[string]any {
	return c.koanfInstance.All()
}
