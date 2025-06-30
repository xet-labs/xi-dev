package data

import "xi/app/lib"

var View = map[string]any{

	"appName":   lib.Env.Get("APP_NAME"),
	"appDomain": lib.Env.Get("APP_DOMAIN"),
	"appUrl":    lib.Env.Get("APP_URL"),

	"navMenu": map[string]any{
		"Home":    "/",
		"Blog":    "/blog",
		"Product": "#",
		"Support": "#",
		"Contact": "#",
	},
}
