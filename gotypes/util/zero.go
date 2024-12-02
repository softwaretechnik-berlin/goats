package util

// Zero returns https://go.dev/ref/spec#The_zero_value of type A.
func Zero[T any]() (zero T) { return }

// IsZero returns true iff the value is https://go.dev/ref/spec#The_zero_value of type A.
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

// IsNonZero returns true iff the value is _not_ https://go.dev/ref/spec#The_zero_value of type A.
func IsNonZero[T comparable](value T) bool {
	return !IsZero(value)
}
