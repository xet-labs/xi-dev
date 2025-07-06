package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerBlogApi() {
	blogApi := r.Group("api/blog")
	{
		blogApi.GET("", cntr.BlogApi.Index)
		blogApi.GET("/:id", cntr.BlogApi.Show)
		blogApi.POST("/:id", cntr.BlogApi.Post)
		blogApi.PUT("/:id", cntr.BlogApi.Put)
		blogApi.DELETE("/:id", cntr.BlogApi.Delete)
	}
}
