package goinsp

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

type reflectGenericallyMeaningfulInterface[TypeName ~string, ImportPath ~string, Length constraints.Integer] interface {
	Name() TypeName
	PkgPath() ImportPath
	String() string
	Kind() reflect.Kind
	Len() Length
	NumField() int
}

// reflectTypeInterface captures the commonalities between this package's Type and reflect.Type.
type reflectTypeInterface[TypeName ~string, ImportPath ~string, Length constraints.Integer, StructField reflectStructFieldInterface, Type any] interface {
	reflectGenericallyMeaningfulInterface[TypeName, ImportPath, Length]
	Elem() Type
	Field(i int) StructField
	Implements(u Type) bool
}

var _ reflectTypeInterface[string, string, int, reflect.StructField, reflect.Type] = reflect.Type(nil)
var _ reflectTypeInterface[TypeName, ImportPath, uint, StructField, Type] = Type(nil)

// GenType is an interface representing generalized types.
// This includes concrete types like `Bar` and `Foo[Bar]`, but also parameterized types like `Bar`.
type GenType interface {
	reflectGenericallyMeaningfulInterface[TypeName, ImportPath, uint]

	// WithoutTypeArguments returns a representation of this type without type arguments.
	//
	// Invoked on a concrete realization of a generic type, this will return a GenType representing the parameterized type.
	// Otherwise, it returns the type it was invoked on.
	WithoutTypeArguments() GenType

	Comment() PotentiallyUnavailable[NoneWhenZero[string]]
}

type Type interface {
	GenType
	reflectTypeInterface[TypeName, ImportPath, uint, StructField, Type]
	//PackageName() PotentiallyUnavailable[NoneWhenZero[PackageName]]
	Comment() PotentiallyUnavailable[NoneWhenZero[string]]
}
