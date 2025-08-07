package htmlfn

import "time"

func IsSlice(v any) bool {
	_, ok := v.([]string)
	return ok
}

func Len(arr []string) int {
	return len(arr)
}

func Slice(arr []string, start, end int) []string {
	if start >= len(arr) {
		return []string{}
	}
	if end > len(arr) {
		end = len(arr)
	}
	return arr[start:end]
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}
