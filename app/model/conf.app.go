package model

type AppConf struct {
	Mode         string   `json:"mode,omitempty"`
	Port         string   `json:"port,omitempty"`
	Env          string   `json:"env,omitempty"`
	EnvFiles     []string `json:"env_files,omitempty"`
	SslCert      string   `json:"ssl_cert,omitempty"`
	SslCertFiles []string `json:"ssl_cert_files,omitempty"`
	TlsCert      string   `json:"tls_cert,omitempty"`
	TlsCertFiles []string `json:"tls_cert_files,omitempty"`

	ForceCachePage bool 	`json:"force_cache_Page,omitempty"`

}
