package routes

import (
	"fmt"
	"html/template"
	"log"
	"xi/app/cntr"
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

	for title, page := range conf.View.Pages {
		route := "/t" + page["route"].(string)
		r.GET(route, cntr.PageTmpl(title, page["tmpl"].(string)))
	}

	// Register routes
	rt.registerCore()
	rt.registerBlog()
	rt.registerRes()
	rt.registerStatic()

	// Optional: Middleware (e.g. gzip)
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	log.Println("✅ Routes registered..")
	for _, route := range r.Routes() {
		fmt.Printf("%s\t%s\t→ %s\n", route.Method, route.Path, route.Handler)
	}
}
