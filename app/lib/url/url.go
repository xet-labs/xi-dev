package url

import "github.com/gin-gonic/gin"

type UrlLib struct {}

var Url = & UrlLib{}

func (u *UrlLib) Full(c *gin.Context) string {
    scheme := "http"
    if c.Request.TLS != nil {
        scheme = "https"
    }

    // c.Request.Host includes host + port
    return scheme + "://" + c.Request.Host + c.Request.RequestURI
}
