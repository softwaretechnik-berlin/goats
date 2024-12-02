package reflective

import (
	"reflect"
	"strings"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
)

// Note: In order to have value equality regardless of which concrete type we use to get this type,
// we don't want to hold a reference to a concrete type.
// An alternative to the current approach of pre-invoking methods would be to only store the name and pkgPath and use a
// global map from that representation to a concrete type.
type typeConstructor struct {
	name     goinsp.TypeName
	pkgPath  goinsp.ImportPath
	string   string
	kind     reflect.Kind
	len      uint
	numField int
}

func (t typeConstructor) Name() goinsp.TypeName {
	return t.name
}

func (t typeConstructor) PkgPath() goinsp.ImportPath {
	return t.pkgPath
}

func (t typeConstructor) String() string {
	return t.string
}

func (t typeConstructor) Kind() reflect.Kind {
	return t.kind
}

func (t typeConstructor) Len() uint {
	if t.kind != reflect.Array {
		panic(t)
	}
	return t.len
}

func (t typeConstructor) NumField() int {
	if t.kind != reflect.Struct {
		panic(t)
	}
	return t.numField
}

func (t typeConstructor) WithoutTypeArguments() goinsp.GenType {
	return t
}

func (t typeConstructor) Comment() goinsp.PotentiallyUnavailable[goinsp.NoneWhenZero[string]] {
	return goinsp.PotentiallyUnavailable[goinsp.NoneWhenZero[string]]{}
}

func newGenericType(t typeAdaptor) goinsp.GenType {
	gen := typeConstructor{
		upToOpeningBrace(t.Name()),
		t.PkgPath(),
		upToOpeningBrace(t.String()),
		t.Kind(),
		0,
		0,
	}
	switch gen.kind {
	case reflect.Slice:
		// nothing
	case reflect.Array:
		gen.len = t.Len()
	case reflect.Struct:
		gen.numField = t.NumField()
	default:
		if !strings.ContainsRune(t.reflected.Name(), '[') {
			panic(t)
		}
	}
	return gen
}

func upToOpeningBrace[S ~string](s S) S {
	prefix, _, _ := strings.Cut(string(s), "[")
	return S(prefix)
}
