package str

import (
	"sort"
	"strings"
)

type StrLib struct{}

func (s *StrLib) Unique(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func (s *StrLib) UniqueSort(in []string) []string {
	out := s.Unique(in)   // reuse unique logic
	sort.Strings(out)     // sort in place
	return out
}

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
