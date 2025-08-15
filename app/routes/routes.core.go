package routes

import (
	"xi/app/cntr"
	"xi/app/lib/cfg"
)

func (rt *RouteStruct) registerCore() {
	home := cfg.View.Pages["home"]
	r.GET(home.Route, cntr.Page.Tcnt("home", home.File))
}
