package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

var z = ts.ImportedName("zod", "z")

func zTypeFunc(name ts.Identifier, args ...ts.Source) zodAnyType {
	return zodAnyType{ts.InvokeMethod(z, name, args...)}
}
