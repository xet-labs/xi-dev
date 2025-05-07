// cntr/blog.go
package cntr

import (
	"net/http"
	"time"

	"xi/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlogCntr struct {
	blog  models.Blog
	blogs []models.Blog
	db    *gorm.DB
}

var Blog = &BlogCntr{
	blog:  models.Blog{},
	blogs: []models.Blog{},
	db:    (&models.Blog{}).DB().Debug(),
}

var db = Blog.db

// GET /blog
func (b *BlogCntr) Index(c *gin.Context) {
	if err := b.db.Preload("User").Find(&b.blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}
	c.JSON(http.StatusOK, b.blogs)
}

// GET /blog/:id
func (b *BlogCntr) Show(c *gin.Context) {

	var blog models.Blog

	if err := db.Joins("User").First(&blog, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}

// POST /blog
func (b *BlogCntr) Post(c *gin.Context) {

	if err := c.ShouldBindJSON(&b.blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	now := time.Now()
	b.blog.CreatedAt = &now

	if err := b.db.Create(&b.blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create blog"})
		return
	}
	c.JSON(http.StatusCreated, b.blog)
}

// PUT /blog/:id
func (b *BlogCntr) Put(c *gin.Context) {

	if err := b.db.First(&b.blog, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&b.blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	now := time.Now()
	b.blog.UpdatedAt = &now
	if err := b.db.Save(&b.blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update blog"})
		return
	}
	c.JSON(http.StatusOK, b.blog)
}

// DELETE /blog/:id
func (b *BlogCntr) Delete(c *gin.Context) {

	if err := b.db.Delete(&models.Blog{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}
