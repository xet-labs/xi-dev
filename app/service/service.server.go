// services/server
package service

import (
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// InitServer start the web server
func InitServer(app *gin.Engine) error {
	appPort := cfg.App.Port

	log.Info().Msgf("\a\033[1;94mApp running \033[0;34m'http://localhost:%s'%s\033[0m\n", appPort,
		func() string {
			if url := cfg.Brand.Url; url != "" {
				return ", '" + url + "'"
			}
			return ""
		}())

	// Start Web-Server
	if err := app.Run(":" + appPort); err != nil {
		log.Error().Msgf("Failed to start server: %v", err)
		return err
	}

	return nil
}
