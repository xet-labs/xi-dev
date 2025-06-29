// routes
package routes

import (
	"xi/app/cntr"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/gzip"
)

func Init(r *gin.Engine) {

	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", cntr.Root)

	r.GET("/res/css/app.css", cntr.Res.Css)

	blog := r.Group("/blog")
	{
		blog.GET("", cntr.Blog.Index)
		blog.POST("/:id", cntr.Blog.Post)
		blog.GET("/:id", cntr.Blog.Show)
		blog.PUT("/:id", cntr.Blog.Put)
		blog.DELETE("/:id", cntr.Blog.Delete)
	}

	r.NoRoute(func(c *gin.Context) { c.File("./public" + c.Request.URL.Path) })
	// r.Static("/assets", "./assets")

	// debug
	r.GET("/d", cntr.D)
}
