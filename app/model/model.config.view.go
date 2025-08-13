package model

type ViewConf struct {
	CssDir      []string             `json:"css_dir,omitempty"`
	TemplateDir []string             `json:"template_dir,omitempty"`
	Default     PageParam            `json:"default,omitempty"`
	Pages       map[string]PageParam `json:"pages,omitempty"`
}

type PageParam struct {
	Title       string         `json:"title,omitempty"`
	Route       string         `json:"route,omitempty"`
	Description string         `json:"description,omitempty"`
	RenderType  string         `json:"render_type"`
	File        string         `json:"file,omitempty"`
	Layout      string         `json:"layout,omitempty"`
	FullHtml    string         `json:"full_html,omitempty"`
	BodyHtml    string         `json:"body_html,omitempty"`
	PartHtml    string         `json:"part_html,omitempty"`
	SubBrand    string         `json:"sub_brand,omitempty"`
	Menu        string         `json:"menu,omitempty"`
	Css         []string       `json:"css,omitempty"`
	Js          []string       `json:"js,omitempty"`
	Js99        []string       `json:"js99,omitempty"`
	HeadLib     []string       `json:"head_lib,omitempty"`
	Lib         []string       `json:"lib,omitempty"`
	Lib99       []string       `json:"lib99,omitempty"`
	Meta        PageMeta       `json:"meta"`
	App         AppConf        `json:"app"`
	NavMenu     []MenuItem     `json:"nav_menu,omitempty"`
	Data        map[string]any `json:"data"`
}

type PageMeta struct {
	Canonical string   `json:"canonical,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type MenuItem struct {
	Label string `json:"label,omitempty"`
	Href  string `json:"href,omitempty"`
}
