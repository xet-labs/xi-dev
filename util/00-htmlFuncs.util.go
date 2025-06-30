package util

import (
	"html"
	"html/template"
	"net/url"
	"strings"

	"xi/util/vtl"
)

var HtmlFuncs = template.FuncMap{
	"csslink":    vtl.Csslink,
	"jslink":     vtl.Jslink,
	"htmlEscape": html.EscapeString,
	"join":       strings.Join,
	"urlEscape":  url.QueryEscape,
}
