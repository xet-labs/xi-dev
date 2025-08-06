package routes

import (
	"xi/app/cntr"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

func (rt *RouteStruct) registerDebug() {
	r.GET("/d", cntr.Debug(r))
	r.GET("/d/m", func(c *gin.Context) {
		// c.Data(200, "application/json", lib.Cfg.AllJsonPretty())
		c.Data(200, "application/json", lib.Cfg.AllJsonStructPretty())
	})
	r.GET("/d/ms", func(c *gin.Context) {
		// c.Data(200, "application/json", lib.Cfg.AllJsonPretty())
		c.Data(200, "application/json", lib.Cfg.AllJsonStructPretty())
	})
	r.GET("/d/j", func(c *gin.Context) {
		// c.Data(200, "application/json", lib.Cfg.AllJsonPretty())
		c.Data(200, "application/json", lib.Cfg.AllJsonPretty())
	})
		r.GET("/d/js", func(c *gin.Context) {
		// c.Data(200, "application/json", lib.Cfg.AllJsonPretty())
		c.Data(200, "application/json", lib.Cfg.AllJsonStructPretty())
	})
}
