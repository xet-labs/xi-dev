package view

import (
	"net/http"
	"time"
	"xi/app/lib/db"
	"xi/app/lib/minify"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Render struct {
	c    *gin.Context
	data []byte
	ok   bool
	err  error
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

func (v *ViewLib) OutJson(c *gin.Context, css []byte, args ...string) bool {
	return true
}


// Render from Cache
func (v *ViewLib) OutCache(c *gin.Context, rdbKey string) Render {
	cache, err := db.Rdb.GetBytes(rdbKey)
	return Render{c: c, data: cache, ok: err == nil}
}

// Helpers methods
func (r Render) Html() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/html; charset=utf-8", r.data)
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

func (r Render) Json() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "application/json", r.data)
		return true
	}
	return false
}