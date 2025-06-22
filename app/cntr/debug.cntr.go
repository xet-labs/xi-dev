package cntr

import (
	"net/http"
	"xi/app/util"
	"github.com/gin-gonic/gin"
)

func D(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": util.Env("APP_NAME", "--"),
		"status":  "up",
	})
}