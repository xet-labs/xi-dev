package util

import (
	"fmt"
)

func ToStr(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}