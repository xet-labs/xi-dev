package lib

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type EnvLib struct {
	app  map[string]any // store Runtime env
	sys  map[string]any // store Systems env
	once sync.Once
	rw   sync.RWMutex
}

var Env = &EnvLib{
	sys: make(map[string]any),
}

func init() { Env.Init() }

// Init, ensure single-time env initialization
func (e *EnvLib) Init() { e.once.Do(e.InitForce) }

// InitForce, reload .env and OS env variables forcibly
func (e *EnvLib) InitForce() {
	e.rw.Lock()
	defer e.rw.Unlock()

	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  Env couldnt be loaded: %v", err)
	} else {
		log.Println("✅ Env loaded..")
	}

	for _, kv := range os.Environ() {
		if parts := strings.SplitN(kv, "=", 2); len(parts) == 2 {
			e.sys[parts[0]] = parts[1]
		}
	}
}

// Get returns string value for key or fallback
func (e *EnvLib) Get(key string, fallback ...string) string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.sys[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// Set adds or updates a key
func (e *EnvLib) Set(key string, value any) {
	e.rw.Lock()
	defer e.rw.Unlock()
	e.sys[key] = value
}

// Unset deletes a key
func (e *EnvLib) Unset(key string) {
	e.rw.Lock()
	defer e.rw.Unlock()
	delete(e.sys, key)
}

// Raw returns value (any type) or fallback
func (e *EnvLib) Raw(key string, fallback ...any) any {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if val, ok := e.sys[key]; ok {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return nil
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
func (e *EnvLib) All() map[string]any {
	e.rw.RLock()
	defer e.rw.RUnlock()

	snapshot := make(map[string]any, len(e.sys))
	for k, v := range e.sys {
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
