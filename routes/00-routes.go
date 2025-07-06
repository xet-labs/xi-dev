package routes

import (
	"html/template"
	"log"
	"xi/conf"
	"xi/util"

	"github.com/gin-gonic/gin"
)

type RouteStruct struct {
	Engine    *gin.Engine
	Tmpl      *template.Template
	templates []string
}

var Route = &RouteStruct{
	templates: conf.View.Templates,
}
var r *gin.Engine

// Initializes all routes and templates
func (rt *RouteStruct) Init(engine *gin.Engine) {
	rt.Engine = engine
	r = engine

	rt.Tmpl = util.GetNewTmpl("base", ".html", rt.templates...) // Load templates
	r.SetHTMLTemplate(rt.Tmpl)

	// Register routes
	rt.registerCore()
	rt.registerBlogApi()
	rt.registerRes()
	rt.registerStatic()
	rt.registerDebug()

	// Optional: Middleware (e.g. gzip)
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	log.Println("✅ Routes registered..")
}
