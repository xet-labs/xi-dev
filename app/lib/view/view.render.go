package view

import (
	"net/http"
	"time"
	"xi/app/lib/db"
	"xi/app/lib/minify"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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
