// main
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"xi/app/services"
	"xi/app/utils"
	"xi/routes"
)

func init() {
	services.InitEnv()
	services.InitDB()
}

func main() {

	// Init Gin router
	gin.SetMode(utils.Env("GIN_MODE", "debug"))
	app := gin.Default()

	// Register routes
	routes.Init(app)

	// Init server
	if err := services.InitServer(app); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
