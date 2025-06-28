// cntr/blog.go
package cntr

import (
	"log"
	"net/http"
	"time"

	"xi/app/model"
	"xi/app/lib"
	"xi/app/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BlogCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
	blog  models.Blog
	blogs []models.Blog
}

// Singleton controller
var Blog = &BlogCntr{
	db:    service.DB(),
	rdb:   lib.Redis.Cli(),
	blog:  models.Blog{},
	blogs: []models.Blog{},
}

// GET /blog
func (b *BlogCntr) Index(c *gin.Context) {
	var blogs []models.Blog
	redisKey := "blogs:all"

	// Try Redis cache
	if err := lib.Redis.GetJson(redisKey, &blogs); err == nil {
		c.JSON(http.StatusOK, blogs)
		return
	}

	// Cache miss or error - load from DB
	if err := Blog.db.Preload("User").Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}

	// Set cache asynchronously
	go func(data []models.Blog) {
		if err := lib.Redis.SetJson(redisKey, data, 10*time.Minute); err != nil {
			log.Println("Redis SET error (Index):", err)
		}
	}(blogs)

	c.JSON(http.StatusOK, blogs)
}

// GET /blog/:id
func (b *BlogCntr) Show(c *gin.Context) {
	id := c.Param("id")
	redisKey := "blogs:id:" + id
	var blog models.Blog

	// Try Redis
	if err := lib.Redis.GetJson(redisKey, &blog); err == nil {
		c.JSON(http.StatusOK, blog)
		return
	}

	// Cache miss - load from DB
	if err := Blog.db.Preload("User").First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Cache the result
	go func(data models.Blog) {
		if err := lib.Redis.SetJson(redisKey, data, 10*time.Minute); err != nil {
			log.Println("Redis SET error (Show):", err)
		}
	}(blog)

	c.JSON(http.StatusOK, blog)
}

// POST /blog
func (b *BlogCntr) Post(c *gin.Context) {
	var blog models.Blog

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	blog.CreatedAt = ptrTime(time.Now())

	if err := Blog.db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create blog"})
		return
	}

	// Invalidate blog list cache
	b.rdb.Del(lib.RedisCtx, "blogs:all")

	c.JSON(http.StatusCreated, blog)
}

// PUT /blog/:id
func (b *BlogCntr) Put(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog

	if err := Blog.db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	blog.UpdatedAt = ptrTime(time.Now())

	if err := Blog.db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(lib.RedisCtx, "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, blog)
}

// DELETE /blog/:id
func (b *BlogCntr) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := Blog.db.Delete(&models.Blog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(lib.RedisCtx, "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

// Utility: return pointer to time
func ptrTime(t time.Time) *time.Time {
	return &t
}
