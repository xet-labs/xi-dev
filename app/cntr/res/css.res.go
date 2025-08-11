package res

import (
	"time"

	"xi/app/lib"
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CssRes struct {
	Dir    []string
	RdbTTL time.Duration
	Files  []string
}

var Css = &CssRes{
	RdbTTL: 12 * time.Hour,
}

// Css handler: serves combined+cssMin CSS (Redis cached)
func (c *CssRes) Get(g *gin.Context) {
	if c.Files == nil {
		c.getFiles()
	}
	rdbKey := g.Request.URL.String()

	if lib.View.OutCache(g, rdbKey).Css() {
		return
	}

	lib.View.OutCss(g, lib.File.MergeByte(c.Files), rdbKey)
}


func (c *CssRes) getFiles() {
	var err error
	Css.Files, err = lib.File.GetExt(".css", cfg.View.CssDir...)
	if err != nil {
		log.Error().Err(err).Msg("Controller CSS files")
	}
	if len(Css.Files) == 0 {
		log.Warn().Msgf("Controller CSS no files found in directory: %s", lib.Util.QuoteSlice(cfg.View.CssDir))
	}
}