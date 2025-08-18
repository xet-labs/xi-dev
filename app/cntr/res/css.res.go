package res

import (
	"strings"
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
	data map[string][]string // key: baseName, value: list of CSS file paths
}

var Css = &CssRes{
	RdbTTL: 12 * time.Hour,
}


// func (ac *AssetCache) GetOrBuild(baseName string) ([]string, error) {
// 	ac.mu.RLock()
// 	files, ok := ac.data[baseName]
// 	ac.mu.RUnlock()
// 	if ok {
// 		return files, nil
// 	}

// 	// Build the list of files from partials/baseName/
// 	dir := filepath.Join("partials", baseName)
// 	var found []string
// 	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !d.IsDir() && strings.HasSuffix(path, ".css") {
// 			found = append(found, path)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	ac.mu.Lock()
// 	ac.data[baseName] = found
// 	ac.mu.Unlock()
// 	return found, nil
// }

// Css handler: serves combined+cssMin CSS (Redis cached)
func (r *CssRes) Get(c *gin.Context) {
	base := strings.TrimSuffix(c.Param("name"), ".css")

	if r.Files == nil {
		r.getFiles()
	}
	rdbKey := c.Request.URL.String()

	if lib.View.OutCache(c, rdbKey).Css() {
		return
	}

	lib.View.OutCss(c, lib.File.MergeByte(r.Files), rdbKey)
}

func (c *CssRes) getFiles() {
	var err error
	Css.Files, err = lib.File.GetWithExt(".css", cfg.View.CssDir...)
	if err != nil {
		log.Error().Err(err).Msg("Controller CSS files")
	}
	if len(Css.Files) == 0 {
		log.Warn().Msgf("Controller CSS no files found in directory: %s", lib.Util.QuoteSlice(cfg.View.CssDir))
	}
}
