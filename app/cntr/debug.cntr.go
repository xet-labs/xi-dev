package cntr

import (
	"net/http"
	"xi/app/utils"
	"github.com/gin-gonic/gin"
)

func D(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": utils.Env("APP_NAME", "--"),
		"status":  "up",
	})
}