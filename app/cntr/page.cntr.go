package cntr

import (
	"bytes"
	"html/template"
	"log"
	mapUtil "maps"
	"net/http"
	"os"
	"time"

	"xi/app/lib"
	"xi/conf"

	// "xi/util"

	"github.com/gin-gonic/gin"
)

// var db = lib.DB.GetCli()

func Page(tmpl *template.Template, title, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisKey := "page:" + c.Request.URL.String()

		// Try Redis cache
		if data, err := lib.Redis.GetBytes(redisKey); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		rawCnt, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
		}

		// Prep page data
		page := make(map[string]any)
		mapUtil.Copy(page, conf.View.PageData)
		mapUtil.Copy(page, conf.View.Pages[title])
		page["url"] = c.Request.URL.String()

		// Parse the content into the cloned template (defines "content")
		t, _ := tmpl.Clone()
		if _, err := t.Parse(string(rawCnt)); err != nil {
			log.Fatal(err)
		}

		// Execute layout/tmpl using the merged data
		var cnt bytes.Buffer
		err = t.ExecuteTemplate(&cnt, "layout/base", gin.H{
			"P": page,
		})
		if err != nil {
			log.Fatal(err)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", cnt.Bytes())

		// Cache the result
		go func(data []byte) {
			if err := lib.Redis.SetBytes(redisKey, data, 10*time.Minute); err != nil {
				log.Printf("Redis SET err (%s): %v", redisKey, err)
			}
		}(cnt.Bytes())
	}
}

func PageTmpl(title, path string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Prep page data
		page := make(map[string]any)
		mapUtil.Copy(page, conf.View.PageData)
		mapUtil.Copy(page, conf.View.Pages[title])
		page["url"] = c.Request.URL.String()

		c.HTML(http.StatusOK, path, gin.H{
			"P": page,
		})
	}
}
