package res

import (
	"strings"
	"sync"
	"time"

	"xi/app/lib"
	"xi/app/lib/cfg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CssRes struct {
	data    map[string][]string // key: baseName, value: list of CSS file paths
	BaseDir string
	RdbTTL  time.Duration
	
	once sync.Once
	mu   sync.RWMutex
}

var Css = &CssRes{
	data:    make(map[string][]string),
	BaseDir: cfg.View.CssBaseDir,
	RdbTTL:  12 * time.Hour,
}

// Css handler: serves combined+cssMin CSS (Redis cached)
func (r *CssRes) Get(c *gin.Context) {
	rdbKey := c.Request.RequestURI
	base := cfg.View.CssBaseDir + "/" + strings.TrimSuffix(c.Param("name"), ".css")

	// Return cache
	if lib.View.OutCache(c, rdbKey).Css() {
		return
	}

	if _, ok := r.data[base]; !ok {
		var (
			files []string
			err   error
		)
		files, err = lib.File.GetWithExt(".css", base)
		if err != nil {
			log.Error().Err(err).Str("Dir", base).Msg("Controller CSS files")
			return
		}

		r.mu.Lock()
		r.data[base] = files
		r.mu.Unlock()
	}
	
	lib.View.OutCss(c, lib.File.MergeByte(r.data[base]), rdbKey)
}
