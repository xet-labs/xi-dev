package util

func MergeMapTo(m1, m2 map[string]any) map[string]any {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func MergeMaps(maps ...map[string]any) map[string]any {
	merged := map[string]any{}
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}