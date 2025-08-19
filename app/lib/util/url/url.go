package url

import "github.com/gin-gonic/gin"

type UrlLib struct {}

func (u *UrlLib) Host(c *gin.Context) string {
    scheme := "http"
    if c.Request.TLS != nil {
        scheme = "https"
    }

    // c.Request.Host includes host + port
    return scheme + "://" + c.Request.Host
}

func (u *UrlLib) Full(c *gin.Context) string {

    // c.Request.Host includes host + port
    return u.Host(c) + c.Request.RequestURI
}
