package lib

import (
	"log"
	"maps"
)

func (c *ConfLib) ConfPostView() {
	
	c.Hook.AddPost("ConfPostView - setup", c.ConfPostView)

	pageDefault := c.GetMap("view.default")
	pages := c.GetMap("view.pages")
	
	rawJson := c.AllMap()
	viewPages, ok := rawJson["view"].(map[string]any)["pages"].(map[string]any)
	if !ok {
		log.Printf("Config Err Post::view 'view.pages' is missing or not a valid map inside intermediate JSON")
		return
	}


	for page, val := range pages {
		// pageKey := fmt.Sprintf("view.pages.%s", page)
		pageConf, ok := val.(map[string]any)
		if !ok {
			log.Printf("Config Err Post::view::%s Page data clone failed", page)
			continue
		}

		// Deep copy viewDefault into rawConf to avoid mutation
		rawConf := make(map[string]any)
		maps.Copy(rawConf, pageDefault)
		maps.Copy(rawConf, pageConf)

		// Store into final merged output
		viewPages[page] = rawConf
	}

	c.DataJson = []byte(rawJson)
}