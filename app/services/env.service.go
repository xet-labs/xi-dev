// services/env
package services

import (
	"log"

	"xi/app/global"
	"github.com/joho/godotenv"
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