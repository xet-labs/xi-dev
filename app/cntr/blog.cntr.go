// cntr/blog.go
package cntr

import (
	"html/template"
	"net/http"
	"sync"
	"xi/app/lib"
	"xi/app/lib/cfg"
	"xi/app/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BlogCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
	blog  model.Blog
	blogs []model.Blog

	mu   sync.RWMutex
	once sync.Once
}

// Singleton controller
var Blog = &BlogCntr{
	db:    lib.Db.GetCli(),
	blog:  model.Blog{},
	blogs: []model.Blog{},
}

// GET /blog
// func (b *BlogCntr) Index(c *gin.Context) {}

func (b *BlogCntr) Show(c *gin.Context) {
	rdbKey := c.Request.RequestURI

	if lib.View.OutCache(c, rdbKey).Html() {
		return
	} // Try cache

	// On cache miss fetch data from DB
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // blog ID or slug

	blog := model.Blog{}
	if err := BlogApi.Validate(rawUID, rawID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b.mu.Lock()
	if err := BlogApi.ShowCore(&blog, rawUID, rawID); err != nil { // Fallback to DB
		status := http.StatusNotFound
		if err == ErrInvalidUID {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	b.mu.Unlock()

	p := cfg.View.Pages["blogs"]
	b.PrepMeta(c, &p.Meta, &blog)
	p.Rt = map[string]any{
		"B":       blog,
		"Content": template.HTML(blog.Content),
	}

	lib.View.OutHtmlLyt(c, p, rdbKey)
}

func (b *BlogCntr) PrepMeta(c *gin.Context, meta *model.PageMeta, raw *model.Blog) {
	meta.Type = "Article"
	meta.Title = raw.Title
	meta.URL = lib.Util.Url.Full(c)
	meta.AltJson = lib.Util.Url.Host(c) + "/api" + c.Request.RequestURI
	meta.Description = raw.Description
	meta.Img.URL = lib.Util.Url.Host(c) + raw.FeaturedImg
	meta.Tags = raw.Tags
	meta.Author.Name = raw.User.Name
	meta.Author.Img = raw.User.ProfileImg
	meta.Author.URL = lib.Util.Url.Host(c) + "/@" + raw.User.Username
	meta.CreatedAt = raw.CreatedAt
	meta.UpdatedAt = raw.UpdatedAt
	// meta.Category = raw.Tags
}

// POST api/blog/uid/id
func (b *BlogCntr) Post(c *gin.Context) {}

// PUT api/blog/uid/id
func (b *BlogCntr) Put(c *gin.Context) {}

// DELETE api/blog/uid/id
func (b *BlogCntr) Delete(c *gin.Context) {}
