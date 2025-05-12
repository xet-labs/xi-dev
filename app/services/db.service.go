// services/db
package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"xi/app/global"
	"xi/app/utils"
	"xi/config"
)

var (
	DBs         = make(map[string]*gorm.DB)
	dbMu        sync.RWMutex
	RedisClients = make(map[string]*redis.Client)
)

// Init initializes all enabled databases based on the config
func InitDB() {
	if global.DBInitialized { return }

	log.Println("Init databases...")
	for name, conf := range config.DB {
		if !conf.Enable {
			log.Printf("⚠️  DB '%s' skipped", name)
			continue
		}

		//- DBUser fallback 
		if conf.User == "" { conf.User = conf.Database + "_u" }
		if conf.Pass == "" { conf.Pass = utils.Env("DB_PASS") }

		switch conf.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				conf.User, conf.Pass, conf.Host, conf.Port, conf.Database, conf.Charset)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to %s DB: %v", name, err)
			}
			dbMu.Lock()
			DBs[name] = db
			dbMu.Unlock()
			log.Printf("✅ DB '%s' connected (MySQL)", name)
			continue

		case "sqlite":
			db, err := gorm.Open(sqlite.Open(conf.Database), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to %s DB: %v", name, err)
			}
			dbMu.Lock()
			DBs[name] = db
			dbMu.Unlock()
			log.Printf("✅ DB '%s' connected (SQLite)", name)
			continue

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     conf.Host + ":" + conf.Port,
				Password: conf.Pass,
				DB:       0,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis DB: %v", err)
			}

			RedisClients[name] = rdb
			log.Printf("✅ DB '%s' connected (Redis)", name)
			continue

		default:
			log.Printf("❌ Unsupported DB driver '%s' for DB '%s'", conf.Driver, name)
			continue
		}
	}

	global.DBInitialized = true
}

// DB safely returns the DB instance by name
func DB(name ...string) *gorm.DB {

	if !global.DBInitialized { InitDB() }

	// Use read lock to safely access DB map
	dbMu.RLock()
	defer dbMu.RUnlock()

	dbName := utils.Env("DB_DEFAULT", "XI")
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	if db, ok := DBs[dbName]; ok {
		return db
	}

	log.Printf("⚠️  Requested unknown DB '%s'", dbName)
	return nil
}

func Redis(name ...string) *redis.Client {
	if !global.DBInitialized {
		InitDB()
	}

	// dbMu.RLock()
	// defer dbMu.RUnlock()

	dbName := utils.Env("DB_REDIS_DEFAULT", "Redis")
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	if rdb, ok := RedisClients[dbName]; ok {
		return rdb
	}

	log.Printf("⚠️  Requested unknown Redis client '%s'", dbName)
	return nil
}

