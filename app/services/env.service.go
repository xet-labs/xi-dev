// services/env
package services

import (
	"log"

	"github.com/joho/godotenv"
	"xi/app/global"
)

func InitEnv() {
	if global.EnvInitialized {
		return
	}

	if err := godotenv.Load(); err == nil {
		log.Println("Init env...")
	} else {
		log.Println("Err loading env")
	}

	global.EnvInitialized = true
}