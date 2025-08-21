package routes

import (
	"xi/app/lib"
	"xi/app/lib/cfg"
)

func (rt *RouteStruct) registerDynamic() {
	for _, page := range cfg.View.Pages {
		if page.Enable == nil || *page.Enable{
			r.GET(page.Route, lib.View.OutPageHandler(page))
		}
	}
}
