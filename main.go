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
	// Init Gin Engine
	gin.SetMode(Env.Get("APP_MODE", "release"))
	app := gin.Default()

	// Init routes
	routes.Route.Init(app)

	// Init server
	service.InitServer(app)
}
