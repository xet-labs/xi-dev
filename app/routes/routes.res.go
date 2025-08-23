package routes

import (
	"xi/app/cntr"

	"github.com/gin-gonic/gin"
)

func (rt *RouteStruct) registerRes() {
	r.GET("/res/css/*name", cntr.Res.Css.Get)

	r.NoRoute(func(c *gin.Context) { c.File("public" + c.Request.URL.Path) })
}
