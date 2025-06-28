package res

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"xi/lib"
	// "xi/app/util"
)

var (
	cssCache string
	once     sync.Once
	redisKey = lib.Redis.Key("res:app.css")
	cacheTTL = time.Hour * 12
)

// getCSSFiles recursively finds all .css files in the partials folder
func getCSSFiles() []string {
	dir := "views/partial/"
	var files []string

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".css") {
			files = append(files, path)
		}
		return nil
	})

	return files
}

// mergeFiles reads and combines all CSS files with start/end banners
func mergeFiles(files []string) string {
	var sb strings.Builder

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		// name := filepath.Base(file)
		// sb.WriteString(fmt.Sprintf("/* ------ [ STRT - %s ] ------ */\n", name))
		sb.Write(data)
		// sb.WriteString(fmt.Sprintf("\n/* ------ [ ENDS - %s ] ------ */\n\n", name))
	}
	return sb.String()
}

// minifyCSS removes comments, newlines, and excessive spaces
func minifyCSS(css string) string {
	// Remove /* comments */
	reComments := regexp.MustCompile(`(?s)/\*.*?\*/`)
	css = reComments.ReplaceAllString(css, "")

	// Remove spaces around selectors, braces, colons, semicolons
	reSpaces := regexp.MustCompile(`\s*([{}:;,])\s*`)
	css = reSpaces.ReplaceAllString(css, "$1")

	// Remove unnecessary semicolons before }
	css = strings.ReplaceAll(css, ";}", "}")

	// Collapse multiple spaces and remove newlines
	reMultiSpaces := regexp.MustCompile(`\s+`)
	css = reMultiSpaces.ReplaceAllString(css, " ")

	return strings.TrimSpace(css)
}

// Css handler: merges and serves combined+minified CSS
func Css(c *gin.Context) {
	c.Header("Content-Type", "text/css")

	files := getCSSFiles()
	if len(files) == 0 {
		c.String(http.StatusOK, "// no css files found")
		return
	}

	finalCSS := minifyCSS(mergeFiles(files))

	c.String(http.StatusOK, finalCSS)

	// Optional: write minified CSS to disk
	// _ = os.WriteFile("public/app.css", []byte(finalCSS), 0644)
}
