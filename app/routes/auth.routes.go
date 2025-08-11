package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerAuth() {

	authApi := r.Group("/api")
	{
		authApi.GET("/login", cntr.Auth.ShowLogin)
		authApi.GET("/login/:uid/:id", cntr.Auth.ShowLogin)
		authApi.POST("login", cntr.Auth.Login)

		authApi.GET("/logout", cntr.Auth.ShowLogout)
		authApi.GET("/logout/:uid/:id", cntr.Auth.ShowLogout)
		authApi.POST("/logout", cntr.Auth.Logout)

		authApi.GET("signup", cntr.Auth.ShowSignup)
		authApi.GET("/signup/:uid/:id", cntr.Auth.ShowSignup)
		authApi.POST("/signup", cntr.Auth.Signup)

		authApi.GET("/signout", cntr.Auth.ShowSignout)
		authApi.GET("/signout/:uid/:id", cntr.Auth.ShowSignout)
		authApi.POST("/signout", cntr.Auth.Signout)
	}

	login := r.Group("/login")
	{
		login.GET("", cntr.Auth.ShowLogin)
		login.GET("/:uid/:id", cntr.Auth.ShowLogin)
		login.POST("", cntr.Auth.Login)
	}
	logout := r.Group("/logout")
	{
		logout.GET("/logout", cntr.Auth.ShowLogout)
		logout.POST("/logout", cntr.Auth.Logout)
	}
	signup := r.Group("/signup")
	{
		signup.GET("", cntr.Auth.ShowSignup)
		signup.GET("/:uid/:id", cntr.Auth.ShowSignup)
		signup.POST("", cntr.Auth.Signup)
	}
	signout := r.Group("/logout")
	{
		signout.GET("/signout", cntr.Auth.ShowSignout)
		signout.POST("/signout", cntr.Auth.Signout)
	}
}
