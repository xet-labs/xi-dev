package view

import (
	"net/http"
	"xi/app/lib/db"

	"github.com/gin-gonic/gin"
)

type Render struct {
	c    *gin.Context
	data []byte
	ok   bool
	err  error
}

// Render from Cache
func (v *ViewLib) OutCache(c *gin.Context, rdbKey string) Render {
	cache, err := db.Rdb.GetBytes(rdbKey)
	return Render{c: c, data: cache, ok: err == nil, err: err}
}

// Helpers methods
func (r Render) Html() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/html; charset=utf-8", r.data)
		return true
	}
	return false
}

func (r Render) Css() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "text/css; charset=utf-8", r.data)
		return true
	}
	return false
}

func (r Render) Json() bool {
	if r.ok {
		r.c.Data(http.StatusOK, "application/json", r.data)
		return true
	}
	return false
}
