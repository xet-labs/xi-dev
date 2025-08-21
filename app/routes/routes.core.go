package routes

import (
	"xi/app/cntr"
	"xi/app/lib/cfg"
)

func (rt *RouteStruct) registerCore() {
	home := cfg.View.Pages["home"]
	r.GET("/", cntr.Page.Tcnt("home", home.ContentFile))
}
