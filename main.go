// main
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"xi/app/service"
	"xi/app/util"
	"xi/routes"
)

func init() {
	services.InitEnv()
	services.InitDB()
}

func main() {

	// Init Gin router
	gin.SetMode(util.Env("GIN_MODE", "debug"))
	app := gin.Default()

	// Register routes
	routes.Init(app)

	// Init server
	if err := services.InitServer(app); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
