package res

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"xi/app/lib"
	"xi/util"

	"github.com/gin-gonic/gin"
)

var (
	CssDir      = "views/partials"
	redisKey    = "res:app.css"
	CssRedisTTL = 12 * time.Hour
	cssOnce     sync.Once
	cssFiles    []string
)

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
	if cached, err := lib.Redis.Get(redisKey); err == nil && cached != "" {
		c.String(http.StatusOK, cached)
		return
	}

	// Cache miss: load + generate
	cssOnce.Do(func() {
		cssFiles, _ = util.GetFilesWithExt(CssDir, ".css")
	})

	if len(cssFiles) == 0 {
		c.String(http.StatusOK, "// no css files found")
		return
	}

	css := minifyCSS(mergeFiles(cssFiles))
	c.String(http.StatusOK, css)

	// Store to Redis asynchronously
	go func(data string) {
		if err := lib.Redis.Set(redisKey, data, CssRedisTTL); err != nil {
			log.Printf("Redis SET err (%s): %v", redisKey, err)
		}
	}(css)

}
