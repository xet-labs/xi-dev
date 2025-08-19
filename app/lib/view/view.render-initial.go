package view

import (
	"bytes"
	"net/http"
	"time"
	"xi/app/lib/db"
	"xi/app/lib/util"
	"xi/app/lib/minify"
	"xi/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Render and Cache Minified HTML
func (v *ViewLib) OutHtmlLyt(c *gin.Context, P model.PageParam, args ...string) bool {
	
	// Render html via template
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, util.Str.Fallback(P.Layout, "layout/default"), gin.H{"P": P}); err != nil {
		log.Error().Err(err).Str("Page",c.Request.RequestURI).Msg("View.OutHtmlLyt.ExecTemplate")
		c.Status(http.StatusInternalServerError)
		return false
	}
	
	argsLen := len(args)
	
	// Minify HTML
	pageMin, err := minify.Minify.Html(page.Bytes())
	if err != nil {
		// Serve the response with optional cache if rdbKey is provided in args[0]
		c.Data(http.StatusOK, "text/html; charset=utf-8", page.Bytes())
		log.Error().Err(err).Str("Page",c.Request.RequestURI).Msg("View.OutHtmlLyt.minify")

		if argsLen > 0 && args[0] != "" {
			go func(data any) { db.Rdb.Set(args[0], data, 10*time.Minute) }(pageMin)
		}
		return true
	}

	// Serve the response with optional cache if rdbKey is provided in args[0]
	c.Data(http.StatusOK, "text/html; charset=utf-8", pageMin)
	if argsLen > 0 && args[0] != "" {
		go func(data any) { db.Rdb.Set(args[0], data, 10*time.Minute) }(pageMin)
	}
	return true
}

func (v *ViewLib) OutCss(c *gin.Context, css []byte, args ...string) bool {
	// Handle empty content
	if len(css) == 0 {
		c.Status(http.StatusNoContent) // 204
		return true
	}

	// Minify the CSS
	cssMin, err := minify.Minify.CssHybrid(css)
	if err != nil {
		c.Data(http.StatusOK, "text/css; charset=utf-8", css)
		log.Error().Err(err).Msg("View OutCss Minify")
		return true
	}

	// Serve the response with optional cache if rdbKey is provided in args[0]
	c.Data(http.StatusOK, "text/css; charset=utf-8", cssMin)
	if len(args) > 0 && args[0] != "" {
		go func(data any) { db.Rdb.Set(args[0], data, 10*time.Minute) }(cssMin)
	}
	return true
}

func (v *ViewLib) OutJson(c *gin.Context, rdbKey string, P model.PageParam, cache ...bool) bool {
	return true
}
