package view

import (
	"html"
	"html/template"
	"net/url"
	"strings"

		"github.com/rs/zerolog/log"

	"xi/app/lib/file"
	"xi/view/htmlfn"
)

var HtmlFn = template.FuncMap{
	"formatTime":   htmlfn.FormatTime,
	"isSlice":      htmlfn.IsSlice,
	"len":          htmlfn.Len,
	"linkCss":      htmlfn.LinkCss,
	"linkCssSlice": htmlfn.LinkCssSlice,
	"linkJs":       htmlfn.LinkJs,
	"linkJsSlice":  htmlfn.LinkJsSlice,
	"linkLib":      htmlfn.LinkLib,
	"linkLibSlice": htmlfn.LinkLibSlice,
	"slice":        htmlfn.Slice,
	"htmlEscape":   html.EscapeString,
	"join":         strings.Join,
	"urlEscape":    url.QueryEscape,
}

func (v *ViewLib) NewTmpl(Name, ext string, dirs ...string) *template.Template {
	files, err := file.File.GetExt(ext, dirs...)
	if err != nil {
		log.Fatal().Msgf("View NewTmpl: Load err: %v", err)
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
		if rawTcli, err := tcli.Clone(); err == nil {
			v.RawTcli = rawTcli
		} else {
			log.Fatal().Msgf("View NewTmpl: Clone err: %s: %v", Name, err)
		}
	}

	return tcli
}
