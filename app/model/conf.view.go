package model

type ViewConf struct {
	CssDir      []string             `json:"css_dir,omitempty"`
	TemplateDir []string             `json:"template_dir,omitempty"`
	Default     PageParam            `json:"default,omitempty"`
	Pages       map[string]PageParam `json:"pages,omitempty"`
}

type PageParam struct {
	Route      string         `json:"route,omitempty"`
	RenderType string         `json:"render_type"`
	File       string         `json:"file,omitempty"`
	Layout     string         `json:"layout,omitempty"`
	FullHtml   string         `json:"full_html,omitempty"`
	BodyHtml   string         `json:"body_html,omitempty"`
	PartHtml   string         `json:"part_html,omitempty"`
	SubBrand   string         `json:"sub_brand,omitempty"`
	Menu       string         `json:"menu,omitempty"`
	Css        []string       `json:"css,omitempty"`
	Js         []string       `json:"js,omitempty"`
	Js99       []string       `json:"js99,omitempty"`
	LibHead    []string       `json:"lib_head,omitempty"`
	Lib        []string       `json:"lib,omitempty"`
	Lib99      []string       `json:"lib99,omitempty"`
	NavMenu    []MenuItem     `json:"nav_menu,omitempty"`
	Meta       PageMeta       `json:"meta"`
	Data       map[string]any `json:"data"`
}

type MenuItem struct {
	Label string `json:"label,omitempty"`
	Href  string `json:"href,omitempty"`
}
