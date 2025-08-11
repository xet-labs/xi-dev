package cfg

import (
	"xi/app/model"
)

var (
	// Defaults (to be overridden at build time)
	BuildDate     = "0-0-0"
	BuildRevision = "dev"
	BuildVersion  = "vx.x.x"
)

var globalStatic = map[string]any{
	"build": model.BuildConf{
		Date:     BuildDate,
		Revision: BuildRevision,
		Version:  BuildVersion,
	},
}