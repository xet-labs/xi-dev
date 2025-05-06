// utils/env.go
package utils

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"xi/app/global"
)

// Env returns the value of an environment variable, or an optional default value if not set.
func Env(key string, fallback ...string) string {
	if !global.EnvInitialized {

		if err := godotenv.Load(); err == nil {
			log.Println("Init env...")
		}
		global.EnvInitialized = true
	}

	if val := os.Getenv(key); val != "" {
		return val
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return ""
}
