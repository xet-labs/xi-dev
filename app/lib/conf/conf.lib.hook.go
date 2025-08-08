package conf

import (
	"log"
	"maps"
	// "xi/app/cfg"
)

func (c *ConfLib) ConfPostView() {

	// c.Hook.AddPost("ConfPostView - setup", c.ConfPostView)

	// Fetch defaults and pages
	pageDefault := c.GetMap("view.default")
	pages := c.GetMap("view.pages")

	rawJson := c.AllMap()
	if rawJson == nil {
		log.Printf("⚠️  [Conf] Postview ERR: No Json data to operate on")
		return
	}
	viewData, ok := rawJson["view"].(map[string]any)
	if !ok {
		log.Println("⚠️  [Conf] Postview ERR: 'view' is missing or not a map")
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
			log.Printf("⚠️  [Conf] Postview ERR: '%s' data setup failed", page)
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
	c.postSetup(rawJson)
	
}
