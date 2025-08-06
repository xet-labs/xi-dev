package lib

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"xi/app/model"

	"github.com/gin-gonic/gin"
)

type ViewLib struct {
	Ecli    *gin.Engine
	Tcli    *template.Template
	RawTcli *template.Template
}

var View = &ViewLib{}

// Render from Cache
func (v *ViewLib) RenderCache(c *gin.Context, page string, refKey string) bool {
	if cached, err := Redis.GetBytes(refKey); err == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", cached)
		return true
	}
	return false
}

// Render and Cache Minified HTML
func (v *ViewLib) RenderAndCache(c *gin.Context, page string, refKey string, P model.PageParam) bool {
	var buf bytes.Buffer
	if err := v.Tcli.ExecuteTemplate(&buf, page, gin.H{"P": P}); err != nil {
		log.Printf("Render error for %s: %v", refKey, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Minify HTML
	minified, err := Minify.Html(buf.Bytes())
	if err != nil {
		log.Printf("Minify error for %s: %v", refKey, err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		return true
	}

	// Send response and cache
	c.Data(http.StatusOK, "text/html; charset=utf-8", minified)
	go func(data any) {
		if err := Redis.Set(refKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", refKey, err)
		}
	}(minified)

	return true
}
