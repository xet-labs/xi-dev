package	maps

type MapsLib struct{}

func (s *MapsLib) AddIfNotEmpty(m map[string]any, key string, val string) {
	if val != "" {
		m[key] = val
	}
}

func (s *MapsLib) AddIfNotEmptySlice(m map[string]any, key string, vals []string) {
	if len(vals) > 0 {
		m[key] = vals
	}
}