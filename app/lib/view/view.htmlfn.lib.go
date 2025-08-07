package view

import (
	"html"
	"html/template"
	"log"
	"net/url"
	"strings"

	"xi/app/lib/file"
	"xi/view/htmlfn"
)

var HtmlFn = template.FuncMap{
	"formatTime": htmlfn.FormatTime,
	"isSlice":    htmlfn.IsSlice,
	"len":        htmlfn.Len,
	"linkCss":    htmlfn.Csslink,
	"linkJs":     htmlfn.Jslink,
	"slice":      htmlfn.Slice,
	"htmlEscape": html.EscapeString,
	"join":       strings.Join,
	"urlEscape":  url.QueryEscape,
}


func (v *ViewLib) NewTmpl(Name, ext string, dirs ...string) *template.Template {
	files, err := file.File.GetExt(ext, dirs...)
	if err != nil {
		log.Fatalf("Route: template load error: %v", err)
	}

	tcli := template.Must(template.New(Name).
		Funcs(HtmlFn).
		// Funcs(HtmlFuncs).
		// Funcs(timeutil.Funcs).
		ParseFiles(files...),
	)

	// Store instance globally so it can be used alter by other functions for rendering pages
	if v.Tcli == nil {
		v.Tcli = tcli
		if rawTcli, err := tcli.Clone(); err == nil{
			v.RawTcli = rawTcli
		} else {
			log.Printf("[View] NewTmpl clone ERR: %s: %v", Name, err)
		}
	}

	return tcli
}
