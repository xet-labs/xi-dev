package cntr

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidUID      = errors.New("invalid UID")
	ErrInvalidSlug     = errors.New("invalid slug")
	ErrBlogNotFound    = errors.New("blog not found")
)

func Rc(c *gin.Context, page string, refKey string) bool {
	if cached, err := lib.Redis.GetBytes(refKey); err == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", cached)
		return true
	}
	return false
}

func Rrc(c *gin.Context, page string, refKey string, P map[string]any) bool {
	// Render
	var buf bytes.Buffer
	if err := lib.View.Tcli.ExecuteTemplate(&buf, page, gin.H{"P": P}); err != nil {
		log.Printf("Render error for %s: %v", refKey, err)
		c.Status(http.StatusInternalServerError)
		return false
	}

	// Return and cache data
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
	go func(data []byte) {
		if err := lib.Redis.SetBytes(refKey, data, 10*time.Minute); err != nil {
			log.Printf("Redis SET err (%s): %v", refKey, err)
		}
	}(buf.Bytes())

	return true
}
