package cfg

import (
	"xi/app/model"
)

type Lib struct{}

var (
	global = &model.Config{}
	Build  = &global.Build
	App    = &global.App
	Db     = &global.Db
	View   = &global.View
)

func GetStatic() any { return globalStatic }

func Set(cfg model.Config) { *global = cfg }

func Get() *model.Config { return global }

func Update(cfg model.Config) {
	*global = cfg
	Build = &global.Build
	App = &global.App
	Db = &global.Db
	View = &global.View
}
