package routes

import (
	"html/template"
	"log"
	"xi/app/lib"
	"xi/conf"
	"xi/util"

	"github.com/gin-gonic/gin"
)

type RouteStruct struct {
	Ecli      *gin.Engine
	Tcli      *template.Template
	templates []string
}

var Route = &RouteStruct{
	templates: conf.View.TemplateDir,
}

var r *gin.Engine

// Initializes all routes and templates
func (rt *RouteStruct) Init(engine *gin.Engine) {
	rt.Ecli = engine
	r = rt.Ecli

	rt.Tcli = util.NewTmpl("base", ".html", rt.templates...) // Load templates
	r.SetHTMLTemplate(rt.Tcli)

	// propagate to global lib.View
	lib.View.Ecli = rt.Ecli
	lib.View.Tcli = rt.Tcli
	lib.View.RawTcli, _ = rt.Tcli.Clone()

	// Register routes
	rt.registerAuth()
	rt.registerCore()
	rt.registerBlog()
	rt.registerRes()
	rt.registerDebug()

	// Optional: Middleware (e.g. gzip)
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	log.Println("✅ Routes registered..")
}
