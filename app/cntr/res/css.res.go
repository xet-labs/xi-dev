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
	Files    map[string][]string // key: baseName, value: list of CSS file paths
	BaseDir string
	RdbTTL  time.Duration

	CacheFilePath bool
	once sync.Once
	mu   sync.RWMutex
}

var Css = &CssRes{
	Files:    make(map[string][]string),
	BaseDir: cfg.View.CssBaseDir,
	RdbTTL:  12 * time.Hour,

	CacheFilePath: false,
}

// Css handler: serves combined+cssMin CSS (Redis cached)
func (r *CssRes) Get(c *gin.Context) {
	rdbKey := c.Request.RequestURI
	base := cfg.View.CssBaseDir + "/" + strings.TrimSuffix(c.Param("name"), ".css")
	
	log.Info().Str("filepath", base).Msg("Css")

	if lib.View.OutCache(c, rdbKey).Css() {
		return 	// Send cache
	}

	if _, ok := r.Files[base]; !r.CacheFilePath || !ok {
		var (
			files []string
			err   error
		)
		files, err = lib.Util.File.GetWithExt(".css", base)
		if err != nil {
			log.Error().Err(err).Str("Dir", base).Msg("Controller CSS files")
			return
		}

		r.mu.Lock()
		r.Files[base] = files
		r.mu.Unlock()
	}

	lib.View.OutCss(c, lib.Util.File.MergeByte(r.Files[base]), rdbKey)
}
