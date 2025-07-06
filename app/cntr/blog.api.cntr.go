// cntr/blog.go
package cntr

import (
	"log"
	"net/http"
	"strconv"
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
	// Default pagination values
	q_page := c.DefaultQuery("Page", "1")
	q_limit := c.DefaultQuery("Limit", "6")

	pageNum, err1 := strconv.Atoi(q_page)
	limitNum, err2 := strconv.Atoi(q_limit)
	if err1 != nil || err2 != nil || pageNum <= 0 || limitNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page or Limit"})
		return
	}

	offset := (pageNum - 1) * limitNum

	// Optional: use paginated cache key
	redisKey := "blogs:page:" + q_page + ":limit:" + q_limit
	var blogs []model.Blog

	// Try Redis cache
	if err := lib.Redis.GetJson(redisKey, &blogs); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"blogsExhausted": len(blogs) == 0,
			"blogs":       blogs,
		})
		return
	}

	// Cache miss or DB fallback
	if err := b.db.Preload("User").
		Where("status IN ?", []string{"published", "published_hidden"}).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limitNum).
		Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blogsExhausted": len(blogs) == 0,
		"blogs":       blogs,
	})

	// Async Redis set
	go func(data []model.Blog) {
		if err := lib.Redis.SetJson(redisKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", redisKey, err)
		}
	}(blogs)

}


// GET /blog/:id
func (b *BlogCntr) Show(c *gin.Context) {
	id := c.Param("id")
	redisKey := "blogs:id:" + id
	var blog model.Blog

	// Try Redis
	if err := lib.Redis.GetJson(redisKey, &blog); err == nil {
		c.JSON(http.StatusOK, blog)
		return
	}

	// Cache miss - load from DB
	if err := b.db.Preload("User").First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)

	// Cache the result
	go func(data model.Blog) {
		if err := lib.Redis.SetJson(redisKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", redisKey, err)
		}
	}(blog)
}

// POST /blog
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

// PUT /blog/:id
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

// DELETE /blog/:id
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
