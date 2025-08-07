package view

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"xi/app/lib/db"
	"xi/app/lib/minify"
	"xi/app/model"

	"github.com/gin-gonic/gin"
)
type ViewLib struct {
	Ecli    *gin.Engine
	Tcli    *template.Template
	RawTcli *template.Template
}

type RenderData struct {
	c *gin.Context
	data []byte
	ok   bool
	err  error
}

var View = &ViewLib{}

func (r RenderData) Html() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/html; charset=utf-8", r.data)
		return true
	}
	return false
}

func (r RenderData) Json() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "application/json", r.data)
		return true
	}
	return false
}

func (r RenderData) Css() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/css; charset=utf-8", r.data)
		return true
	}
	return false
}

// Render from Cache
func (v *ViewLib) RenderCache(c *gin.Context, rdbKey string) RenderData {
	cache, err := db.Rdb.GetBytes(rdbKey);
	return RenderData{c:c, data: cache, ok: err == nil}
}

// Render and Cache Minified HTML
func (v *ViewLib) RenderAndCache(c *gin.Context, rdbKey string, P model.PageParam) bool {
	var buf bytes.Buffer
	if err := v.Tcli.ExecuteTemplate(&buf, P.Layout, gin.H{"P": P}); err != nil {
		log.Printf("[View] RenderAndCache ERR: %s: %v", rdbKey, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Minify HTML
	minified, err := minify.Minify.Html(buf.Bytes())
	if err != nil {
		log.Printf("[View] RenderAndCache ERR: for %s: %v", rdbKey, err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		return true
	}

	// Send response and cache
	c.Data(http.StatusOK, "text/html; charset=utf-8", minified)
	go func(data any) { db.Rdb.Set(rdbKey, data, 10*time.Minute) }(minified)
	return true
}
