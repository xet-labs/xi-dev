package cfg

import (
	"xi/app/model"
)

type CfgLib struct{}

// var Cfg = &CfgLib{}

var (
	global = &model.Config{}
	App    = &global.App
	Db     = &global.Db
	View   = &global.View
)

func Set(cfg model.Config) {
	*global = cfg
}

func Update(cfg model.Config) {
	*global = cfg
	App = &global.App
	Db = &global.Db
	View = &global.View
}

func Get() *model.Config {
	return global
}
