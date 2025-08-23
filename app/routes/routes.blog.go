package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerBlog() {
	blogApi := r.Group("api/blog") // route /api/blog
	{
		blogApi.GET("", cntr.BlogApi.Index)
		blogApi.GET("/:uid/:id", cntr.BlogApi.Show)
		blogApi.POST("/:uid/:id", cntr.BlogApi.Post)
		blogApi.PUT("/:uid/:id", cntr.BlogApi.Put)
		blogApi.DELETE("/:uid/:id", cntr.BlogApi.Delete)
	}

	// r.GET("", cntr.Blog.Index)         // route /blog
	blogs := r.Group("/blog/:uid/:id") // route /blog/*
	{
		blogs.GET("", cntr.Blog.Show)
		blogs.POST("", cntr.Blog.Post)
		blogs.PUT("", cntr.Blog.Put)
		blogs.DELETE("", cntr.Blog.Delete)
	}
}
