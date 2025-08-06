package cfg

import "xi/app/lib/schema"

var Global schema.Config

var (
	App  *schema.AppConf
	Db   *schema.DbConf
	View *schema.ViewConf
)