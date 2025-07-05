package routes

import (
	"xi/app/cntr"
	"xi/conf"
)

func (rt *RouteStruct) registerTest() {
	for title, page := range conf.View.Pages {
		route := "/t" + page["route"].(string)
		r.GET(route, cntr.PageTmpl(title, page["tmpl"].(string)))
	}

}
