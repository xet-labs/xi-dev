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

type DBService struct {
	once sync.Once
}

var DB = &DBService{}

func init() {
	lib.DB.RegisterLazyInit(DB.Init)
	lib.Redis.RegisterLazyInit(DB.Init)
}

// InitForce initializes all DBs and Redis clients (forced)
func (d *DBService) InitForce() {
	// Set global Redis and DB defaults
	lib.DB.SetDefault(config.DbConf.DefDb)
	lib.Redis.SetCtx(context.Background())
	lib.Redis.SetDefault(config.DbConf.RedisDefRdb)
	lib.Redis.SetPrefix(config.DbConf.RedisPrefix)

	for name, conf := range config.DB {
		if !conf.Enable {
			continue
		}

		// Fallback for DB credentials
		if conf.User == "" {
			conf.User = conf.Database + "_u"
		}
		if conf.Pass == "" {
			conf.Pass = lib.Env.Get("DB_PASS")
		}

		switch conf.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				conf.User, conf.Pass, conf.Host, conf.Port, conf.Database, conf.Charset)
			dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", name, err)
			}
			lib.DB.SetCli(name, dbConn)
			log.Printf("✅ DB connected '%s' (MySQL)", name)

		case "sqlite":
			dbConn, err := gorm.Open(sqlite.Open(conf.Database), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", name, err)
			}
			lib.DB.SetCli(name, dbConn)
			log.Printf("✅ DB connected '%s' (SQLite)", name)

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     conf.Host + ":" + conf.Port,
				Password: conf.Pass,
				DB:       conf.RedisDB,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis '%s': %v", name, err)
			}
			lib.Redis.SetCli(name, rdb)
			log.Printf("✅ RDB connected '%s'", name)

		default:
			log.Printf("⚠️  Unsupported DB driver '%s' for DB '%s'", conf.Driver, name)
		}
	}

}

// Init initializes DBs once
func (d *DBService) Init() {
	d.once.Do(d.InitForce)
}
