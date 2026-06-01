package anyx

// Coalesce returns the first non-zero value.
func Coalesce[T comparable](vs ...T) T {
	var zero T
	for _, v := range vs {
		if v != zero {
			return v
		}
	}
	return zero
}
