package cntr

import (
	"bytes"
	"log"
	"maps"
	"net/http"
	"os"
	"time"
	"xi/app/lib"
	"xi/conf"

	"github.com/gin-gonic/gin"
)

type PageCntr struct{}

var Page = &PageCntr{}

// Uses Go template directly to render a file-based page
func (p *PageCntr) Tmpl(title, tmpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := p.buildData(c, title)
		c.HTML(http.StatusOK, tmpl, gin.H{"P": data})
	}
}

// Renders raw HTML file inside a base layout and caches it
func (p *PageCntr) Tcnt(title, rawPath string, ttl ...time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		refKey := "page:" + c.Request.URL.String()

		// Return from cache if available
		if data, err := lib.Redis.GetBytes(refKey); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		rawContent, err := os.ReadFile(rawPath)
		if err != nil {
			log.Printf("Tcnt: error reading file: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		rendered, err := p.renderTcnt(c, title, string(rawContent))
		if err != nil {
			log.Printf("Tcnt: render error: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", rendered)

		// Async cache
		go func(data []byte) {
			expire := 10 * time.Minute
			if len(ttl) > 0 {
				expire = ttl[0]
			}
			if err := lib.Redis.SetBytes(refKey, data, expire); err != nil {
				log.Printf("Redis SET err (%s): %v", refKey, err)
			}
		}(rendered)
	}
}

// renderTcnt renders inline content inside base layout
func (p *PageCntr) renderTcnt(c *gin.Context, title, content string) ([]byte, error) {
	t, err := lib.View.RawTcli.Clone()
	if err != nil {
		return nil, err
	}

	if _, err := t.Parse(content); err != nil {
		return nil, err
	}

	P := p.buildData(c, title)
	// data, _ := json.MarshalIndent(P, "", "  ")
	// fmt.Println(string(data))
	
	var out bytes.Buffer
	if err := t.ExecuteTemplate(&out, conf.View.Layout, gin.H{"P": P}); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// Combines global and per-page config data
func (p *PageCntr) buildData(c *gin.Context, title string) map[string]any {
	data := make(map[string]any)
	maps.Copy(data, conf.View.PageData)

	if page, ok := conf.View.Pages[title]; ok {
		maps.Copy(data, page)
	}

	data["url"] = c.Request.URL.String()
	return data
}
