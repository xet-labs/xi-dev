package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"xi/app/lib/conf"
	"xi/app/lib/cfg"
	"xi/app/lib/env"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Central utility
type DbLib struct {
	clients   map[string]*gorm.DB
	defaultCli string
	mu        sync.RWMutex
	once      sync.Once
	lazyInit  func()
}

// Global instance
var Db = &DbLib{
	defaultCli: "database",
	clients:   make(map[string]*gorm.DB),
}

// RegisterLazyFn sets a callback for deferred initialization.
func (d *DbLib) RegisterLazyFn(fn func()) {
	d.lazyInit = fn
}

// Init initializes DBs once
func init() { Db.Init() }
func (d *DbLib) Init() { d.once.Do(d.InitForce) }

func (d *DbLib) initPre() {
	conf.Conf.Init()
	
	// Set global Redis and DB defaults
	d.SetDefault(cfg.Db.DbDefault)
	Rdb.SetCtx(context.Background())
	Rdb.SetDefault(cfg.Db.RdbDefault)
	Rdb.SetPrefix(cfg.Db.RdbPrefix)
}
func (d *DbLib) initPost() {}

// Initializes all DBs and Redis clients (forced)
func (d *DbLib) InitForce() {
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
			c.Pass = env.Env.Get("DB_PASS")
		}

		switch c.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				c.User, c.Pass, c.Host, c.Port, c.Db, c.Charset)
			dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			Db.SetCli(profile, dbConn)
			log.Printf("✅ [DB] \tConnected '%s' (MySQL)", profile)

		case "sqlite":
			dbConn, err := gorm.Open(sqlite.Open(c.Db), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			Db.SetCli(profile, dbConn)
			log.Printf("✅ [DB] \tConnected '%s' (SQLite)", profile)

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     c.Host + ":" + c.Port,
				Password: c.Pass,
				DB:       c.Rdb,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis '%s': %v", profile, err)
			}
			Rdb.SetCli(profile, rdb)
			log.Printf("✅ [DB] \tConnected '%s' (Redis)", profile)

		default:
			log.Printf("⚠️  Unsupported DB driver '%s' for DB '%s'", c.Driver, profile)
		}
	}
	
	d.initPost()
}
