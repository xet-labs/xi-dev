package routes

import (
	"xi/app/lib"
	"xi/app/lib/cfg"

		"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

type RouteStruct struct {
	templates []string
}

var Route = &RouteStruct{
	templates: cfg.View.TemplateDir,
}

var r *gin.Engine

// Initializes all routes and templates
func (rt *RouteStruct) Init(engine *gin.Engine) {
	lib.View.Ecli = engine	// Engine_cli
	r = engine

	// Register templates
	r.SetHTMLTemplate(lib.View.NewTmpl("base", ".html", rt.templates...))

	// Register routes
	rt.registerCore()
	rt.registerBlog()
	rt.registerRes()
	rt.registerAuth()
	rt.registerDebug()

	// Optional: Middleware (e.g. gzip)
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	log.Info().Msg("Routes registered")
}
