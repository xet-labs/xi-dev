package util

import (
	"strings"
)

type utilLib struct{}

var Util = &utilLib{}

func (u *utilLib) QuoteSlice(items []string) string {
	quoted := make([]string, len(items))
	for i, v := range items {
		quoted[i] = "'" + v + "'"
	}
	return strings.Join(quoted, ", ")
}
