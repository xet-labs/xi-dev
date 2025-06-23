// cntr/blog.go
package cntr

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"xi/app/model"
	"xi/app/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BlogCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
	ctx   context.Context
	blog  models.Blog
	blogs []models.Blog
}

// Singleton controller
var Blog = &BlogCntr{
	db:  service.DB(),
	rdb: service.Redis(),
	ctx: context.Background(),

	blog:  models.Blog{},
	blogs: []models.Blog{},
}

// Alias to reduce verbosity (optional)
var db = Blog.db

// GET /blog
func (b *BlogCntr) Index(c *gin.Context) {
	var blogs []models.Blog
	cacheKey := "blogs:all"
	
	// Try Redis cache
	if cached, err := b.rdb.Get(b.ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &blogs); err == nil {
			c.JSON(http.StatusOK, blogs)
			return
		}
	}

	// Cache miss or error - load from DB
	if err := db.Preload("User").Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}

	// Set cache asynchronously
	go func(data []models.Blog) {
		if jsonData, err := json.Marshal(data); err == nil {
			if err := b.rdb.Set(b.ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
				log.Println("Redis SET error (Index):", err)
			}
		}
	}(blogs)

	c.JSON(http.StatusOK, blogs)
}

// GET /blog/:id
func (b *BlogCntr) Show(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "blogs:id:" + id
	var blog models.Blog

	// Try Redis
	if cached, err := b.rdb.Get(b.ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &blog); err == nil {
			c.JSON(http.StatusOK, blog)
			return
		}
	}

	// Cache miss - load from DB
	if err := db.Preload("User").First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Cache the result
	go func(data models.Blog) {
		if jsonData, err := json.Marshal(data); err == nil {
			if err := b.rdb.Set(b.ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
				log.Println("Redis SET error (Show):", err)
			}
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

	if err := db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create blog"})
		return
	}

	// Invalidate blog list cache
	b.rdb.Del(b.ctx, "blogs:all")

	c.JSON(http.StatusCreated, blog)
}

// PUT /blog/:id
func (b *BlogCntr) Put(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog

	if err := db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	blog.UpdatedAt = ptrTime(time.Now())

	if err := db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(b.ctx, "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, blog)
}

// DELETE /blog/:id
func (b *BlogCntr) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&models.Blog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}

	// Invalidate caches
	b.rdb.Del(b.ctx, "blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

// Utility: return pointer to time
func ptrTime(t time.Time) *time.Time {
	return &t
}
