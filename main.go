// main
package main

import (
	"xi/app/lib"
	"xi/app/service"
	"xi/routes"

	"github.com/gin-gonic/gin"
)
var Env = lib.Env

func init() {
	service.Init()
}

func main() {
	// Init Gin router
	gin.SetMode(Env.Get("APP_MODE", "debug"))
	app := gin.Default()

	// Register routes
	routes.Init(app)

	// Init server
	service.InitServer(app)
}
