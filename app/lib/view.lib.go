package lib

import (
	"html/template"
	"sync"

	"github.com/gin-gonic/gin"
)

type ViewLib struct {
	Ecli *gin.Engine
	Tcli *template.Template
	RawTcli  *template.Template

	rw   sync.RWMutex
	once sync.Once
	// lazyInit func()
}

// Global instance
var View = &ViewLib{}
