package routes

import (
	"xi/app/cntr"
	"xi/conf"
)

func (rt *RouteStruct) registerCore() {
	// r.GET("/", cntr.Page("pages/home"))
	for title, page := range conf.View.Pages {
		r.GET(page["route"].(string), cntr.Page(rt.Tmpl, title, page["file"].(string)))
	}
}
