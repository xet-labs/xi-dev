// services/server
package service

import (
	"log"
	"xi/app/util"
	// "xi/app/global"
	"github.com/gin-gonic/gin"
)

// InitServer start the web server
func InitServer(app *gin.Engine) error {
	appPort := util.Env("APP_PORT", "5000")

	log.Printf("\a\033[1;94mServer running \033[0;34m'http://localhost:%s'%s\033[0m\n", appPort,
	func() string {
		if url := util.Env("APP_URL"); url != "" {
			return ", 'http://" + url + "'"
		}
		return ""
	}())

	// Start Web-Server
	return app.Run(":" + appPort) //&& { global.ServerInitialized = true}
}
