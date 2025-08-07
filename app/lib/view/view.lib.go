package view

import (
	"bytes"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"xi/app/lib/db"
	"xi/app/lib/file"
	"xi/app/lib/minify"
	"xi/app/model"
	"xi/view/htmlfn"

	"github.com/gin-gonic/gin"
)

type ViewLib struct {
	Ecli    *gin.Engine
	Tcli    *template.Template
	RawTcli *template.Template
}

var View = &ViewLib{}

var HtmlFuncs = template.FuncMap{
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
		Funcs(HtmlFuncs).
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

// Render from Cache
func (v *ViewLib) RenderCache(c *gin.Context, page string, refKey string) bool {
	if cached, err := db.Rdb.GetBytes(refKey); err == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", cached)
		return true
	}
	return false
}

// Render and Cache Minified HTML
func (v *ViewLib) RenderAndCache(c *gin.Context, page string, refKey string, P model.PageParam) bool {
	var buf bytes.Buffer
	if err := v.Tcli.ExecuteTemplate(&buf, page, gin.H{"P": P}); err != nil {
		log.Printf("[View] RenderAndCache ERR: %s: %v", refKey, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Minify HTML
	minified, err := minify.Minify.Html(buf.Bytes())
	if err != nil {
		log.Printf("[View] RenderAndCache ERR: for %s: %v", refKey, err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		return true
	}

	// Send response and cache
	c.Data(http.StatusOK, "text/html; charset=utf-8", minified)
	go func(data any) {
		if err := db.Rdb.Set(refKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", refKey, err)
		}
	}(minified)

	return true
}
