package cfg

import "xi/app/model"

var (
	global = &model.Config{}
	App    = &global.App
	Db     = &global.Db
	View   = &global.View
)

func Update(cfg model.Config) {
	*global = cfg // copies all values into the pointed struct
}

func Get() *model.Config {
	return global
}
