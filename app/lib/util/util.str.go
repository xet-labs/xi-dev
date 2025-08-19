package util

func PtrNotNilThen[T any](val string, ptr *T) string {
	if ptr != nil {
		return val
	}
	return ""
}