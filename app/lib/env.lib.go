// Optimized EnvLib for performance, safety, and clarity
package lib

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// EnvLib provides a thread-safe environment loader with caching
type EnvLib struct {
	vars map[string]interface{}
	once sync.Once
	rw   sync.RWMutex
}

// Global singleton instance
var Env = &EnvLib{
	vars: make(map[string]interface{}),
}

func init() {
	Env.Init()
}

// InitForce reloads .env and OS environment variables forcibly
func (e *EnvLib) InitForce() {
	e.rw.Lock()
	defer e.rw.Unlock()

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not loaded")
	} else {
		log.Println("✅ Env loaded")
	}

	for _, kv := range os.Environ() {
		if parts := strings.SplitN(kv, "=", 2); len(parts) == 2 {
			e.vars[parts[0]] = parts[1]
		}
	}
}

// Init ensures single-time environment initialization
func (e *EnvLib) Init() {
	e.once.Do(e.InitForce)
}

// Get returns string value for key or fallback
func (e *EnvLib) Get(key string, fallback ...string) string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.vars[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// Raw returns value (any type) or fallback
func (e *EnvLib) Raw(key string, fallback ...interface{}) interface{} {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.vars[key]; ok {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return nil
}

// Set adds or updates a key
func (e *EnvLib) Set(key string, value interface{}) {
	e.rw.Lock()
	defer e.rw.Unlock()
	e.vars[key] = value
}

// Unset deletes a key
func (e *EnvLib) Unset(key string) {
	e.rw.Lock()
	defer e.rw.Unlock()
	delete(e.vars, key)
}

// Bool returns boolean value or fallback
func (e *EnvLib) Bool(key string, fallback ...bool) bool {
	switch v := e.Raw(key).(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	case int:
		return v != 0
	case float64:
		return v != 0.0
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return false
}

// Int returns int value or fallback
func (e *EnvLib) Int(key string, fallback int) int {
	switch v := e.Raw(key).(type) {
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// All returns a copy of all environment variables
func (e *EnvLib) All() map[string]interface{} {
	e.rw.RLock()
	defer e.rw.RUnlock()

	snapshot := make(map[string]interface{}, len(e.vars))
	for k, v := range e.vars {
		snapshot[k] = v
	}
	return snapshot
}

// As parses and returns custom type
func As[T any](env *EnvLib, key string, fallback T, parser func(string) (T, error)) T {
	if val, ok := env.Raw(key).(string); ok {
		if parsed, err := parser(val); err == nil {
			return parsed
		}
	}
	return fallback
}
