package conf

import (
	"errors"
	"maps"

	"github.com/rs/zerolog/log"
)

func (c *ConfLib) ConfPostView() error {

	// c.Hook.AddPost("ConfPostView - setup", c.ConfPostView)

	// Fetch defaults and pages
	pageDefault := c.GetMap("view.default")
	pages := c.GetMap("view.pages")

	rawJson := c.AllMap()
	if rawJson == nil {
		err := "Config Postview: No Json data to operate on"
		log.Error().Msg(err)
		return errors.New(err)
	}
	viewData, ok := rawJson["view"].(map[string]any)
	if !ok {
		err := "Config Postview: 'view' is missing or not a map"
		log.Error().Msg(err)
		return errors.New(err)
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
			log.Warn().Msgf("Config Postview: '%s' data setup failed", page)
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
	if err := c.postSetup(rawJson); err != nil {
		log.Error().Msgf("Config Postview: failed to sync: %v", err)
		return err
	}
	return nil
}
