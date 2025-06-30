package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerBlog() {
	blog := r.Group("/blog")
	{
		blog.GET("", cntr.Blog.Index)
		blog.POST("/:id", cntr.Blog.Post)
		blog.GET("/:id", cntr.Blog.Show)
		blog.PUT("/:id", cntr.Blog.Put)
		blog.DELETE("/:id", cntr.Blog.Delete)
	}
}
