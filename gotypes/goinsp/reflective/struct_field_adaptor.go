package reflective

import (
	"reflect"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
)

type fieldAdaptor struct {
	reflected reflect.StructField
}

func (f fieldAdaptor) IsExported() bool {
	return f.reflected.IsExported()
}

func (f fieldAdaptor) Type() goinsp.Type {
	return typeAdaptor{f.reflected.Type}
}
