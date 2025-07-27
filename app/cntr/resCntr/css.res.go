package resCntr

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"xi/app/lib"
	"xi/conf"
	"xi/util"

	"github.com/gin-gonic/gin"
)

var (
	CssDir      = conf.View.CssDir
	refKey    = "res:app.css"
	CssRedisTTL = 12 * time.Hour
	cssOnce     sync.Once
	cssFiles    []string
)

func init() {
	var err error
	cssFiles, err = util.GetExtFiles(".css", CssDir...)
	if err != nil {
		log.Println("Error loading CSS files:", err)
	}
	if len(cssFiles) == 0 {
		log.Println("Error no CSS files found")
	}
}

// mergeFiles reads and combines all CSS file content
func mergeFiles(files []string) string {
	var sb strings.Builder
	for _, file := range files {
		if data, err := os.ReadFile(file); err == nil {
			sb.Write(data)
		}
	}
	return sb.String()
}

// minifyCSS removes comments and compresses whitespace
func minifyCSS(css string) string {
	css = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(css, "")
	css = regexp.MustCompile(`\s*([{}:;,])\s*`).ReplaceAllString(css, "$1")
	css = strings.ReplaceAll(css, ";}", "}")
	css = regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	return strings.TrimSpace(css)
}

// Css handler: serves combined+minified CSS (Redis cached)
func Css(c *gin.Context) {
	c.Header("Content-Type", "text/css")

	// Try Redis cache
	if cached, err := lib.Redis.Get(refKey); err == nil && cached != "" {
		c.String(http.StatusOK, cached)
		return
	}

	css := minifyCSS(mergeFiles(cssFiles))
	c.String(http.StatusOK, css)

	// Store to Redis asynchronously
	go func(data string) {
		if err := lib.Redis.Set(refKey, data, CssRedisTTL); err != nil {
			log.Printf("Redis SET err (%s): %v", refKey, err)
		}
	}(css)

}
