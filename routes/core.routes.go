package routes

import (
	"xi/app/cntr"
	"xi/app/lib"
)

func (rt *RouteStruct) registerCore() {
	cfg:=lib.Cfg
	r.GET(cfg.Get("view.pages.home.route").(string), cntr.Page.Tcnt("home", cfg.Get("view.pages.home.file").(string)))
}
