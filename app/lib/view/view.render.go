package view

import (
	"bytes"
	"net/http"
	"time"

	"xi/app/lib/db"
	"xi/app/lib/minify"
	"xi/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Render struct {
	c    *gin.Context
	data []byte
	ok   bool
	err  error
}

func (r Render) Html() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/html; charset=utf-8", r.data)
		return true
	}
	return false
}

func (r Render) Json() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "application/json", r.data)
		return true
	}
	return false
}

func (r Render) Css() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/css; charset=utf-8", r.data)
		return true
	}
	return false
}

// Render from Cache
func (v *ViewLib) OutCache(c *gin.Context, rdbKey string) Render {
	cache, err := db.Rdb.GetBytes(rdbKey)
	return Render{c: c, data: cache, ok: err == nil}
}

// Render and Cache Minified HTML
func (v *ViewLib) OutHtmlLyt(c *gin.Context, P model.PageParam, args ...string) bool {
	argsLen := len(args)

	// Render html via template
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, P.Layout, gin.H{"P": P}); err != nil {
		log.Error().Msgf("View OutHtmlLyt execTemplate: %s: %v", c.Request.URL, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Minify HTML
	pageMin, err := minify.Minify.Html(page.Bytes())
	if err != nil {
		log.Error().Msgf("View OutHtmlLyt minify: for %s: %v", c.Request.URL, err)

		// Serve the response with optional cache if rdbKey is provided in args[0]
		c.Data(http.StatusOK, "text/html; charset=utf-8", page.Bytes())
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
	argsLen := len(args)

	// Minify the CSS
	cssMin, err := minify.Minify.CssHybrid(css)
	if err != nil {
		log.Error().Msgf("View OutCss Minify: %v", err)
		c.Data(http.StatusOK, "text/css; charset=utf-8", css)
		return true
	}

	// Serve the response with optional cache if rdbKey is provided in args[0]
	c.Data(http.StatusOK, "text/css; charset=utf-8", cssMin)
	if argsLen > 0 && args[0] != "" {
		go func(data any) { db.Rdb.Set(args[0], data, 10*time.Minute) }(cssMin)
	}
	return true
}

func (v *ViewLib) OutJson(c *gin.Context, rdbKey string, P model.PageParam, cache ...bool) bool {
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, P.Layout, gin.H{"P": P}); err != nil {
		log.Printf("[View] OutJson ERR: %s: %v", rdbKey, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Minify HTML
	pageMin, err := minify.Minify.Html(page.Bytes())
	if err != nil {
		log.Printf("[View] OutJson ERR: for %s: %v", rdbKey, err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", page.Bytes())
		return true
	}

	// Send response and cache
	c.Data(http.StatusOK, "text/html; charset=utf-8", pageMin)
	if cache[0] {
		go func(data any) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(pageMin)
	}
	return true
}
