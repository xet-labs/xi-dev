package view

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"xi/app/lib/cfg"
)

type ViewLib struct {
	Ecli      *gin.Engine        // Gin Engine
	Tcli      *template.Template // Current Template Cli
	RawTcli   *template.Template // Clean Template Cli
	templates []string
}

var View = &ViewLib{
	templates: cfg.View.TemplateDir,
}
