// config/db
package config

import (
	"xi/app/lib"
)

var DbConf = struct {
	DefDb       string
	RedisDefRdb string
	RedisPrefix string
	MysqlDb     string
	PostgresDb  string
}{
	DefDb:       lib.Env("DB_DEFAULT", "sql"),
	RedisDefRdb: lib.Env("DB_REDIS_DEFAULT", "redis"),
	RedisPrefix: lib.Env("APP_ABBR", "redis"),
}

var DB = map[string]DBParam{
	"sql": {
		Database:      lib.Env("DB_XI", "XI"),
		User:          lib.Env("DB_XI_USER"),
		Pass:          lib.Env("DB_XI_PASS"),
		Driver:        lib.Env("DB_XI_DRIVER", "mysql"),
		Host:          lib.Env("DB_XI_HOST", "127.0.0.1"),
		Port:          lib.Env("DB_XI_PORT", "3306"),
		UnixSocket:    lib.Env("DB_SOCKET", ""),
		Charset:       lib.Env("DB_CHARSET", "utf8mb4"),
		Collation:     lib.Env("DB_COLLATION", "utf8mb4_unicode_ci"),
		Prefix:        "",
		PrefixIndexes: true,
		Strict:        true,
		Engine:        "",
		Enable:        true,
	},
	"redis": {
		Driver:  "redis",
		RedisDB: lib.EnvInt("DB_REDIS", 0),
		Host:    lib.Env("DB_REDIS_HOST", "127.0.0.1"),
		Port:    lib.Env("DB_REDIS_PORT", "6379"),
		Pass:    lib.Env("DB_REDIS_PASS", ""),
		Enable:  true,
	},
	"app": {
		Driver:    lib.Env("DB_APP_DRIVER", "sqlite"),
		Database:  lib.Env("DB_APP_NAME", "app"),
		Charset:   "utf8mb4",
		Collation: "utf8mb4_unicode_ci",
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
