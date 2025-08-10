package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Ensures lazyInit runs once
func (d *DbLib) lazyFnOnce() {
	d.once.Do(func() {
		if d.lazyInit != nil {
			d.lazyInit()
		}
	})
}

// Get returns the DB instance by name or default
func (d *DbLib) GetCli(name ...string) *gorm.DB {
	d.Init()
	d.mu.RLock()
	defer d.mu.RUnlock()

	dbName := d.defaultCli
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	if db, ok := d.clients[dbName]; ok {
		return db
	}

	log.Warn().Msgf("DB '%s' not found", dbName)
	return nil
}

// Set sets a DB by name
func (d *DbLib) SetCli(name string, db *gorm.DB) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.clients[name] = db
}

// SetDefault sets the default DB name
func (d *DbLib) SetDefault(name string) {
	d.defaultCli = name
}

// You can similarly add Redis setters/getters if needed
