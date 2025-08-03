package cfg

import "xi/app/lib"

type ViewConf = struct {
	CssDir    []string
	Layout    string
	TemplateDir []string
	PageData  map[string]any
	Pages     map[string]map[string]any
}

type MenuItem struct {
	Label string
	Href  string
}

var View = &ViewConf{
	CssDir:    []string{"view/partials"},
	Layout:    "layout/base",
	TemplateDir: []string{"view/layout", "view/partials", "view/pages"},

	PageData: map[string]any{
		"appName":   lib.Env.Get("APP_NAME"),
		"appDomain": lib.Env.Get("APP_DOMAIN"),
		"appUrl":    lib.Env.Get("APP_URL"),

		"appTags":         []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
		"appFeatured_img": []string{"/res/static/brand/brand.svg"},

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
		"home": {
			"route":       "/",
			"file":		   "view/pages/home.part.html",
			"tmpl":        "page/home",
			"menu":        "Home",
			"title":       "XetIndustries",
			"description": "A Platform for Makers, Creators, and Developers.",
			"meta": map[string]any{
				"canonical": "https://xetindustries.com/",
				"tags":      []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
			},
		},
		"blogs": {
			"route":       "/blog",
			"tmpl":        "layout/blogs",
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
		"blog": {
			"route":    "/blog/*",
			"tmpl":     "layout/blog",
			"menu":     "Blog",
			"subBrand": "Blog",
			"lib_hljs": true,
			"js99": []string{
				"/res/js/jquery/jquery.min.js",
				"/res/js/app.js",
				// "/res/js/blog.js",
				"/res/util/copy-code.util.js",
				"/res/util/h2-clickable.util.js",
			},
		},
		// "Contact": {
		// 	"route": "/contact",
		// 	"file":  "views/pages/contact.html",
		// },
	},
}
