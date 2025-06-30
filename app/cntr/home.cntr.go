package cntr

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	PAGE := map[string]any{
		"title":       "Xet Blog",
		"description": "Cool dev stuff",
		"canonical":   "https://xetlabs.com/blog",
		"tags":        []string{"XetIndustries", "Go", "Gin"},
		"ogType":      "article",
		"metaTitle":   "Xet Blog | XetIndustries",
		"url": "",
		"featured_img": []string{
			"/res/static/brand/brand.svg",
		},
		"x": map[string]any{
			"site":     "@xet",
			"creator":  "@xet-author",
			"category": "Dev",
		},
	}

	c.HTML(http.StatusOK, "pages/home", gin.H{
		"PAGE": PAGE,
	})

}
