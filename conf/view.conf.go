package conf

import "xi/app/lib"

var View = &ViewStruct{

	Templates: []string{"views/layout", "views/partials", "views/pages"},

	Layout: "layout/base",

	PageData: map[string]any{
		"appName":   lib.Env.Get("APP_NAME"),
		"appDomain": lib.Env.Get("APP_DOMAIN"),
		"appUrl":    lib.Env.Get("APP_URL"),

		"tags":         []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
		"featured_img": []string{"/res/static/brand/brand.svg"},

		"navMenu": []MenuItem{
			{Label: "Home", Href: "/"},
			{Label: "Blog", Href: "/blog"},
			{Label: "Product", Href: "#"},
			{Label: "Support", Href: "#"},
			{Label: "Contact", Href: "#"},
		},
		"js99": []string{
			"/res/js/jquery/jquery.min.js",
			"/res/js/app.js",
		},
	},

	Pages: map[string]map[string]any{
		"Home": {
			"route":       "/",
			"file":        "views/pages/home.html",
			"tmpl":        "pages/home",
			"menu":        "Home",
			"title":       "XetIndustries",
			"description": "XetIndustries is a Collaborative Platform for Makers, Creators, and Developers.",
			"meta": map[string]any{
				"canonical": "https://xetindustries.com/",
				"tags":      []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
			},
		},
		"Blog": {
			"route":       "/blog",
			"file":        "views/pages/blogs.html",
			"tmpl":        "pages/blogs",
			"menu":        "Blog",
			"subBrand":    "Blog",
			"title":       "Blog | XetIndustries",
			"description": "Discover Insightful Articles, Expert Tips and Inspiring Ideas. Write, Share and Connect with a thriving like-minded community.",
			"js99": []string{
				"/res/js/jquery/jquery.min.js",
				"/res/js/app.js",
				"/res/js/blogs.js",
			},
			"meta": map[string]any{
				"canonical": "https://xetindustries.com/blog",
			},
		},
		// "Contact": {
		// 	"route": "/contact",
		// 	"file":  "views/pages/contact.html",
		// },
	},
}

type ViewStruct = struct {
	Layout    string
	Templates []string
	PageData  map[string]any
	Pages     map[string]map[string]any
}

type MenuItem struct {
	Label string
	Href  string
}
