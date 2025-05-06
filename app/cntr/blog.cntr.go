// cntr/blog.go
package cntr

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"xi/app/models"
)

type BlogCntr struct {
	blog  	models.Blog
	blogs 	[]models.Blog
	db		*gorm.DB
}

var Blog = &BlogCntr{
	blog: 	models.Blog{},
	blogs: 	[]models.Blog{},
}

func init(){
	Blog.db=Blog.blog.DB()
}

// GET /blogs
func (b *BlogCntr) Index(c *gin.Context) {
	if err := b.db.Preload("User").Find(&b.blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blogs"})
		return
	}
	c.JSON(http.StatusOK, b.blogs)
}


// GET /blogs/:id
func (b *BlogCntr) Show(c *gin.Context) {
	
	if err := b.db.First(&b.blog, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	c.JSON(http.StatusOK, b.blog)
}

// POST /blogs
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

// PUT /blogs/:id
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

// DELETE /blogs/:id
func (b *BlogCntr) Delete(c *gin.Context) {
	
	if err := b.db.Delete(&models.Blog{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}
