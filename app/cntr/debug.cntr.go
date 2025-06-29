package cntr

import (
	"net/http"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

func D(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":   lib.Env("APP_NAME", "--"),
		"status": "up",
	})
}
