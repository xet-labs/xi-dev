package cntr

import (
	"net/http"
	"xi/data"
	"xi/util"

	"github.com/gin-gonic/gin"
)

func Page(page string) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var PAGE = map[string]any{
			"url":         c.Request.URL.String(),
			"title":       "XetIndustries",
			"currentMenu": "Home",
			"description": "XetIndustries is a Collaborative Platform for Makers, Creators, and Developers.",
			"canonical":   "https://xetindustries.com/",
			"tags":        []string{"XetIndustries", "Xet Industries", "xetindustries", "xet industries", "Xtreme Embedded Tech Industries"},
		}
		
		util.MergeMapTo(PAGE, data.View)

		c.HTML(http.StatusOK, page, gin.H{
			"P": PAGE,
            // any other data you want to pass
        })
    }
}
