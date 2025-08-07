package lib

import (
	"xi/app/lib/conf"
	"xi/app/lib/db"
	"xi/app/lib/env"
	"xi/app/lib/file"
	"xi/app/lib/minify"
	"xi/app/lib/view"
	// "xi/app/lib/hook"
)

var (
	Conf = conf.Conf
	Db = db.Db
	Rdb = db.Rdb
	Env = env.Env
	File = file.File
	// Hook = &hook.Hook{}
	Minify = minify.Minify
	View = view.View
)