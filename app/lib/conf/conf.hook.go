package conf

import (
	"errors"
	"maps"

	"github.com/rs/zerolog/log"
)

func (c *ConfLib) ConfPostView() {

	c.Hook.AddPost("ConfPostView - setup", c.ConfPostView)

	// Fetch defaults and pages
	pageDefault := c.GetMap("view.default")
	pages := c.GetMap("view.pages")

	rawJson := c.AllMap()
	if rawJson == nil {
		log.Error().Err(error(errors.New("no Json data to operate on"))).Msg("Config PostView")
		return
	}
	viewData, ok := rawJson["view"].(map[string]any)
	if !ok {
		log.Error().Err(error(errors.New("'view' is missing or not a map"))).Msg("Config PostView")
		return
	}

	// Ensure "pages" exists inside viewData
	viewPages, ok := viewData["pages"].(map[string]any)
	if !ok {
		// Create it if missing
		viewPages = make(map[string]any)
		viewData["pages"] = viewPages
	}

	// Merge defaults into each page
	for page, val := range pages {
		pageConf, ok := val.(map[string]any)
		if !ok {
			log.Warn().Str("page", page).Msg("Config Postview: Page data setup failed")
			continue
		}

		// Copy defaults first, then page-specific config
		rawConf := make(map[string]any)
		if pageDefault != nil {
			maps.Copy(rawConf, pageDefault)
		}
		maps.Copy(rawConf, pageConf)

		viewPages[page] = rawConf
	}

	// Save merged config
	if err := c.postProcess(rawJson); err != nil {
		log.Error().Err(err).Msg("Config Postview: failed to sync")
		return
	}
}
