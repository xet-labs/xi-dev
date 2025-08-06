package schema

type Config struct {
	App  AppConf  `json:"app"`
	DB   DbConf   `json:"db"`
	View ViewConf `json:"view"`
}

type AppConf struct {
	Abbr        string   `json:"abbr"`
	Name        string   `json:"name"`
	Domain      string   `json:"domain"`
	Url         string   `json:"url"`
	Mode        string   `json:"mode"`
	Port        string   `json:"port"`
	Tags        []string `json:"tags"`
	Logo        []string `json:"logo"`
	FeaturedImg []string `json:"featured_img"`
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
	CssDir      []string             `json:"css_dir"`
	TemplateDir []string             `json:"template_dir"`
	Pages       map[string]PageParam `json:"pages"`
}

type PageParam struct {
	App         AppConf  `json:"app"`
	Layout      string     `json:"layout"`
	NavMenu     []MenuItem `json:"nav_menu"`
	Route       string     `json:"route"`
	Title       string     `json:"title,omitempty"`
	File        string     `json:"file,omitempty"`
	Template    string     `json:"template,omitempty"`
	Menu        string     `json:"menu,omitempty"`
	Description string     `json:"description,omitempty"`
	SubBrand    string     `json:"sub_brand,omitempty"`
	LibHLJS     bool       `json:"lib_hljs,omitempty"`
	Js99        []string   `json:"js99,omitempty"`
	Meta        *PageMeta  `json:"meta,omitempty"`
}

type PageMeta struct {
	Canonical string   `json:"canonical"`
	Tags      []string `json:"tags,omitempty"`
}

type MenuItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}
