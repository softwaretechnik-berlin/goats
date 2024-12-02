package goinsp

import (
	"reflect"
)

type reflectStructFieldInterface interface {
	IsExported() bool
}

var _ reflectStructFieldInterface = reflect.StructField{}
var _ reflectStructFieldInterface = StructField{}

type StructField struct {
	Name      string
	Tag       reflect.StructTag
	Anonymous bool
	impl      StructFieldImpl
}

func (f StructField) IsExported() bool {
	return f.impl.IsExported()
}

func (f StructField) Type() Type {
	return f.impl.Type()
}

func NewStructField(name string, tag reflect.StructTag, anonymous bool, impl StructFieldImpl) StructField {
	return StructField{name, tag, anonymous, impl}
}

type StructFieldImpl interface {
	IsExported() bool
	Type() Type
}
