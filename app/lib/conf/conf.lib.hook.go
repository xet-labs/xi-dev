package conf

import (
	"log"
	"maps"
	// "xi/app/cfg"
)

func (c *ConfLib) ConfPostView() {
	
	// c.Hook.AddPost("ConfPostView - setup", c.ConfPostView)

	pageDefault := c.GetMap("view.default")
	pages := c.GetMap("view.pages")
	
	rawJson := c.AllMap()
	viewPages, ok := rawJson["view"].(map[string]any)["pages"].(map[string]any)
	if !ok {
		log.Printf("⚠️  [Conf] Postview ERR: 'view.pages' is missing or not a valid map inside intermediate JSON")
		return
	}

	for page, val := range pages {
		pageConf, ok := val.(map[string]any)
		if !ok {
			log.Printf("⚠️  [Conf] Postview ERR: '%s' data setup failed", page)
			continue
		}

		// Deep copy viewDefault into rawConf to avoid mutation
		rawConf := make(map[string]any)
		maps.Copy(rawConf, pageDefault)
		maps.Copy(rawConf, pageConf)

		viewPages[page] = rawConf
	}

	c.postSetup(rawJson)

	// fmt.Printf("-->\n%s\n%s\n", cfg.Global, "")
	// fmt.Printf("--> Default:\n%v\nMap:\n%v\nJson:\n%s\n", pageDefault, c.AllMap(), c.AllJsonPretty())
	// fmt.Printf("--> Default:\n%v\nMap:\n%v\nJson:\n%s\n", pageDefault, c.AllMapStruct(), c.AllJsonStruct())
}