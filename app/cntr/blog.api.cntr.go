// cntr/blog.api.go
package cntr

import (
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

type BlogApiCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
	blog  model.Blog
	blogs []model.Blog
}

// Singleton controller
var BlogApi = &BlogApiCntr{
	db:    lib.Db.GetCli(),
	blog:  model.Blog{},
	blogs: []model.Blog{},
}

// GET /blog or /blog?Page=2&Limit=6
func (b *BlogApiCntr) Index(c *gin.Context) {
	page := c.DefaultQuery("Page", "1")
	limit := c.DefaultQuery("Limit", "6")
	pageNum, err1 := strconv.Atoi(page)
	limitNum, err2 := strconv.Atoi(limit)

	if err1 != nil || err2 != nil || pageNum <= 0 || limitNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page or Limit"})
		return
	}

	// Try cache
	blogs := []model.Blog{}
	rdbKey := "/api/blog/?Page=" + page + "&Limit=" + limit
	if err := lib.Rdb.GetJson(rdbKey, &blogs); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"blogsExhausted": len(blogs) == 0,
			"blogs":          blogs,
		})
		return
	}

	// Fallback to DB
	offset := (pageNum - 1) * limitNum
	if err := b.IndexCore(&blogs, offset, limitNum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blogsExhausted": len(blogs) == 0,
		"blogs":          blogs,
	})

	// Async cache
	go func(data any) { lib.Rdb.SetJson(rdbKey, data, 10*time.Minute) }(blogs)
}

func (b *BlogApiCntr) IndexCore(blogs *[]model.Blog, offset, limit int) error {
	return b.db.
		Preload("User").
		Where("status IN ?", []string{"published", "published_hidden"}).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(blogs).
		Error
}

// GET api/blog/uid/id
func (b *BlogApiCntr) Show(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // blog ID or slug

	// Try cache
	blog := model.Blog{}
	rdbKey := "/api/blog/" + rawUID + "/" + rawID
	if err := lib.Rdb.GetJson(rdbKey, &blog); err == nil {
		c.JSON(http.StatusOK, blog)
		return
	}

	// Validate parameters
	if err := b.Validate(rawUID, rawID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fallback to DB
	if err := b.ShowCore(&blog, rawUID, rawID); err != nil {
		status := http.StatusNotFound
		if err == ErrInvalidUID {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)

	// Cache asynchronously
	go func(data model.Blog) { lib.Rdb.SetJson(rdbKey, data, 10*time.Minute) }(blog)
}

// FetchBlog fetches a blog and stores it in the given pointer.
// It returns the Redis key and any error.
func (b *BlogApiCntr) ShowCore(dest *model.Blog, rawUID, rawID string) error {
	var err error

	// Case 1: @username format
	if username, ok := strings.CutPrefix(rawUID, "@"); ok {
		// DB fallback
		err = b.db.Preload("User").
			Joins("JOIN users ON users.uid = blogs.uid").
			Where("users.username = ? AND (blogs.slug = ? OR blogs.id = ?)", username, rawID, rawID).
			First(dest).Error

		// Case 2: UID (numeric)
	} else if isNumeric(rawUID) {

		err = b.db.Preload("User").
			Where("uid = ? AND (slug = ? OR id = ?)", rawUID, rawID, rawID).
			First(dest).Error

		// Invalid UID format
	} else {
		return ErrInvalidUID
	}

	if err != nil {
		return ErrBlogNotFound
	}

	return nil
}

// POST api/blog/uid/id
func (b *BlogApiCntr) Post(c *gin.Context) {
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

	c.JSON(http.StatusCreated, blog)

	// Invalidate blog list cache
	lib.Rdb.Del("blogs:all")
}

// PUT api/blog/uid/id
func (b *BlogApiCntr) Put(c *gin.Context) {
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

	c.JSON(http.StatusOK, blog)

	// Invalidate caches
	lib.Rdb.Del("blogs:all", "blogs:id:"+id)
}

// DELETE api/blog/uid/id
func (b *BlogApiCntr) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := b.db.Delete(&model.Blog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}

	// Invalidate caches
	lib.Rdb.Del("blogs:all", "blogs:id:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

// --HELPERS--
// Utility: return pointer to time
func ptrTime(t time.Time) *time.Time {
	return &t
}
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (b *BlogApiCntr) Validate(rawUID, rawID string) error {
	if strings.HasPrefix(rawUID, "@") {
		if !lib.Validate.Uname(rawUID) {
			return ErrInvalidUserName
		}
	} else if !lib.Validate.UID(rawUID) {
		return ErrInvalidUID
	}
	if !lib.Validate.Slug(rawID) {
		return ErrInvalidSlug
	}
	return nil
}
