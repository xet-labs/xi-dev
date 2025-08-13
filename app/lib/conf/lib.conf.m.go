package conf

import (
	"encoding/json"
	"log"
	"xi/app/lib/cfg"
	"xi/app/model"
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
	var cfgRaw map[string]any
	if err := c.koanfCli.Unmarshal("", &cfgRaw); err != nil {
		log.Printf("⚠️  [Conf] AllMap WRN: failed to unmarshal: %v", err)
		return map[string]any{} // Return an empty map on error
	}
	return cfgRaw
}

func (c *ConfLib) AllJson() []byte {
	out, err := json.Marshal(c.AllMap())
	if err != nil {
		log.Printf("⚠️  [Conf] AllJson WRN: failed to marshal: %v", err)
		return []byte("{}")
	}
	return out
}

func (c *ConfLib) AllJsonPretty() []byte {
	out, err := json.MarshalIndent(c.AllMap(), "", "  ")
	if err != nil {
		log.Printf("⚠️  [Conf] AllJson WRN: failed to marshal: %v", err)
		return []byte("{}")
	}
	return out
}

func (c *ConfLib) AllMapStruct() *model.Config {
	cfgRaw := cfg.Get()
	// if err := c.koanfCli.Unmarshal("", &cfgRaw); err != nil {
	// 	log.Printf("⚠️  [Conf] AllMapStruct WRN: failed to unmarshal: %v", err)
	// 	return model.Config{} // Return zero value on error
	// }
	if err := json.Unmarshal(c.AllJson(), &cfgRaw); err != nil {
		log.Printf("⚠️  [Conf] AllMapStruct WRN: failed to unmarshal: %v", err)
		return &model.Config{}
	}
	return cfgRaw
}

func (c *ConfLib) AllJsonStruct() []byte {
	out, err := json.Marshal(c.AllMapStruct())
	if err != nil {
		log.Printf("⚠️  [Conf] AllJson WRN: failed to marshal: %v", err)
		return []byte("{}")
	}
	return out
}

func (c *ConfLib) AllJsonStructPretty() []byte {
	out, err := json.MarshalIndent(c.AllMapStruct(), "", "  ")
	if err != nil {
		log.Printf("⚠️  [Conf] AllJson WRN: failed to marshal: %v", err)
		return []byte("{}")
	}
	return out
}
