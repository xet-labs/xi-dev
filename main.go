// main
package main

import (
	"xi/app/lib"
	"xi/app/routes"
	"xi/app/service"

	"github.com/gin-gonic/gin"
)

func init() {
	service.Init()
}

func main() {
	// Init Gin Engine
	gin.SetMode(lib.Env.Get("APP_MODE", "release"))
	app := gin.Default()

	// Init routes
	routes.Route.Init(app)

	// Init server
	err := service.InitServer(app)
	if err != nil {
		return
	}
}
