package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerMm() {
	r.GET("/mm", cntr.Res.Css.Get)
}
