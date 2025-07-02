package cntr

import (
	"net/http"
	"xi/data"
	"xi/util"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	var PAGE = map[string]any{
		"url":c.Request.URL.String(),
		"title":       "XetIndustries",
		"currentMenu": "Home",
		"excerpt":     "XetIndustries is a Collaborative Platform for Makers, Creators, and Developers.",
		"canonical":   "https://xetindustries.com/",
		"tags":        []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
	}
	util.MergeMapTo(PAGE, data.View)

	c.HTML(http.StatusOK, "pages/home", gin.H{
		"P": PAGE,
	})

}
