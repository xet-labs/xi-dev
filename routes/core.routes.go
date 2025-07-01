package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerCore() {
	r.GET("/", cntr.Home)
	r.GET("/res/css/app.css", cntr.Res.Css)

	r.GET("/d", cntr.D) // debug
}
