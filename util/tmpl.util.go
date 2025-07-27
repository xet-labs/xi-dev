package util

import (
	"html"
	"html/template"
	"log"
	"net/url"
	"strings"

	"xi/util/vtl"
)

var HtmlFuncs = template.FuncMap{
	"formatTime": vtl.FormatTime,
	"isSlice":    vtl.IsSlice,
	"len":        vtl.Len,
	"linkCss":    vtl.Csslink,
	"linkJs":     vtl.Jslink,
	"slice":      vtl.Slice,
	"htmlEscape": html.EscapeString,
	"join":       strings.Join,
	"urlEscape":  url.QueryEscape,
}

func NewTmpl(Name, ext string, dirs ...string) *template.Template {
	files, err := GetExtFiles(ext, dirs...)
	if err != nil {
		log.Fatalf("Route: template load error: %v", err)
	}

	tmpl := template.Must(template.New(Name).
		Funcs(HtmlFuncs).
		// .Funcs(strutil.Funcs).
		// .Funcs(timeutil.Funcs).
		ParseFiles(files...),
	)

	return tmpl
}
