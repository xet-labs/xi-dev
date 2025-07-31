// cntr/blog.go
package cntr

import (
	"html/template"
	"maps"
	"net/http"
	"net/url"
	"xi/app/lib"
	"xi/app/model"
	"xi/conf"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BlogCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
	blog  model.Blog
	blogs []model.Blog
}

// Singleton controller
var Blog = &BlogCntr{
	db:    lib.DB.GetCli(),
	blog:  model.Blog{},
	blogs: []model.Blog{},
}

// GET /blog or /blog?Page=2&Limit=6
func (b *BlogCntr) Index(c *gin.Context) {
	refKey := "/blog"

	// Try cache
	if lib.View.RenderCache(c, "layout/blog", refKey) {
		return
	}

	// Build data
	P := make(map[string]any)
	maps.Copy(P, conf.View.PageData)
	if pd, ok := conf.View.Pages["blogs"]; ok {
		maps.Copy(P, pd)
	}
	P["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, "layout/blogs", refKey, P)
}

func (b *BlogCntr) Show(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // blog ID or slug
	refKey := "/blog/" + url.QueryEscape(rawUID) + "/" + url.QueryEscape(rawID)
	var blog model.Blog

	// Try cache
	if lib.View.RenderCache(c, "layout/blog", refKey) {
		return
	}

	// Validate params
	if err := BlogApi.Validate(rawUID, rawID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fallback to DB
	if err := BlogApi.ShowCore(&blog, rawUID, rawID); err != nil {
		status := http.StatusNotFound
		if err == ErrInvalidUID {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prep data
	P := make(map[string]any)
	maps.Copy(P, conf.View.PageData)
	if pd, ok := conf.View.Pages["blog"]; ok {
		maps.Copy(P, pd)
	}
	P["B"] = BlogView{
		Blog:    blog,
		Content: template.HTML(blog.Content),
	}
	P["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, "layout/blog", refKey, P)
}

// POST api/blog/uid/id
func (b *BlogCntr) Post(c *gin.Context) {}

// PUT api/blog/uid/id
func (b *BlogCntr) Put(c *gin.Context) {}

// DELETE api/blog/uid/id
func (b *BlogCntr) Delete(c *gin.Context) {}

// --HELPERS--
type BlogView struct {
	model.Blog
	Content template.HTML
}
