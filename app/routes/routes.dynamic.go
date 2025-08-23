package routes

import (
	"xi/app/lib"
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
)

func (rt *RouteStruct) registerDynamic() {
	for name, page := range cfg.View.Pages {
		if page != nil && (page.Enable == nil || *page.Enable) {
			// r.GET(page.Route, lib.View.PageHandler(name))
			r.GET(page.Route, func(c *gin.Context) { lib.View.Page(c, cfg.View.Pages[name]) })
		}
	}
}
