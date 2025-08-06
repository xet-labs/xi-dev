package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"xi/app/lib"
	"xi/app/cfg"

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
	// Ensure the 'DB.Init' is called in lib.{DB, Redis} so that core env is setup
	lib.DB.RegisterLazyFn(DB.Init)
	lib.Redis.RegisterLazyFn(DB.Init)
}
// Init initializes DBs once
func (d *DBService) Init() { d.once.Do(d.InitForce) }

func (d *DBService) initPre() {
	lib.Cfg.Init()
	
	// Set global Redis and DB defaults
	lib.DB.SetDefault(cfg.Db.DbDefault)
	lib.Redis.SetCtx(context.Background())
	lib.Redis.SetDefault(cfg.Db.RdbDefault)
	lib.Redis.SetPrefix(cfg.Db.RdbPrefix)
}
func (d *DBService) initPost() {}

// Initializes all DBs and Redis clients (forced)
func (d *DBService) InitForce() {
	d.initPre()

	if cfg.Db.Conn == nil {
		log.Println("⚠️  [DB] WRN: No connections were configured")
	}
	for profile, c := range cfg.Db.Conn {
		if !c.Enable {
			continue
		}

		// Fallback for DB credentials
		if c.User == "" {
			c.User = c.Db + "_u"
		}
		if c.Pass == "" {
			c.Pass = lib.Env.Get("DB_PASS")
		}

		switch c.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				c.User, c.Pass, c.Host, c.Port, c.Db, c.Charset)
			dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			lib.DB.SetCli(profile, dbConn)
			log.Printf("✅ DB connected '%s' (MySQL)", profile)

		case "sqlite":
			dbConn, err := gorm.Open(sqlite.Open(c.Db), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			lib.DB.SetCli(profile, dbConn)
			log.Printf("✅ DB connected '%s' (SQLite)", profile)

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     c.Host + ":" + c.Port,
				Password: c.Pass,
				DB:       c.Rdb,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis '%s': %v", profile, err)
			}
			lib.Redis.SetCli(profile, rdb)
			log.Printf("✅ DB connected '%s' (Redis)", profile)

		default:
			log.Printf("⚠️  Unsupported DB driver '%s' for DB '%s'", c.Driver, profile)
		}
	}
	
	d.initPost()
}
