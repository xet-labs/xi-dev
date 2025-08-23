package view

import (
	"html"
	"html/template"
	"net/url"
	"strings"
	"xi/app/lib/util"
	"xi/app/lib/view/htmlfn"

	"github.com/rs/zerolog/log"
)

var HtmlFn = template.FuncMap{
	"formatTime":   htmlfn.FormatTime,
	"genMeta":      htmlfn.GenMeta,
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

func (v *ViewLib) NewTmpl(name, ext string, dirs ...string) *template.Template {
	files, err := util.File.GetWithExt(ext, dirs...)
	if err != nil {
		log.Error().Err(err).Str("cli", name).Str("template-dir", util.Str.QuoteSlice(dirs)).
			Msg("View NewTmpl: couldnt get template files")
	}

	tcli := template.Must(template.New(name).
		Funcs(HtmlFn).
		// Funcs(HtmlFuncs).
		// Funcs(timeutil.Funcs).
		ParseFiles(files...),
	)

	// Store instance globally so it can be used alter by other functions for rendering pages
	if v.Tcli == nil {
		v.Tcli = tcli
		rawTcli, err := tcli.Clone()
		if err != nil {
			log.Error().Err(err).Str("cli", name).
				Msg("View NewTmpl: cli")
		}
		v.RawTcli = rawTcli
	}
	return tcli
}
