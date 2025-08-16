package cfg

// --------------------
// Build-time constants
// --------------------
// These are injected at build time with -ldflags "-X"
var (
	BuildDate     = "0-0-0"
	BuildName     = "dev"
	BuildRevision = "00000"
	BuildVersion  = "v0.0.0"
)
