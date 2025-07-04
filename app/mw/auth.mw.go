package mw

import (
	"sync"

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type AuthMw struct {
	once sync.Once
	rw   sync.RWMutex
}

// Global singleton instance
var Auth = &AuthMw{}

func (a *AuthMw) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user == nil {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
