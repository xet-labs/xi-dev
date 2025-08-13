package cfg

import (
	"xi/app/model"
)

type Lib struct{}

var (
	global = &model.Config{}
	App    = &global.App
	Build  = &global.Build
	Db     = &global.Db
	Path   = &global.Path
	View   = &global.View
)

func GetStatic() any { return globalStatic }

func Set(cfg model.Config) { *global = cfg }

func Get() *model.Config { return global }

func Update(cfg model.Config) {
	*global = cfg

	App = &global.App
	Build = &global.Build
	Db = &global.Db
	Path = &global.Path
	View = &global.View
}
