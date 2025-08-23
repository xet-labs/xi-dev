package view

import (
	"html/template"
	"sync"

	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
)

type ViewLib struct {
	Ecli      *gin.Engine        // Gin Engine
	Tcli      *template.Template // Current Template Cli
	RawTcli   *template.Template // Clean Template Cli
	templates []string

	once sync.Once
	mu   sync.RWMutex
}

var View = &ViewLib{
	templates: cfg.View.TemplateDirs,
}
