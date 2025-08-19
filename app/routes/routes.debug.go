package routes

import (
	"xi/app/cntr"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

func (rt *RouteStruct) registerDebug() {
	r.GET("/d", cntr.Debug(r))
	r.GET("/d/c", func(c *gin.Context) {
		c.Data(200, "application/json", lib.Conf.AllJsonPretty())
	})
	r.GET("/d/cs", func(c *gin.Context) {
		c.Data(200, "application/json", lib.Conf.AllJsonStructPretty())
	})
}
