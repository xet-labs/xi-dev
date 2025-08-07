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

	mu         sync.RWMutex
	once       sync.Once
}

// Singleton controller
var Blog = &BlogCntr{
	db:    lib.Db.GetCli(),
	blog:  model.Blog{},
	blogs: []model.Blog{},
}

// GET /blog or /blog?Page=2&Limit=6
func (b *BlogCntr) Index(c *gin.Context) {
	// rdbKey := "/blog"
	rdbKey := c.Request.URL.String()

	// Try cache
	if lib.View.RenderCache(c, rdbKey).Html() { return }

	// Build data
	b.mu.Lock()
	defer b.mu.Unlock()
	P := cfg.View.Pages["blogs"]
	if P.Data == nil {
		P.Data = make(map[string]any)
	}
	P.Data["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, rdbKey, P)
}

func (b *BlogCntr) Show(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // blog ID or slug
	rdbKey := "/blog/" + rawUID + "/" + rawID
	
	if lib.View.RenderCache(c, rdbKey).Html() { return }

	blog := model.Blog{}
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
	b.mu.Lock()
	defer b.mu.Unlock()
	P := cfg.View.Pages["blog"]
	P.Data["B"] = blog
	P.Data["Content"] = template.HTML(blog.Content)
	P.Data["Url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, rdbKey, P)
}

// POST api/blog/uid/id
func (b *BlogCntr) Post(c *gin.Context) {}

// PUT api/blog/uid/id
func (b *BlogCntr) Put(c *gin.Context) {}

// DELETE api/blog/uid/id
func (b *BlogCntr) Delete(c *gin.Context) {}