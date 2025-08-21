package routes

import (
	"xi/app/lib"
	"xi/app/lib/cfg"
)

func (rt *RouteStruct) registerDynamic() {
	for page := range cfg.View.Pages {
		if cfg.View.Pages[page].Enable == nil || *cfg.View.Pages[page].Enable{
			r.GET(cfg.View.Pages[page].Route, lib.View.OutPageHandler(cfg.View.Pages[page]))
		}
	}
}
