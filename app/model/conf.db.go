package model

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
