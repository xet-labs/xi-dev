package res

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

var (
	CssDir      = "views/partial/"
	CssRedisKey = "res:app.css"
	CssRedisTTL = 12 * time.Hour
	filesOnce   sync.Once
	cssFiles    []string
)

// loadCSSFiles caches the list of CSS files on first use
func loadCSSFiles() {
	_ = filepath.Walk(CssDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".css") {
			cssFiles = append(cssFiles, path)
		}
		return nil
	})

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
	if cached, err := lib.Redis.GetString(CssRedisKey); err == nil && cached != "" {
		c.String(http.StatusOK, cached)
		return
	}

	// Cache miss: load + generate
	filesOnce.Do(func() {
		loadCSSFiles()
	})

	if len(cssFiles) == 0 {
		c.String(http.StatusOK, "// no css files found")
		return
	}

	css := minifyCSS(mergeFiles(cssFiles))
	c.String(http.StatusOK, css)

	// Store to Redis asynchronously
	go func(content string) {
		_ = lib.Redis.SetString(CssRedisKey, content, CssRedisTTL)
	}(css)
}
