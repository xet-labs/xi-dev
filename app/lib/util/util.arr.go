package util

// FallbackArg safely gets args[i] if within bounds and not the zero value.
// Otherwise it returns the provided fallback.
func ArrFallback[T comparable](args []T, i int, fallback T) T {
	if i < len(args) {
		var zero T
		if args[i] != zero {
			return args[i]
		}
	}
	return fallback
}