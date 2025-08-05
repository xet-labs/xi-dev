// config/db
package cfg

import (
	"xi/app/lib"
	"xi/app/schema"
)

var env = lib.Env

var Db = &schema.DbConf{
	DbDefault:  env.Get("DEFAULT_DB", "db"),
	RdbDefault: env.Get("DEFAULT_RDB", "rdb"),
	RdbPrefix:  env.Get("APP_ABBR", "redis"),

	Conn: map[string]schema.DbParam{

		"db": {
			Enable:        true,
			Db:            env.Get("DB", "XI"),
			User:          env.Get("DB_USER"),
			Pass:          env.Get("DB_PASS"),
			Driver:        env.Get("DB_DRIVER", "mysql"),
			Host:          env.Get("DB_HOST", "127.0.0.1"),
			Port:          env.Get("DB_PORT", "3306"),
			Charset:       env.Get("DB_CHARSET", "utf8mb4"),
			Collation:     env.Get("DB_COLLATION", "utf8mb4_unicode_ci"),
			Socket:        env.Get("DB_SOCKET", ""),
			Prefix:        "",
			PrefixIndexes: true,
			Strict:        true,
			Engine:        "",
		},

		"rdb": {
			Enable: env.Bool("RDB_ENABLE", true),
			Rdb:    env.Int("RDB", 0),
			Host:   env.Get("RDB_HOST", "127.0.0.1"),
			Port:   env.Get("RDB_PORT", "6379"),
			Pass:   env.Get("RDB_PASS", ""),
			Driver: "redis",
		},
	},
}
