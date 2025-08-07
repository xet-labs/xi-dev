package cfg

import "xi/app/model"

type CfgLib struct{}

// var Cfg = &CfgLib{}

var (
	global = &model.Config{}
	App    = &global.App
	Db     = &global.Db
	View   = &global.View
)

func Update(cfg model.Config) {
	*global = cfg
}

func Get() *model.Config {
	return global
}
