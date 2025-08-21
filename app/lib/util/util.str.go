package util

func StrIfPtrNotNil[T any](val string, ptr *T) string {
	if ptr != nil {
		return val
	}
	return ""
}

func IfPtrNotNil[T any](val T, ptr *T) T {
	var zero T
	if ptr != nil {
		return val
	}
	return zero
}

