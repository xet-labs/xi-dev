package model

type Config struct {
	App   AppConf   `json:"app"`
	Brand BrandConf `json:"brand"`
	Build BuildConf `json:"build"`
	Db    DbConf    `json:"db"`
	View  ViewConf  `json:"view"`
}

type AppConf struct {
	Mode         string   `json:"mode,omitempty"`
	Port         string   `json:"port,omitempty"`
	Env          string   `json:"env,omitempty"`
	EnvFiles     []string `json:"env_files,omitempty"`
	SslCert      string   `json:"ssl_cert,omitempty"`
	SslCertFiles []string `json:"ssl_cert_files,omitempty"`
	TlsCert      string   `json:"tls_cert,omitempty"`
	TlsCertFiles []string `json:"tls_cert_files,omitempty"`
}

type BrandConf struct {
	Name        string   `json:"name,omitempty"`
	Abbr        string   `json:"abbr,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Url         string   `json:"url,omitempty"`
	Logo        []string `json:"logo"`
	FeaturedImg []string `json:"featured_img,omitempty"`
	Tagline     string   `json:"tagline,omitempty"`
}

type BuildConf struct {
	Date     string `json:"date,omitempty"`
	Name     string `json:"name,omitempty"`
	Revision string `json:"revision,omitempty"`
	Version  string `json:"version,omitempty"`
}
