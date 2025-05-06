// routes
package routes

import (
	"net/http"

	"xi/app/cntr"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/gzip"
)

func Init(r *gin.Engine) {
	
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	
	r.GET("/", cntr.Root)
	
	blog := r.Group("/blog"); {
		blog.GET("", 		cntr.Blog.Index)
		blog.POST("/:id", 	cntr.Blog.Post)
		blog.GET("/:id", 	cntr.Blog.Show)
		blog.PUT("/:id", 	cntr.Blog.Put)
		blog.DELETE("/:id", cntr.Blog.Delete)
	}
	
	resDir := http.FileServer(http.Dir("./public/res"))
	r.GET("/res/*filepath", gin.WrapH(http.StripPrefix("/res/", resDir)))

	// Serve /static/* from ./public/static
	staticDir := http.FileServer(http.Dir("./public/static"))
	r.GET("/static/*filepath", gin.WrapH(http.StripPrefix("/static/", staticDir)))

	r.NoRoute(func(c *gin.Context) { c.File("./public" + c.Request.URL.Path) })
	
	//- debug
	r.GET("/d", cntr.D)
}
