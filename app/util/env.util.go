package util

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envVar         = make(map[string]interface{})
	envInitialized = false
	envLock        sync.RWMutex
)

// InitEnv loads the .env file and populates the envVar map.
func InitEnv() {
	envLock.Lock()
	defer envLock.Unlock()

	if envInitialized {
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found or couldn't be loaded")
	} else {
		log.Println("✅ Env loaded..")
	}

	// Populate envVar from os.Environ
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			envVar[parts[0]] = parts[1]
		}
	}

	envInitialized = true
}

// Env returns a string value from the environment or fallback.
func Env(key string, fallback ...string) string {
	InitEnv()

	envLock.RLock()
	defer envLock.RUnlock()

	if val, ok := envVar[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// EnvGet returns a value from the envVar map as interface{}.
func EnvGet(key string) interface{} {
	InitEnv()

	envLock.RLock()
	defer envLock.RUnlock()
	return envVar[key]
}

// EnvSet sets a value into the envVar map.
func EnvSet(key string, value interface{}) {
	InitEnv()

	envLock.Lock()
	defer envLock.Unlock()
	envVar[key] = value
}

// EnvUnset removes a key from envVar.
func EnvUnset(key string) {
	InitEnv()

	envLock.Lock()
	defer envLock.Unlock()
	delete(envVar, key)
}

// EnvBool returns a boolean environment variable.
func EnvBool(key string, fallback ...bool) bool {
	val := EnvGet(key)

	switch v := val.(type) {
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

// EnvInt returns an integer environment variable.
func EnvInt(key string, fallback int) int {
	val := EnvGet(key)

	switch v := val.(type) {
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// EnvAll returns a copy of the environment variables.
func EnvAll() map[string]interface{} {
	InitEnv()

	envLock.RLock()
	defer envLock.RUnlock()

	copyMap := make(map[string]interface{}, len(envVar))
	for k, v := range envVar {
		copyMap[k] = v
	}
	return copyMap
}
