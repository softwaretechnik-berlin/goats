package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

type zodArray struct {
	zodAnyType
}

var _ ZodArray = zodArray{}

func (a zodArray) Brand(brand string) ZodBranded {
	return chainBrand(a, brand)
}

func (a zodArray) Length(len uint) ZodArray {
	return zodArray{a.chain("length", ts.NumberLiteral(len))}
}

// TODO reconsider
func (a zodArray) DeclaredAs(name ts.Identifier) ZodType {
	return zodArray{zodAnyType{name}}
}
