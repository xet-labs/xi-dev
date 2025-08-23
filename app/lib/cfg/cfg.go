package cfg

import "xi/app/model"

// --------------------
// Runtime config (mutable)
// --------------------
var global = &model.Config{}

// Direct pointers for convenience
var (
	App   = &global.App
	Brand = &global.Brand
	Db    = &global.Db
	View  = &global.View
)

// Static BuildConf (never changes at runtime)
var Build = model.BuildConf{
	Date:     BuildDate,
	Name:     BuildName,
	Revision: BuildRevision,
	Version:  BuildVersion,
}


// Get returns current runtime config
func Get() *model.Config { return global }

// Set replaces the entire config (except Build, which stays static)
func Set(cfg model.Config) {
	cfg.Build = Build         // enforce static build info
	*global = cfg
}

// Update merges in a new config but keeps Build static
func Update(cfg model.Config) {
	cfg.Build = Build
	*global = cfg
	App   = &global.App
	Brand = &global.Brand
	Db    = &global.Db
	View  = &global.View
}