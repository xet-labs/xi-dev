package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerBlogApi() {
	blogApi := r.Group("api/blog")
	{
		blogApi.GET("", cntr.BlogApi.Index)
		blogApi.GET("/:uid/:id", cntr.BlogApi.Show)
		blogApi.POST("/:uid/:id", cntr.BlogApi.Post)
		blogApi.PUT("/:uid/:id", cntr.BlogApi.Put)
		blogApi.DELETE("/:uid/:id", cntr.BlogApi.Delete)
	}

		blog := r.Group("/blog")
	{
		// blog.GET("", cntr.BlogApi.Index)
		blog.GET("/:uid/:id", cntr.BlogApi.Show)
		blog.POST("/:uid/:id", cntr.BlogApi.Post)
		blog.PUT("/:uid/:id", cntr.BlogApi.Put)
		blog.DELETE("/:uid/:id", cntr.BlogApi.Delete)
	}
}
