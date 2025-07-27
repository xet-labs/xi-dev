package routes

import (
	"xi/app/cntr"
)

func (rt *RouteStruct) registerAuth() {

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
