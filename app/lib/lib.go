package lib

import (
	"xi/app/lib/url"
	"xi/app/lib/conf"
	"xi/app/lib/db"
	"xi/app/lib/env"
	"xi/app/lib/file"
	"xi/app/lib/logger"

	// "github.com/rs/zerolog/log"
	"xi/app/lib/hook"
	"xi/app/lib/minify"
	"xi/app/lib/util"
	"xi/app/lib/view"
)

var (
	Conf   = conf.Conf
	Db     = db.Db
	Rdb    = db.Rdb
	Env    = env.Env
	File   = file.File
	Hook   = &hook.Hook{}
	Log    = logger.Logger.Log
	Logger = logger.Logger
	Minify = minify.Minify
	Util   = util.Util
	Url    = url.Url
	View   = view.View
)
