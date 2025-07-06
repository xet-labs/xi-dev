// cntr/blog.go
package cntr

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"xi/app/lib"
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
}

// Singleton controller
var BlogApi = &BlogCntr{
	db:    lib.DB.GetCli(),
	blog:  model.Blog{},
	blogs: []model.Blog{},
}

// GET /blog or /blog?Page=2&Limit=6
func (b *BlogCntr) Index(c *gin.Context) {
	page := c.DefaultQuery("Page", "1")
	limit := c.DefaultQuery("Limit", "6")

	pageNum, err1 := strconv.Atoi(page)
	limitNum, err2 := strconv.Atoi(limit)
	if err1 != nil || err2 != nil || pageNum <= 0 || limitNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page or Limit"})
		return
	}

	offset := (pageNum - 1) * limitNum

	// Optional: use paginated cache key
	redisKey := "blogs:page:" + page + ":limit:" + limit
	var blogs []model.Blog

	// Try cache
	if err := lib.Redis.GetJson(redisKey, &blogs); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"blogsExhausted": len(blogs) == 0,
			"blogs":          blogs,
		})
		return
	}

	// Cache miss or DB fallback
	err := b.db.Preload("User").
		Where("status IN ?", []string{"published", "published_hidden"}).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limitNum).
		Find(&blogs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blogsExhausted": len(blogs) == 0,
		"blogs":          blogs,
	})

	// Async Cache set
	go func(data []model.Blog) {
		if err := lib.Redis.SetJson(redisKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", redisKey, err)
		}
	}(blogs)

}

// GET api/blog/uid/id
func (b *BlogCntr) Show(c *gin.Context) {
	rawUID := c.Param("uid") // could be @username or UID
	rawID := c.Param("id")   // could be blogId or slug
	var blog model.Blog
	var redisKey string
	var err error

	if username, ok :=strings.CutPrefix(rawUID, "@"); ok  {
		redisKey = fmt.Sprintf("blogs:uname:%s:%s", username, rawID)

		// Try Redis
		if err = lib.Redis.GetJson(redisKey, &blog); err == nil {
			c.JSON(http.StatusOK, blog)
			return
		}

		// DB Fallback with JOIN on username
		err = b.db.Preload("User").
			Joins("JOIN users ON users.uid = blogs.uid").
			Where("users.username = ? AND (blogs.slug = ? OR blogs.id = ?)", username, rawID, rawID).
			First(&blog).Error

	} else if isNumeric(rawUID) {
		redisKey = fmt.Sprintf("blogs:uid:%s:%s", rawUID, rawID)

		if err = lib.Redis.GetJson(redisKey, &blog); err == nil {
			c.JSON(http.StatusOK, blog)
			return
		}

		// DB Fallback using UID directly
		err = b.db.Preload("User").
			Where("uid = ? AND (slug = ? OR id = ?)", rawUID, rawID, rawID).
			First(&blog).Error
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user identifier"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Return and cache
	c.JSON(http.StatusOK, blog)

	go func(key string, data model.Blog) {
		if err := lib.Redis.SetJson(key, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", key, err)
		}
	}(redisKey, blog)
}

// POST api/blog/uid/id
func (b *BlogCntr) Post(c *gin.Context) {
	var blog model.Blog

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	blog.CreatedAt = ptrTime(time.Now())

	if err := b.db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create blog"})
		return
	}

	// Invalidate blog list cache
	b.rdb.Del(lib.Redis.GetCtx(), "blogs:all")

	c.JSON(http.StatusCreated, blog)
}

// PUT api/blog/uid/id
func (b *BlogCntr) Put(c *gin.Context) {
	id := c.Param("id")
	var blog model.Blog

	if err := b.db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	blog.UpdatedAt = ptrTime(time.Now())

	if err := b.db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(lib.Redis.GetCtx(), "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, blog)
}

// DELETE api/blog/uid/id
func (b *BlogCntr) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := b.db.Delete(&model.Blog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(lib.Redis.GetCtx(), "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

// Utility: return pointer to time
func ptrTime(t time.Time) *time.Time {
	return &t
}
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
