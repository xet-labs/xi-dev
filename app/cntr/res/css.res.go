package res

import (
	"log"
	"net/http"
	"time"

	"xi/app/lib"
	"xi/conf"
	"xi/util"

	"github.com/gin-gonic/gin"
)

var (
	CssDir      = conf.View.CssDir
	refKey      = "res:app.css"
	CssRedisTTL = 12 * time.Hour
	cssFiles    []string
)

func init() {
	var err error
	cssFiles, err = util.GetExtFiles(".css", CssDir...)
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

	// Try Redis cache
	if cached, err := lib.Redis.Get(refKey); err == nil && cached != "" {
		c.String(http.StatusOK, cached)
		return
	}

	// Merge and minify CSS
	cssCnt := lib.Merge.FilesByte(cssFiles)
	cssMin, err := lib.Minify.CssHybrid(cssCnt)
	if err != nil {
		log.Printf("CSS Minify Err: %v", err)
		c.Data(http.StatusOK, "text/css; charset=utf-8", cssCnt)
		return
	}

	// Serve response
	c.Data(http.StatusOK, "text/css; charset=utf-8", cssMin)


	// Cache it in background
	go func(data any) {
		if err := lib.Redis.Set(refKey, data, CssRedisTTL); err != nil {
			log.Printf("Redis SET err (%s): %v", refKey, err)
		}
	}(cssMin)
}
