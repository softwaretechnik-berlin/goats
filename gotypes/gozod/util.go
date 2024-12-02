package gozod

import (
	"strings"
)

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func collectSlice[T any](cap int, element func(i int) (T, bool)) []T {
	slice := make([]T, 0, cap)
	for i := range cap {
		if value, ok := element(i); ok {
			slice = append(slice, value)
		}
	}
	return slice
}

func mapSlice[B, A any](as []A, f func(A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = f(a)
	}
	return bs
}
