package routes

import (
	"xi/app/cntr"
	"xi/app/cfg"
)

func (rt *RouteStruct) registerCore() {
	home := cfg.View.Pages["home"]
	r.GET(home["route"].(string), cntr.Page.Tcnt("home", home["file"].(string)))
	
}
