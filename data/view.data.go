package data

import "xi/app/lib"

type MenuItem struct {
	Label string
	Href  string
}

var View = map[string]any{
	"appName":      lib.Env.Get("APP_NAME"),
	"appDomain":    lib.Env.Get("APP_DOMAIN"),
	"appUrl":       lib.Env.Get("APP_URL"),
	
	"tags":         []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
	"featured_img": []string{"/res/static/brand/brand.svg"},

	"navMenu": []MenuItem{
		{Label: "Home", Href: "/"},
		{Label: "Blog", Href: "/blog"},
		{Label: "Product", Href: "#"},
		{Label: "Support", Href: "#"},
		{Label: "Contact", Href: "#"},
	},

	"jsInc99": []string{
		"/res/js/jquery/jquery.min.js",
		"/res/js/app.js",
	},
}
