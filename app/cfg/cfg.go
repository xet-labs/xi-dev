package cfg

import "xi/app/schema"

type CfgLib struct{
	schema.Config
}

var Cfg = &CfgLib{}