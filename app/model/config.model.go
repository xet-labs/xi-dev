package model

type Config struct {
	App  AppConf  `json:"app"`
	Db   DbConf   `json:"db"`
	View ViewConf `json:"view"`
}

type AppConf struct {
	Abbr        string   `json:"abbr,omitempty"`
	Name        string   `json:"name,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Url         string   `json:"url,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	Port        string   `json:"port,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Logo        []string `json:"logo,omitempty"`
	FeaturedImg []string `json:"featured_img,omitempty"`
}

type DbConf struct {
	DbDefault  string             `json:"db_default"`
	RdbDefault string             `json:"rdb_default"`
	RdbPrefix  string             `json:"rdb_prefix"`
	Conn       map[string]DbParam `json:"conn"`
}
type DbParam struct {
	Enable        bool   `json:"enable"`
	Db            string `json:"db"`
	Rdb           int    `json:"rdb"`
	User          string `json:"user"`
	Pass          string `json:"pass"`
	Driver        string `json:"driver"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	Engine        string `json:"engine"`
	Socket        string `json:"socket"`
	Charset       string `json:"charset"`
	Collation     string `json:"collation"`
	Prefix        string `json:"prefix"`
	PrefixIndexes bool   `json:"prefixindexes"`
	Strict        bool   `json:"strict"`
}

type ViewConf struct {
	CssDir      []string             `json:"css_dir,omitempty"`
	TemplateDir []string             `json:"template_dir,omitempty"`
	Default     PageParam            `json:"default,omitempty"`
	Page        PageParam            `json:"page,omitempty"`
	Pages       map[string]PageParam `json:"pages"`
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
