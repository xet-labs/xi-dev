package routes

import (
	"xi/app/cntr"
	"xi/conf"
)

func (rt *RouteStruct) registerCore() {
	home := conf.View.Pages["home"]
	r.GET(home["route"].(string), cntr.Page.Tcnt("home", home["file"].(string)))
	
}
