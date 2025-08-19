package str
import "strings"

type StrLib struct{}

func (s *StrLib) Fallback(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}

func (s *StrLib) Fallbacks(vals ...string) string {
	for _, v := range vals {
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
