package view

import (
	"bytes"
	"net/http"
	"time"
	"xi/app/lib/db"
	"xi/app/lib/minify"
	"xi/app/lib/util"
	"xi/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Render and Cache Minified HTML
func (v *ViewLib) OutHtmlLyt(c *gin.Context, p model.PageParam, args ...string) bool {
	
	// Render html via template
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, util.Str.Fallback(p.Layout, "layout/default"), gin.H{"P": p}); err != nil {
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

func (v *ViewLib) Page(c *gin.Context, p model.PageParam, args ...string) bool {
	// Detect Render Type and Format
	if p.Render != "" {
		rType := "part"
		rFormat := "file"

		switch rType { // Render Type
		case "part":
		case "body":
		case "full":
		default:
		}

		switch rFormat { // Render Format
		case "content":
		case "file":
		case "url":
		default:
		}
	}
	// Render html via template
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, util.Str.Fallback(p.Layout, "layout/default"), gin.H{"P": p}); err != nil {
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


