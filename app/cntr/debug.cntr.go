package cntr

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func Debug(r *gin.Engine) gin.HandlerFunc {

return func(c *gin.Context) {
	var routes []string
	var detailed []string

	for _, rt := range r.Routes() {
		method := rt.Method
		if len(method) > 3 {
			method = method[:3]
		}
		routes = append(routes, method+" "+rt.Path)
		detailed = append(detailed, method+" "+rt.Path+" | "+rt.Handler)
	}

	// Sort by route path (strip method prefix for sorting)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i][4:] < routes[j][4:]
	})
	sort.Slice(detailed, func(i, j int) bool {
		return detailed[i][6:] < detailed[j][6:] // "GET | " is 6 chars
	})

	c.JSON(http.StatusOK, gin.H{
		"route":         routes,
		"routeDetailed": detailed,
	})
}



}

