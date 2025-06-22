// config/db
package config

import (
	"xi/app/util"
)

type DBParam struct {
	Database      string
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

var DB = map[string]DBParam{
	"XI": {
		Database:      util.Env("DB_XI", "XI"),
		User:          util.Env("DB_XI_USER"),
		Pass:          util.Env("DB_XI_PASS"),
		Driver:        util.Env("DB_XI_DRIVER", "mysql"),
		Host:          util.Env("DB_XI_HOST", "127.0.0.1"),
		Port:          util.Env("DB_XI_PORT", "3306"),
		UnixSocket:    util.Env("DB_SOCKET", ""),
		Charset:       util.Env("DB_CHARSET", "utf8mb4"),
		Collation:     util.Env("DB_COLLATION", "utf8mb4_unicode_ci"),
		Prefix:        "",
		PrefixIndexes: true,
		Strict:        true,
		Engine:        "",
		Enable:        true,
	},
	"Redis": {
		Driver:   "redis",
		Database: util.Env("DB_REDIS", "0"),
		Host:     util.Env("DB_REDIS_HOST", "127.0.0.1"),
		Port:     util.Env("DB_REDIS_PORT", "6379"),
		Pass:     util.Env("DB_REDIS_PASS", ""),
		Enable:   true,
	},
	"App": {
		Driver:    util.Env("DB_APP_DRIVER", "sqlite"),
		Database:  util.Env("DB_APP_NAME", "app"),
		Charset:   "utf8mb4",
		Collation: "utf8mb4_unicode_ci",
	},
}
