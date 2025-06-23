package util

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envVar        = make(map[string]interface{})
	envLock       = sync.RWMutex{}
	envInitialized = false
)

// splitEnv safely splits "KEY=VALUE" into [key, value]
func splitEnv(kv string) []string {
	for i := 0; i < len(kv); i++ {
		if kv[i] == '=' {
			return []string{kv[:i], kv[i+1:]}
		}
	}
	return []string{kv}
}

// InitEnv loads .env file and populates envVar map
func InitEnv() {
	envLock.Lock()
	defer envLock.Unlock()

	if envInitialized {
		return
	}

	// Load .env from project root
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found or couldn't be loaded")
	} else {
		log.Println("✅ Init env..")
	}

	// Read all env keys/values into map
	for _, kv := range os.Environ() {
		parts := splitEnv(kv)
		if len(parts) == 2 {
			envVar[parts[0]] = parts[1]
		}
	}

	envInitialized = true
}

// Env returns the value of an env variable, with optional fallback
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

// EnvGet returns value as interface{} (or nil if not found)
func EnvGet(key string) interface{} {
	InitEnv()

	envLock.RLock()
	defer envLock.RUnlock()
	return envVar[key]
}

// EnvSet sets or updates an environment variable in the map
func EnvSet(key string, value interface{}) {
	InitEnv()

	envLock.Lock()
	defer envLock.Unlock()
	envVar[key] = value
	os.Setenv(key, ToStr(value))
}

// EnvUnset removes the key from the map and unsets it from real env
func EnvUnset(key string) {
	InitEnv()

	envLock.Lock()
	defer envLock.Unlock()
	delete(envVar, key)
	os.Unsetenv(key)
}

// EnvAll returns a copy of the full envVar map
func EnvAll() map[string]interface{} {
	InitEnv()

	envLock.RLock()
	defer envLock.RUnlock()

	copied := make(map[string]interface{}, len(envVar))
	for k, v := range envVar {
		copied[k] = v
	}
	return copied
}


