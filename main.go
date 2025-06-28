// main
package main

import (
	"xi/app/service"
	"xi/app/util"
	"xi/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	service.Init()
}

func main() {
	// Init Gin router
	gin.SetMode(util.Env("GIN_MODE", "debug"))
	app := gin.Default()

	// Register routes
	routes.Init(app)

	// Init server
	service.InitServer(app)
}
