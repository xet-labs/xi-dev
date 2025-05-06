// services/server
package services

import (
	"log"
	"xi/app/utils"
	// "xi/app/global"
	"github.com/gin-gonic/gin"
)

// InitServer starts the web server
func InitServer(router *gin.Engine) error {
	serverPort := utils.Env("APP_PORT", "5000")

	log.Printf("\033[1;94mServer running \033[0;34m'http://localhost:%s'%s\033[0m\n", serverPort,
	func() string {
		if u := utils.Env("APP_URL"); u != "" {
			return ", 'http://" + u + "'"
		}
		return ""
	}())

	return router.Run(":" + serverPort) //&& { global.ServerInitialized = true}
}
