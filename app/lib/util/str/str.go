package str
import "strings"

type StrLib struct{}

func (s *StrLib) Fallback(str, fallback string) string {
	if str != "" {
		return str
	}
	return fallback
}

func (s *StrLib) Fallbacks(strs ...string) string {
	for _, v := range strs {
		if v == "" {
			continue
		}
		return v
	}
	return ""
}

func (s *StrLib) NotEmptyThen(val, str string) string {
	if str != "" {
		return val
	}
	return ""
}

func (u *StrLib) QuoteSlice(items []string) string {
	quoted := make([]string, len(items))
	for i, v := range items {
		quoted[i] = "'" + v + "'"
	}
	return strings.Join(quoted, ", ")
}
