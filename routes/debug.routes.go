package routes

import (
	"xi/app/cntr"
	"xi/conf"
)

func (rt *RouteStruct) registerDebug() {
	r.GET("/d", cntr.Debug(r))

	for title, page := range conf.View.Pages {
		route := "/t" + page["route"].(string)
		r.GET(route, cntr.PageTmpl(title, page["tmpl"].(string)))
	}
}
