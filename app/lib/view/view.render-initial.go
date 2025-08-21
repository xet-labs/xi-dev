package view

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
	"xi/app/lib/cfg"
	"xi/app/lib/db"
	"xi/app/lib/minify"
	"xi/app/lib/util"
	"xi/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Render and Cache Minified HTML
func (v *ViewLib) OutHtmlLyt(c *gin.Context, p *model.PageParam, args ...string) bool {

	// Render html via template
	page := bytes.Buffer{}
	if err := v.Tcli.ExecuteTemplate(&page, util.Str.Fallback(p.Layout, "layout/default"), gin.H{"P": p}); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("page", c.Request.URL.Path).Msg("View OutHtmlLyt, ExecTemplate")
		return false
	}
	// page := buf.Bytes()

	// Minify HTML
	pageMin, err := minify.Minify.Html(page.Bytes())
	if err != nil {
		// Serve the response with optional cache if rdbKey is provided in args[0]
		c.Data(http.StatusOK, "text/html; charset=utf-8", page.Bytes())
		log.Error().Err(err).Str("page", c.Request.URL.Path).Msg("View OutHtmlLyt, minify")

		if p.Cache == nil || *p.Cache || cfg.App.ForceCachePage {
			rdbKey := util.ArrFallback(args, 0, c.Request.URL.Path)
			go func(data any) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(page)
		}
		return true
	}

	// Serve the response with optional cache if rdbKey is provided in args[0]
	c.Data(http.StatusOK, "text/html; charset=utf-8", pageMin)
	if p.Cache == nil || *p.Cache || cfg.App.ForceCachePage {
		rdbKey := util.ArrFallback(args, 0, c.Request.URL.Path)
		go func(data any) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(pageMin)
	}
	return true
}

func (v *ViewLib) OutPageHandler(p model.PageParam, args ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v.OutPage(c, &p, args...)
	}
}

func (v *ViewLib) OutPage(c *gin.Context, p *model.PageParam, args ...string) bool {
	rdbKey := c.Request.URL.Path

	// Try cache
	if v.OutCache(c, rdbKey).Html() {
		return true
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	p.Rt = map[string]any{}
	// Detect Render.Content Format
	switch p.Render {
	case "url":
	case "md":
	case "file":
		contentBytes, err := os.ReadFile(p.ContentFile)
		if err != nil {
			log.Error().Err(err).Str("page", c.Request.URL.Path).Msg("View Page, Read-file")
			c.Status(http.StatusInternalServerError)
			return false
		}
		p.Rt["Content"] = template.HTML(string(contentBytes))

	case "content":
		p.Rt["Content"] = template.HTML(p.Content)
	}

	// Render Type
	var page []byte
	switch p.Layout {
	case "raw":
		switch v := p.Rt["Content"].(type) {
		case []byte:
			page = v
		case string:
			page = []byte(v)
		case template.HTML:
			page = []byte(v)
		default:
			c.Status(http.StatusInternalServerError)
			log.Error().Str("type", fmt.Sprintf("%T", v)).Str("Page", c.Request.URL.Path).Msg("View Page, Unsupported content type in p.Rt[\"content\"]")
			return false
		}

	default:
		buf := bytes.Buffer{}
		if err := v.Tcli.ExecuteTemplate(&buf, util.Str.Fallback(p.Layout, "layout/default"), gin.H{"P": p}); err != nil {
			log.Error().Err(err).Str("page", c.Request.URL.Path).Msg("View Page, ExecTemplate")
			c.Status(http.StatusInternalServerError)
			return false
		}
		page = buf.Bytes()
	}

	// Minify HTML
	pageMin, err := minify.Minify.Html(page)
	if err != nil {
		// Serve the response with optional cache if rdbKey is provided in args[0]
		c.Data(http.StatusOK, "text/html; charset=utf-8", page)
		log.Error().Err(err).Str("page", c.Request.URL.Path).Msg("View.OutHtmlLyt.minify")

		if p.Cache == nil || *p.Cache || cfg.App.ForceCachePage {
			go func(data *[]byte) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(&page)
		}
		return true
	}

	// Serve the response with optional cache if rdbKey is provided in args[0]
	c.Data(http.StatusOK, "text/html; charset=utf-8", pageMin)
	if p.Cache == nil || *p.Cache || cfg.App.ForceCachePage {
		go func(data *[]byte) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(&pageMin)
	}
	return true
}
