package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type ConfLib struct {
	ConfigFiles []string
	data map[string]any
	rw sync.RWMutex
}

// Singleton instance
var Cfg = &ConfLib{
	ConfigFiles: []string{"app/config.json", "conf/config.json"},
}

// Regex to parse {{VAR:-default}} format
var envPattern = regexp.MustCompile(`\{\{\s*([A-Z0-9_]+)(:-([^}]*))?\s*\}\}`)

// Load loads and parses the config JSON file
func (c *ConfLib) Load(filePath ...string) error {
	configFiles := c.ConfigFiles
	if len(filePath) > 0 { configFiles = filePath }

	for _, path := range configFiles {
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Skipping config %s: %v\n", path, err)
			continue
		}

		expanded := c.ExpandJsonVar(string(raw))

		var parsed map[string]any
		if err := json.Unmarshal([]byte(expanded), &parsed); err != nil {
			return fmt.Errorf("json decode (%s): %w", path, err)
		}

		c.rw.Lock()
		c.data = deepMerge(c.data, parsed)
		c.rw.Unlock()

		fmt.Printf("✅ Loaded config: %s\n", path)
	}

	return nil
}

// ExpandJsonVar resolves {{VAR:-fallback}} in raw JSON
func (c *ConfLib) ExpandJsonVar(input string) string {
	return envPattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := envPattern.FindStringSubmatch(match)
		key := parts[1]
		fallback := parts[3]
		val, ok := os.LookupEnv(key)
		if ok {
			return val
		}
		return fallback
	})
}

// Get retrieves a value using dot-notation path
func (c *ConfLib) Get(path string) any {
	c.rw.RLock()
	defer c.rw.RUnlock()

	parts := strings.Split(path, ".")
	var cur any = c.data

	for _, part := range parts {
		switch typed := cur.(type) {
		case map[string]any:
			cur = typed[part]
		default:
			return nil
		}
	}
	return cur
}

// GetMap returns a map for a path (e.g., object)
func (c *ConfLib) GetMap(path string) map[string]any {
	val := c.Get(path)
	if m, ok := val.(map[string]any); ok {
		return m
	}
	return nil
}

// GetArray returns array/slice at path
func (c *ConfLib) GetArray(path string) []any {
	val := c.Get(path)
	if arr, ok := val.([]any); ok {
		return arr
	}
	return nil
}

func (c *ConfLib) JSON() string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	out, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(out)
}


func (c *ConfLib) All() map[string]any {
	c.rw.RLock()
	defer c.rw.RUnlock()

	// Optionally return a deep copy if you want to avoid mutation from outside
	return c.data
}