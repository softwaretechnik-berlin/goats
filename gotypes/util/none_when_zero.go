package util

import (
	"encoding/json"
	"fmt"
)

// NoneWhenZero explicitly signals that its value might be https://go.dev/ref/spec#The_zero_value for that type.
//
// E.g., if we declare `userID UUID`, we typically mean that `userID` always has a valid id, even though it's
// technically possible for me to give it the zero value UUID, which doesn't actually identify a user.
// If we intend for this value to potentially be empty, we can make this explicit by declaring `userID NoneWhenZero[UUID]`.
// This signals to others that this value could be empty, and gives the additional behaviour described below.
//
// NoneWhenZero is equipped with utility methods for working with the value.
//
// When the value is the zero value, the NoneWhenZero is serialized as `null` in JSON.
//
// This type is appropriate to use with relatively small data types.
// Larger data types should use [Optional].
type NoneWhenZero[A comparable] struct {
	V A `tsgen:",value,nullable"`
}

var _ json.Marshaler = (*NoneWhenZero[any])(nil)
var _ json.Unmarshaler = (*NoneWhenZero[any])(nil)

func AsNoneWhenZero[A comparable](value A) NoneWhenZero[A] {
	return NoneWhenZero[A]{value}
}

func AsNoneWhenZeros[T comparable](values []T) []NoneWhenZero[T] {
	return CastElementsUnsafe[NoneWhenZero[T]](values)
}

func NoneZero[A comparable]() NoneWhenZero[A] {
	return NoneWhenZero[A]{}
}

func NoneWhenZeroFromOptional[A comparable](o Optional[A]) NoneWhenZero[A] {
	if o.HasValue {
		return AsNoneWhenZero(o.V)
	}
	return NoneZero[A]()
}

func (n NoneWhenZero[A]) IsNone() bool {
	return IsZero(n)
}

func (n NoneWhenZero[A]) Is(value A) bool {
	return n.V == value
}

func (n NoneWhenZero[A]) HasValue() bool {
	return !n.IsNone()
}

func (n NoneWhenZero[A]) All(predicate func(A) bool) bool {
	return MapNoneWhenZeroWithDefault(n, true, predicate)
}

func (n NoneWhenZero[A]) Exists(predicate func(A) bool) bool {
	return MapNoneWhenZeroWithDefault(n, false, predicate)
}

func (n NoneWhenZero[A]) IfPresent(whenPresent func(A)) {
	if n.HasValue() {
		whenPresent(n.V)
	}
}

func (n NoneWhenZero[A]) Filter(predicate func(A) bool) NoneWhenZero[A] {
	if predicate(n.V) {
		return n
	}
	return NoneZero[A]()
}

// MarshalJSON implements json.Marshaler, representing zero values as `null`
// and non-zero values as their normal JSON representation (without any wrapping struct).
func (n NoneWhenZero[A]) MarshalJSON() ([]byte, error) {
	if n.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(n.V)
}

// MarshalJSON implements json.Unmarshaler, unmarshalling the representation given by MarshalJSON.
func (n *NoneWhenZero[A]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		*n = NoneZero[A]()
		return nil
	}
	return json.Unmarshal(bytes, &n.V)
}

func (n NoneWhenZero[A]) Get() (A, bool) {
	return n.V, n.HasValue()
}

func (n NoneWhenZero[A]) GetOrElse(fallback A) A {
	if n.HasValue() {
		return n.V
	}
	return fallback
}

func (n NoneWhenZero[A]) GetOrElseFunc(fallback func() A) A {
	if n.HasValue() {
		return n.V
	}
	return fallback()
}

func (n NoneWhenZero[A]) MustGet() A {
	if n.IsNone() {
		panic(fmt.Sprintf("Unexpectedly getting zero value of type %T", n.V))
	}
	return n.V
}

func (n NoneWhenZero[A]) Or(other NoneWhenZero[A]) NoneWhenZero[A] {
	if n.IsNone() {
		return other
	}
	return n
}

func (n NoneWhenZero[A]) ToOption() Optional[A] {
	return TupleAsOptional(n.Get())
}

func (n NoneWhenZero[A]) ToPtr() *A {
	if n.IsNone() {
		return nil
	}
	return &n.V
}

func (n NoneWhenZero[A]) ToSlice() []A {
	return MapNoneWhenZeroWithDefault(n, nil, Singleton[A])
}

func FoldNoneWhenZero[B any, A comparable](value NoneWhenZero[A], initial B, f func(B, A) B) B {
	if value.IsNone() {
		return initial
	}
	return f(initial, value.V)
}

func MapNoneWhenZero[B, A comparable](value NoneWhenZero[A], f func(A) B) NoneWhenZero[B] {
	return AsNoneWhenZero(MapNoneWhenZeroWithDefault(value, Zero[B](), f))
}

func MapNoneWhenZeroWithDefault[B any, A comparable](value NoneWhenZero[A], whenNone B, f func(A) B) B {
	if value.IsNone() {
		return whenNone
	}
	return f(value.V)
}

func MapNoneWhenZeroWithDefaultFunc[B any, A comparable](value NoneWhenZero[A], whenNone func() B, f func(A) B) B {
	if value.IsNone() {
		return whenNone()
	}
	return f(value.V)
}

func FlatMapNoneWhenZero[B, A comparable](value NoneWhenZero[A], f func(A) NoneWhenZero[B]) NoneWhenZero[B] {
	return MapNoneWhenZeroWithDefaultFunc(value, NoneZero[B], f)
}

// Combine applies a function to the values of a and b if both have values.
// It returns NoneZero if one of a or b has no value.
func Combine[C, B, A comparable](a NoneWhenZero[A], b NoneWhenZero[B], f func(A, B) C) NoneWhenZero[C] {
	return FlatMapNoneWhenZero(a, func(a A) NoneWhenZero[C] {
		return MapNoneWhenZero(b, func(b B) C {
			return f(a, b)
		})
	})
}

func MapNoneWhenZeroToOptional[B any, A comparable](value NoneWhenZero[A], f func(A) B) Optional[B] {
	return MapNoneWhenZeroWithDefault(value, Optional[B]{}, func(a A) Optional[B] { return AsOptional(f(a)) })
}
