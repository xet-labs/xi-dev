package routes

import (
	"xi/app/lib"
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
)

type RouteStruct struct{}

var (
	Route = &RouteStruct{}
	r     *gin.Engine
)

// Initializes all routes and templates
func (rt *RouteStruct) Init(engine *gin.Engine) {
	lib.View.Ecli = engine // Engine_cli
	r = engine

	// Register templates
	r.SetHTMLTemplate(lib.View.NewTmpl("main", ".html", cfg.View.TemplateDirs...))

	// Register routes
	rt.registerMm()
	rt.registerStatic()
	rt.registerCore()
	rt.registerBlog()
	rt.registerAuth()
	rt.registerDebug()
	rt.registerRes()

	// Optional: Middleware (e.g. gzip)
	// r.Use(gzip.Gzip(gzip.DefaultCompression))
}
