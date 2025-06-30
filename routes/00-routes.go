package routes

import (
	"log"
	"xi/util"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/gzip"
)

// Struct and Global instance
type RouteStruct struct {
	Engine *gin.Engine
}
var Route = &RouteStruct{}

// Globally defined so other methods can access the *engine via 'r' after Init method updates it
var r *gin.Engine

// Initializes all routes and templates
func (rt *RouteStruct) Init(engine *gin.Engine) {
	rt.Engine = engine
	r = engine 

	// Load templates
	r.SetHTMLTemplate(util.LoadTemplates("views", ".html"))

	// Optional middleware
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Register grouped routes
	rt.registerCore()
	rt.registerBlog()
	rt.registerStatic()

	log.Println("✅ Routes registered..")
}
