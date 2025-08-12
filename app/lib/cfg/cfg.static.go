package cfg

import (
	"xi/app/model"
)

var (
	// Defaults (to be overridden at build time)
	BuildDate     = "0-0-0"
	BuildName     = "dev"
	BuildRevision = "00000"
	BuildVersion  = "v0.0.0"
)

var globalStatic = map[string]any{
	"build": model.BuildConf{
		Date:     BuildDate,
		Name:     BuildName,
		Revision: BuildRevision,
		Version:  BuildVersion,
	},
}

func init() {
	*Build = model.BuildConf{
		Date:     BuildDate,
		Name:     BuildName,
		Revision: BuildRevision,
		Version:  BuildVersion,
	}
}
