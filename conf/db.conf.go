// config/db
package conf

import (
	"xi/app/lib"
)
var Env = lib.Env

var Db = struct {
	Def         string
	RedisDef    string
	RedisPrefix string
	DB          map[string]DBParam

}{
	Def:         Env.Get("DB_DEFAULT", "DB"),
	RedisDef:    Env.Get("DB_REDIS_DEFAULT", "RDB"),
	RedisPrefix: Env.Get("APP_ABBR", "redis"),

	DB: map[string]DBParam {

		"DB": {
			Enable:        true,
			Database:      Env.Get("DB_XI", "XI"),
			User:          Env.Get("DB_XI_USER"),
			Pass:          Env.Get("DB_XI_PASS"),
			Driver:        Env.Get("DB_XI_DRIVER", "mysql"),
			Host:          Env.Get("DB_XI_HOST", "127.0.0.1"),
			Port:          Env.Get("DB_XI_PORT", "3306"),
			UnixSocket:    Env.Get("DB_SOCKET", ""),
			Charset:       Env.Get("DB_CHARSET", "utf8mb4"),
			Collation:     Env.Get("DB_COLLATION", "utf8mb4_unicode_ci"),
			Prefix:        "",
			PrefixIndexes: true,
			Strict:        true,
			Engine:        "",
		},

		"RDB": {
			Enable:  Env.Bool("DB_REDIS_Enabled", true),
			Driver:  "redis",
			RedisDB: Env.Int("DB_REDIS", 0),
			Host:    Env.Get("DB_REDIS_HOST", "127.0.0.1"),
			Port:    Env.Get("DB_REDIS_PORT", "6379"),
			Pass:    Env.Get("DB_REDIS_PASS", ""),
		},
	},
}


type DBParam struct {
	Database      string
	RedisDB       int
	User          string
	Pass          string
	Driver        string
	Host          string
	Port          string
	UnixSocket    string
	Charset       string
	Collation     string
	Prefix        string
	PrefixIndexes bool
	Strict        bool
	Engine        string
	Enable        bool
}
