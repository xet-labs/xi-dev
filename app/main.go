package main

import (
	"xi/app/lib/cfg"
	"xi/app/routes"
	"xi/app/service"

	"github.com/gin-gonic/gin"
	// "xi/app/lib"
)

func main() {
	// Init services
	service.Init()

	// Init Gin Engine
	gin.SetMode(cfg.App.Mode)
	app := gin.Default()

	// Init routes
	routes.Route.Init(app)

	// Init server

	err := service.InitServer(app)
	if err != nil {
		return
	}
}
