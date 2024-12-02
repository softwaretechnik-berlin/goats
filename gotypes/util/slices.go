package util

import "unsafe"

// Map maps a slice from one type to another using an element mapping function.
func Map[B, A any](as []A, f func(A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = f(a)
	}
	return bs
}

// CaseElementsUnsafe
func CastElementsUnsafe[U, T any](values []T) []U {
	// Sanity check
	var in T
	var out U
	if unsafe.Sizeof(out) != unsafe.Sizeof(in) || unsafe.Alignof(out) != unsafe.Alignof(in) {
		panic("Incompatible types")
	}
	// Cast
	return *((*[]U)(unsafe.Pointer(&values)))
}

func Singleton[T any](values T) []T {
	return []T{values}
}

func Slice[T any](values ...T) []T {
	return values
}
