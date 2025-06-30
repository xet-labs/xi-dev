// config/db
package config

import (
	"xi/app/lib"
)

var Env = lib.Env

var DbConf = struct {
	DefDb       string
	RedisDefRdb string
	RedisPrefix string
	MysqlDb     string
	PostgresDb  string
}{
	DefDb:       Env.Get("DB_DEFAULT", "DB"),
	RedisDefRdb: Env.Get("DB_REDIS_DEFAULT", "eeefeef"),
	RedisPrefix: Env.Get("APP_ABBR", "redis"),
}

var DB = map[string]DBParam{
	"DB": {
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
		Enable:        true,
	},
	"RDB": {
		Driver:  "redis",
		RedisDB: Env.Int("DB_REDIS", 0),
		Host:    Env.Get("DB_REDIS_HOST", "127.0.0.1"),
		Port:    Env.Get("DB_REDIS_PORT", "6379"),
		Pass:    Env.Get("DB_REDIS_PASS", ""),
		Enable:  true,
	},
	"app": {
		Driver:    Env.Get("DB_APP_DRIVER", "sqlite"),
		Database:  Env.Get("DB_APP_NAME", "app"),
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
