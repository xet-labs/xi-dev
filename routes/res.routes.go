package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerRes() {
	r.GET("/res/css/app.css", cntr.Res.Css)
}