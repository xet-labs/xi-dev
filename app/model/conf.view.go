package model

type ViewConf struct {
	CssDir       string                `json:"css_dir,omitempty"`
	CssDirs      []string              `json:"css_dirs,omitempty"`
	CssBaseDir   string                `json:"css_base_dir,omitempty"`
	TemplateDir  string                `json:"template_dir,omitempty"`
	TemplateDirs []string              `json:"template_dirs,omitempty"`
	Default      *PageParam            `json:"default,omitempty"`
	Pages        map[string]*PageParam `json:"pages,omitempty"`
}

type PageParam struct {
	Enable      *bool          `json:"enable,omitempty"`
	Route       string         `json:"route,omitempty"`
	Cache       *bool          `json:"cache,omitempty"`
	Layout      string         `json:"layout,omitempty"`
	Render      string         `json:"render,omitempty"`
	Content     string         `json:"content,omitempty"`
	ContentFile string         `json:"content_file,omitempty"`
	ContentUrl  string         `json:"content_url,omitempty"`
	SubBrand    string         `json:"sub_brand,omitempty"`
	Menu        string         `json:"menu,omitempty"`
	Css         []string       `json:"css,omitempty"`
	Js          []string       `json:"js,omitempty"`
	Js99        []string       `json:"js99,omitempty"`
	LibHead     []string       `json:"lib_head,omitempty"`
	Lib         []string       `json:"lib,omitempty"`
	Lib99       []string       `json:"lib99,omitempty"`
	NavMenu     []MenuItem     `json:"nav_menu,omitempty"`
	Meta        PageMeta       `json:"meta,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
	Rt          map[string]any `json:"_runtime,omitempty"` // Runtime data
}

type MenuItem struct {
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"` // Button, Link,
	Href  string `json:"href,omitempty"`
	Link  string `json:"link,omitempty"`
	Data  string `json:"data,omitempty"`
}
