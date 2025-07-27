package routes

import (
	"xi/app/cntr"
	// "xi/conf"
)

func (rt *RouteStruct) registerDebug() {
	r.GET("/d", cntr.Debug(r))
}
