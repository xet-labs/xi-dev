package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerCore() {
	r.GET("/", cntr.Page("pages/home"))
	r.GET("/b", cntr.Page("pages/blogs"))
	r.GET("/res/css/app.css", cntr.Res.Css)

	r.GET("/d", cntr.D) // debug
}
