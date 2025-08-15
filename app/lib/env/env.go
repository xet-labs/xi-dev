package env

import (
	"fmt"
	"os"
	"strconv"

	"sync"

	"xi/app/lib/logger"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvLib struct {
	app map[string]any // store Runtime env
	sys map[string]any // store Systems env

	once sync.Once
	mu   sync.RWMutex
}

var Env = &EnvLib{
	sys: make(map[string]any),
}

// Init, ensure single-time env initialization
func (e *EnvLib) Init() { e.once.Do(e.InitCore) }

// InitCore, reload .env and OS env variables forcibly
func (e *EnvLib) InitCore() {
	logger.Logger.Init()
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := os.Stat(".env"); err == nil || !os.IsNotExist(err) {
		if err := godotenv.Load(".env"); err == nil {
			log.Info().Str("file", ".env").Msg("Env loaded")
		}
	}

	// for _, kv := range os.Environ() {
	// 	if parts := strings.SplitN(kv, "=", 2); len(parts) == 2 {
	// 		e.sys[parts[0]] = parts[1]
	// 	}
	// }
}

func (e *EnvLib) Load(file string) error {
	if file != "" {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				return fmt.Errorf("godot: %w", err)
			}
		} else if os.IsNotExist(err) {
			return err
		} else {
			return err
		}
	}
	return nil

}

// Get returns string value for key or fallback
func (e *EnvLib) Get(key string, fallback ...string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()

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
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sys[key] = value
}

// Unset deletes a key
func (e *EnvLib) Unset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.sys, key)
}

// Raw returns value (any type) or fallback
func (e *EnvLib) Raw(key string, fallback ...any) any {
	e.mu.RLock()
	defer e.mu.RUnlock()

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
	e.mu.RLock()
	defer e.mu.RUnlock()

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
