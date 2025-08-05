package lib

import (
	"log"
	"sync"

	"gorm.io/gorm"
)

// Central utility
type DBLib struct {
	clients   map[string]*gorm.DB
	defaultCli string
	mu        sync.RWMutex
	once      sync.Once
	lazyInit  func()
}

// Global instance
var DB = &DBLib{
	defaultCli: "database",
	clients:   make(map[string]*gorm.DB),
}

// RegisterLazyFn sets a callback for deferred initialization.
func (d *DBLib) RegisterLazyFn(fn func()) {
	d.lazyInit = fn
}

// Ensures lazyInit runs once
func (d *DBLib) lazyFnOnce() {
	d.once.Do(func() {
		if d.lazyInit != nil {
			d.lazyInit()
		}
	})
}

// Get returns the DB instance by name or default
func (d *DBLib) GetCli(name ...string) *gorm.DB {
	d.lazyFnOnce()
	d.mu.RLock()
	defer d.mu.RUnlock()

	dbName := d.defaultCli
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	if db, ok := d.clients[dbName]; ok {
		return db
	}

	log.Printf("⚠️  Database '%s' not found", dbName)
	return nil
}

// Set sets a DB by name
func (d *DBLib) SetCli(name string, db *gorm.DB) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.clients[name] = db
}

// SetDefault sets the default DB name
func (d *DBLib) SetDefault(name string) {
	d.defaultCli = name
}

// You can similarly add Redis setters/getters if needed
