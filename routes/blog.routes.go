package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerBlog() {
	blogApi := r.Group("api/blog")
	{
		blogApi.GET("", cntr.BlogApi.Index)
		blogApi.POST("/:id", cntr.BlogApi.Post)
		blogApi.GET("/:id", cntr.BlogApi.Show)
		blogApi.PUT("/:id", cntr.BlogApi.Put)
		blogApi.DELETE("/:id", cntr.BlogApi.Delete)
	}
}
