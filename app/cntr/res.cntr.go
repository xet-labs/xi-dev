package cntr

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Page": "home",
	})
}
