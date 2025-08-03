// config/db
package cfg

import (
	"xi/app/lib"
)

var env = lib.Env

type DbConf struct {
	DefaultDb  string
	DefaultRdb string
	RdbPrefix  string
	Connection map[string]DbParam
}
type DbParam struct {
	Enable        bool
	Db      string
	Rdb       int
	User          string
	Pass          string
	Driver        string
	Host          string
	Port          string
	Socket        string
	Charset       string
	Collation     string
	Prefix        string
	PrefixIndexes bool
	Strict        bool
	Engine        string
}

var Db = &DbConf {
	DefaultDb:  env.Get("DEFAULT_DB", "db"),
	DefaultRdb: env.Get("DEFAULT_RDB", "rdb"),
	RdbPrefix:  env.Get("APP_ABBR", "redis"),

	Connection: map[string]DbParam{

		"db": {
			Enable:        true,
			Db:      env.Get("DB_XI", "XI"),
			User:          env.Get("DB_XI_USER"),
			Pass:          env.Get("DB_XI_PASS"),
			Driver:        env.Get("DB_XI_DRIVER", "mysql"),
			Host:          env.Get("DB_XI_HOST", "127.0.0.1"),
			Port:          env.Get("DB_XI_PORT", "3306"),
			Charset:       env.Get("DB_CHARSET", "utf8mb4"),
			Collation:     env.Get("DB_COLLATION", "utf8mb4_unicode_ci"),
			Socket:        env.Get("DB_SOCKET", ""),
			Prefix:        "",
			PrefixIndexes: true,
			Strict:        true,
			Engine:        "",
		},

		"rdb": {
			Enable:  env.Bool("RDB_Enable", true),
			Rdb: env.Int("RDB", 0),
			Host:    env.Get("RDB_HOST", "127.0.0.1"),
			Port:    env.Get("RDB_PORT", "6379"),
			Pass:    env.Get("RDB_PASS", ""),
			Driver:  "redis",
		},
	},
}
