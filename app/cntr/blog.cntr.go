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
func (b *BlogCntr) Index(c *gin.Context) {
	rdbKey := "/blog"

	if lib.View.OutCache(c, rdbKey).Html() {
		return
	} // Try cache

	b.mu.RLock()
	defer b.mu.RUnlock()

	// Build data
	P := cfg.View.Pages["blogs"]
	P.Meta.URL = lib.Url.Full(c)

	lib.View.OutHtmlLyt(c, P, rdbKey) // Cache renderer
}

func (b *BlogCntr) Show(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // blog ID or slug
	rdbKey := "/blog/" + rawUID + "/" + rawID

	if lib.View.OutCache(c, rdbKey).Html() {
		return
	} // Try cache

	b.mu.Lock()
	defer b.mu.Unlock()

	blog := model.Blog{}
	if err := BlogApi.Validate(rawUID, rawID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := BlogApi.ShowCore(&blog, rawUID, rawID); err != nil { // Fallback to DB
		status := http.StatusNotFound
		if err == ErrInvalidUID {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	P := cfg.View.Pages["blog"]
	b.PrepMeta(&P.Meta, blog, c)
	P.Rt = map[string]any{
		"B":       blog,
		"Content": template.HTML(blog.Content),
	}

	lib.View.OutHtmlLyt(c, P, rdbKey)
}

func (b *BlogCntr) PrepMeta(meta *model.PageMeta, raw model.Blog, c *gin.Context){
	meta.Type = "Article" 
	meta.Title = raw.Title
	meta.URL = lib.Url.Full(c)
	meta.AltJson = lib.Url.Host(c) + "/api" + c.Request.RequestURI
	meta.Description = raw.Description
	meta.Img.URL = lib.Url.Host(c) + raw.FeaturedImg
	meta.Tags = raw.Tags
	meta.Author.Name = raw.User.Name
	meta.Author.Img = raw.User.ProfileImg
	meta.Author.URL = lib.Url.Host(c) + "/@" + raw.User.Username
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
