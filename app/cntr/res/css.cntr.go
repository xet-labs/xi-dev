package cntr

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Replace this with your logic to locate CSS files
func getCSSFiles() []string {
	dir := "../../asset/" // example path
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".css") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// Merge all CSS files into one string with banners
func mergeFiles(files []string) string {
	var combinedCss strings.Builder
	for _, file := range files {
		if content, err := os.ReadFile(file); err == nil {
			name := filepath.Base(file)
			combinedCss.WriteString(fmt.Sprintf("/* ------ [ STRT - %s ] ------ */\n", name))
			combinedCss.Write(content)
			combinedCss.WriteString(fmt.Sprintf("\n/* ------ [ ENDS - %s ] ------ */\n", name))
		}
	}
	return combinedCss.String()
}

// Minify CSS (optional - disabled for debug)
func minifyCSS(css string) string {
	return css // disabled for now (debug mode)

	// Enable with regex minification if needed
	// Use third-party minifiers for production
}

// Gin handler to serve merged CSS
func cssHandler(c *gin.Context) {
	c.Header("Content-Type", "text/css")

	files := getCSSFiles()
	if len(files) == 0 {
		c.String(http.StatusOK, "// nuii css'o fendoz")
		return
	}

	css := minifyCSS(mergeFiles(files))
	c.String(http.StatusOK, css)

	// Optional: write to file
	outputPath := "./public/app.css"
	os.WriteFile(outputPath, []byte(css), 0644)
}
