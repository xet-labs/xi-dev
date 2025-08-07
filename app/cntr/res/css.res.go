package res

import (
	"log"
	"time"

	"xi/app/lib"
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
)

var (
	CssDir      = cfg.View.CssDir
	CssRedisTTL = 12 * time.Hour
	cssFiles    []string
)

func init() {
	var err error
	cssFiles, err = lib.File.GetExt(".css", CssDir...)
	if err != nil {
		log.Println("CSS Err loading files:", err)
	}
	if len(cssFiles) == 0 {
		log.Println("CSS Err no files found")
	}
}

// Css handler: serves combined+cssMin CSS (Redis cached)
func Css(c *gin.Context) {
	c.Header("Content-Type", "text/css")
	rdbKey := c.Request.URL.String()

	if lib.View.OutCache(c, rdbKey).Css() {
		return
	}

	lib.View.OutCss(c, lib.File.MergeByte(cssFiles))
}
