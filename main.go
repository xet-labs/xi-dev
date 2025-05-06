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
	// cntr.InitBlogDB()
	// Init Gin router
	gin.SetMode(utils.Env("GIN_MODE", "debug"))
	router := gin.Default()

	// Register routes
	routes.Init(router)

	// Init server
	log.Printf("Init webserver 'http://localhost:5000'%s",
	func() string { if u := utils.Env("APP_URL"); u != "" { return ", 'http://" + u + "'" }; return "" }())

	if err := router.Run(":5000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
