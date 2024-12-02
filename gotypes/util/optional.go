package util

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/samber/lo"
)

type Optional[A any] struct {
	HasValue bool
	V        A `tsgen:",value,nullable"`
}

var _ json.Marshaler = (*Optional[any])(nil)
var _ json.Unmarshaler = (*Optional[any])(nil)

func AsOptional[A any](value A) Optional[A] {
	return Optional[A]{true, value}
}

func None[A any]() Optional[A] {
	return Optional[A]{}
}

func TupleAsOptional[A any](value A, hasValue bool) Optional[A] {
	return Optional[A]{hasValue, value}
}

func ZeroableAsOptional[A yaml.IsZeroer](value A) Optional[A] {
	return Optional[A]{!value.IsZero(), value}
}

func (o Optional[A]) All(predicate func(A) bool) bool {
	return MapOptionalWithDefault(o, true, predicate)
}

func (o Optional[A]) Exists(predicate func(A) bool) bool {
	return MapOptionalWithDefault(o, false, predicate)
}

func (o Optional[A]) IfPresent(whenPresent func(value A)) {
	if o.HasValue {
		whenPresent(o.V)
	}
}

func (o Optional[A]) Filter(predicate func(A) bool) Optional[A] {
	if o.HasValue && predicate(o.V) {
		return o
	}
	return Optional[A]{}
}

// MarshalJSON implements json.Marshaler, representing zero values as `null`
// and non-zero values as their normal JSON representation (without any wrapping struct).
func (o Optional[A]) MarshalJSON() ([]byte, error) {
	if !o.HasValue {
		return []byte("null"), nil
	}
	return json.Marshal(o.V)
}

// Value implements driver.Valuer
func (o Optional[A]) Value() (value driver.Value, err error) {
	if o.IsNone() {
		return nil, nil
	}
	value = o.V
	for {
		valuer, ok := value.(driver.Valuer)
		if !ok {
			break
		}
		value, err = valuer.Value()
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalJSON implements json.Unmarshaler, unmarshalling the representation given by MarshalJSON.
func (o *Optional[A]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		*o = Optional[A]{}
		return nil
	}
	err := json.Unmarshal(bytes, &o.V)
	o.HasValue = err == nil
	return err
}

func (o Optional[A]) Get() (A, bool) {
	return o.V, o.HasValue
}

func (o Optional[A]) GetOrElse(orElse A) A {
	return lo.Ternary(o.HasValue, o.V, orElse)
}

func (o Optional[A]) GetOrElseDefault() (a A) {
	return lo.Ternary(o.HasValue, o.V, a)
}

func (o Optional[A]) GetOrElseFunc(orElse func() A) A {
	return lo.Ternary(o.HasValue, o.V, orElse())
}

func (o Optional[A]) IsNone() bool {
	return !o.HasValue
}

func (o Optional[A]) MustGet() A {
	if o.HasValue {
		return o.V
	}
	panic(fmt.Sprintf("Cannot get optional value of type %T", o.V))
}

func (o Optional[A]) Or(other Optional[A]) Optional[A] {
	if o.HasValue {
		return o
	}
	return other
}

func FoldOptional[B, A any](o Optional[A], initial B, f func(B, A) B) B {
	if o.HasValue {
		return f(initial, o.V)
	}
	return initial
}

func MapOptional[B, A any](o Optional[A], f func(A) B) Optional[B] {
	return MapOptionalWithDefault(o, Optional[B]{}, func(a A) Optional[B] { return AsOptional(f(a)) })
}

func MapOptionalWithDefault[B, A any](o Optional[A], whenNone B, f func(A) B) B {
	if o.HasValue {
		return f(o.V)
	}
	return whenNone
}

func MapOptionalWithDefaultFunc[B, A any](o Optional[A], whenNone func() B, f func(A) B) B {
	if o.HasValue {
		return f(o.V)
	}
	return whenNone()
}

func MapOptionalToNoneWhenZero[B comparable, A any](o Optional[A], f func(A) B) NoneWhenZero[B] {
	return MapOptionalWithDefault(o, NoneZero[B](), func(a A) NoneWhenZero[B] { return AsNoneWhenZero(f(a)) })
}

func ToNoneWhenZero[A comparable](o Optional[A]) NoneWhenZero[A] {
	return AsNoneWhenZero(o.V)
}
