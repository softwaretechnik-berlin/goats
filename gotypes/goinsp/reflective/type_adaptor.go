package reflective

import (
	"reflect"
	"strings"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
)

type typeAdaptor struct{ reflected reflect.Type }

func (t typeAdaptor) PkgPath() goinsp.ImportPath {
	return goinsp.ImportPath(t.reflected.PkgPath())
}

func (t typeAdaptor) Comment() goinsp.PotentiallyUnavailable[goinsp.NoneWhenZero[string]] {
	return goinsp.PotentiallyUnavailable[goinsp.NoneWhenZero[string]]{}
}

func (t typeAdaptor) Name() goinsp.TypeName {
	return goinsp.TypeName(t.reflected.Name())
}

func (t typeAdaptor) String() string {
	return t.reflected.String()
}

func (t typeAdaptor) Implements(u goinsp.Type) bool {
	return t.reflected.Implements(u.(typeAdaptor).reflected)
}

func (t typeAdaptor) Kind() reflect.Kind {
	return t.reflected.Kind()
}

func (t typeAdaptor) Len() uint {
	return uint(t.reflected.Len())
}

func (t typeAdaptor) NumField() int {
	return t.reflected.NumField()
}

func (t typeAdaptor) Key() goinsp.Type {
	return typeAdaptor{t.reflected.Key()}
}

func (t typeAdaptor) Elem() goinsp.Type {
	return typeAdaptor{t.reflected.Elem()}
}

func (t typeAdaptor) Field(i int) goinsp.StructField {
	field := t.reflected.Field(i)
	return goinsp.NewStructField(field.Name, field.Tag, field.Anonymous, fieldAdaptor{field})
}

func (t typeAdaptor) WithoutTypeArguments() goinsp.GenType {
	if t.Kind() == reflect.Array || t.Kind() == reflect.Map || t.Kind() == reflect.Slice || strings.ContainsRune(t.reflected.Name(), '[') {
		return newGenericType(t)
	}
	return t
}

func Adapt(t reflect.Type) goinsp.Type {
	return typeAdaptor{t}
}

func TypeOf(i any) goinsp.Type {
	return Adapt(reflect.TypeOf(i))
}

func TypeFor[T any]() goinsp.Type {
	return Adapt(reflect.TypeFor[T]())
}
