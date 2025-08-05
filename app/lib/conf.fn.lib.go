package lib

import (
	"encoding/json"
	"log"
	"maps"
)

func (c *ConfLib) Get(path string) any {
	return c.koanfCli.Get(path)
}

func (c *ConfLib) GetStr(path string) string {
	return c.koanfCli.Get(path).(string)
}

func (c *ConfLib) GetMap(path string) map[string]any {
	if val, ok := c.koanfCli.Get(path).(map[string]any); ok {
		return val
	}
	return map[string]any{}
}

func (c *ConfLib) GetArr(path string) []any {
	if val, ok := c.koanfCli.Get(path).([]any); ok {
		return val
	}
	return []any{}
}

func (c *ConfLib) All() map[string]any {
	return c.koanfCli.All()
}

func (c *ConfLib) AllMap() map[string]any {
	var nested map[string]any
	if err := c.koanfCli.Unmarshal("", &nested); err != nil {
		log.Printf("⚠️ Config failed to unmarshal: %v", err)
		return map[string]any{} // Return an empty map on error
	}
	return nested
}

func (c *ConfLib) AllJSON() []byte {
	if out, err := json.MarshalIndent(c.AllMap(), "", "  "); err == nil {
		log.Printf("⚠️ COnfig failed to marshal: %v", err)
		return out
	}
	return []byte("{}")
}

func (c *ConfLib) PostConf() {
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

	c.DataJson = rawJson
}
