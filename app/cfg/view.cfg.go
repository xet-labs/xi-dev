package cfg

// import (
// 	"xi/app/lib"
// 	"xi/app/lib/schema"
// )

// var View = &schema.ViewConf{
// 	CssDir:      []string{"view/partials"},
// 	TemplateDir: []string{"view/layout", "view/partials", "view/pages"},

// 	Pages: map[string]schema.PageParam{
// 		"Default": {
// 			Layout: "layout/base",
// 			App: schema.AppConf{
// 				Name:        lib.Env.Get("APP_NAME"),
// 				Domain:      lib.Env.Get("APP_DOMAIN"),
// 				Url:         lib.Env.Get("APP_URL"),
// 				Tags:        []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
// 				FeaturedImg: []string{"/res/static/brand/brand.svg"},
// 			},
// 			NavMenu: []schema.MenuItem{
// 				{Label: "Home", Href: "/"},
// 				{Label: "Blog", Href: "/blog"},
// 				{Label: "Product", Href: "#"},
// 				{Label: "Support", Href: "#"},
// 				{Label: "Contact", Href: "#"},
// 			},
// 			Js99: []string{
// 				"/res/js/jquery/jquery.min.js",
// 				"/res/js/app.js",
// 			},
// 		},
// 		"home": {
// 			Route:       "/",
// 			File:        "view/pages/home.part.html",
// 			Template:    "page/home",
// 			Menu:        "Home",
// 			Title:       "XetIndustries",
// 			Description: "A Platform for Makers, Creators, and Developers.",
// 			Meta: &schema.PageMeta{
// 				Canonical: "https://xetindustries.com/",
// 				Tags:      []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
// 			},
// 		},
// 		"blogs": {
// 			Route:       "/blog",
// 			Template:    "layout/blogs",
// 			Menu:        "Blog",
// 			SubBrand:    "Blog",
// 			Title:       "Blog | XetIndustries",
// 			Description: "Discover Insightful Articles, Expert Tips and Inspiring Ideas. Write, Share and Connect with a thriving like-minded community.",
// 			Js99: []string{
// 				"/res/js/jquery/jquery.min.js",
// 				"/res/js/app.js",
// 				"/res/js/blogs.js",
// 			},
// 			Meta: &schema.PageMeta{
// 				Canonical: "https://xetindustries.com/blog",
// 			},
// 		},
// 		"blog": {
// 			Route:    "/blog/*",
// 			Template: "layout/blog",
// 			Menu:     "Blog",
// 			SubBrand: "Blog",
// 			LibHLJS:  true,
// 			Js99: []string{
// 				"/res/js/jquery/jquery.min.js",
// 				"/res/js/app.js",
// 				"/res/util/copy-code.util.js",
// 				"/res/util/h2-clickable.util.js",
// 			},
// 		},
// 	},
// }
