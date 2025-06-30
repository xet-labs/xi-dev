package lib

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Central utility
type EnvLib struct {
	envMap   map[string]interface{}
	rw       sync.RWMutex
	once     sync.Once
}

// Global instance
var Env = &EnvLib{
	envMap: make(map[string]interface{}),
}

func init() {
	Env.Init()
}

// InitForce forcibly reloads environment variables from .env and system.
func (e *EnvLib) InitForce() {
	e.rw.Lock()
	defer e.rw.Unlock()

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found or couldn't be loaded")
	} else {
		log.Println("✅ Env initialized..")
	}

	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			e.envMap[parts[0]] = parts[1]
		}
	}
}

// Init ensures one-time lazy initialization.
func (e *EnvLib) Init() {
	e.once.Do(e.InitForce)
}

// Env returns string value from env with optional fallback.
func (e *EnvLib) Get(key string, fallback ...string) string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.envMap[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// Get retrieves a value (any type).
func (e *EnvLib) Raw(key string, fallback ...interface{}) interface{} {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.envMap[key]; ok {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return nil
}

// Set sets a key to a value.
func (e *EnvLib) Set(key string, value interface{}) {
	e.rw.Lock()
	defer e.rw.Unlock()
	e.envMap[key] = value
}

// Unset deletes a key.
func (e *EnvLib) Unset(key string) {
	e.rw.Lock()
	defer e.rw.Unlock()
	delete(e.envMap, key)
}

// Bool returns a boolean value from env with optional fallback.
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

// Int returns an int value from env or fallback.
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

// All returns a copy of all stored environment variables.
func (e *EnvLib) All() map[string]interface{} {
	e.rw.RLock()
	defer e.rw.RUnlock()

	copyMap := make(map[string]interface{}, len(e.envMap))
	for k, v := range e.envMap {
		copyMap[k] = v
	}
	return copyMap
}

// As allows custom typed parsing of string env vars (Go 1.18+).
func As[T any](env *EnvLib, key string, fallback T, parser func(string) (T, error)) T {
	val, ok := env.Raw(key).(string)
	if !ok {
		return fallback
	}
	if parsed, err := parser(val); err == nil {
		return parsed
	}
	return fallback
}
