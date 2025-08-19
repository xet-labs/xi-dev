package util

import (
	"xi/app/lib/util/file"
	"xi/app/lib/util/maps"
	"xi/app/lib/util/str"
	"xi/app/lib/util/url"
)

type UtilLib struct {
	File file.FileLib
	Map  maps.MapsLib
	Str  str.StrLib
	Url  url.UrlLib
}

var Util = &UtilLib{
	File: file.FileLib{},
	Map:  maps.MapsLib{},
	Str:  str.StrLib{},
	Url:  url.UrlLib{},
}

// expose shortcuts
var (
	File = &Util.File
	Map  = &Util.Map
	Str  = &Util.Str
	Url  = &Util.Url
)
