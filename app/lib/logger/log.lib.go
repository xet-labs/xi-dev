package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogLib struct {
	Log zerolog.Logger
}

var Logger = &LogLib{}

var (
	Log = Logger.Log
)

// func init() {
// 	Logger.Init()
// }

func (l *LogLib) Init() {
	zerolog.TimeFieldFormat = time.RFC3339 // Keep consistent
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		// TimeFormat: time.RFC3339,
		NoColor: false,
		FormatTimestamp: func(i any) string {
			switch v := i.(type) {
			case time.Time:
				return "\x1b[90m" + v.UTC().Format(time.RFC3339) + "\x1b[0m"
			case string:
				// This case happens if zerolog already converted it to a string
				// Instead of parsing, we just show it
				return "\x1b[90m" + v + "\x1b[0m"
			default:
				return ""
			}
		},
		FormatMessage: func(i any) string {
			// if msg, ok := i.(string); ok {
			return "\x1b[0m" + i.(string) + "\x1b[0m"
			// }
			// return "Err printing message !!"
		},
	})

}
