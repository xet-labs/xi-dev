// services/db
package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"xi/app/lib"
	"xi/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DBs       = make(map[string]*gorm.DB)
	DbDef     = config.DbConf.DefDb
	RedisClis = lib.RedisClis
	dbLock    sync.RWMutex
)

// Init initializes all enabled databases based on the config
func InitDB() {
	if lib.EnvBool("DBInitialized") {
		return
	}

	for name, conf := range config.DB {
		if !conf.Enable {
			// log.Printf("⚠️  DB '%s' skipped", name)
			continue
		}

		//- DBUser fallback
		if conf.User == "" {
			conf.User = conf.Database + "_u"
		}
		if conf.Pass == "" {
			conf.Pass = lib.Env("DB_PASS")
		}

		switch conf.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				conf.User, conf.Pass, conf.Host, conf.Port, conf.Database, conf.Charset)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to %s DB: %v", name, err)
			}
			dbLock.Lock()
			DBs[name] = db
			dbLock.Unlock()
			log.Printf("✅ DB connected '%s' (MySQL)", name)
			continue

		case "sqlite":
			db, err := gorm.Open(sqlite.Open(conf.Database), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to %s DB: %v", name, err)
			}
			dbLock.Lock()
			DBs[name] = db
			dbLock.Unlock()
			log.Printf("✅ DB connected '%s' (SQLite)", name)
			continue

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     conf.Host + ":" + conf.Port,
				Password: conf.Pass,
				DB:       conf.RedisDB,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis DB: %v", err)
			}

			RedisClis[name] = rdb
			log.Printf("✅ DB connected '%s' (Redis)", name)
			continue

		default:
			log.Printf("❌ Unsupported DB driver '%s' for DB '%s'", conf.Driver, name)
			continue
		}
	}

	lib.RedisPrefix = config.DbConf.RedisPrefix
	lib.RedisDefRdb = config.DbConf.RedisDefRdb
	if rdb, ok := lib.RedisClis[lib.RedisDefRdb]; ok {
		lib.RedisDefCli = rdb
	}

	lib.EnvSet("DBInitialized", true)
}

// DB safely returns the DB instance by name
func DB(name ...string) *gorm.DB {

	if !lib.EnvBool("DBInitialized") {
		InitDB()
	}

	// Use read lock to safely access DB map
	dbLock.RLock()
	defer dbLock.RUnlock()

	dbName := DbDef
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	if db, ok := DBs[dbName]; ok {
		return db
	}

	log.Printf("⚠️  Requested unknown DB '%s'", dbName)
	return nil
}
