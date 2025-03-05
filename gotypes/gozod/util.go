package gozod

import (
	"unicode"
	"unicode/utf8"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	firstRune, n := utf8.DecodeRuneInString(s)
	return string([]rune{unicode.ToUpper(firstRune)}) + s[n:]
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
