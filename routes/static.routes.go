package routes

import "github.com/gin-gonic/gin"

func (rt *RouteStruct) registerStatic() {
	r.NoRoute(func(c *gin.Context) { c.File("./public" + c.Request.URL.Path) })
	// r.Static("/assets", "./assets")
}
