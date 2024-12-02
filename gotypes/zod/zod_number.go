package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

type zodNumber struct {
	zodAnyType
	int, nonNegative bool
}

var _ ZodNumber = zodNumber{}

func (n zodNumber) Brand(brand string) ZodBranded {
	return chainBrand(n, brand)
}

func (n zodNumber) Int() ZodNumber {
	return zodNumber{n.chain("int"), true, n.nonNegative}
}

func (n zodNumber) NonNegative() ZodNumber {
	return zodNumber{n.chain("nonnegative"), n.int, true}
}

func (n zodNumber) IsInt() bool {
	return n.int
}

func (n zodNumber) IsNonNegative() bool {
	return n.nonNegative
}

// TODO reconsider
func (n zodNumber) DeclaredAs(name ts.Identifier) ZodType {
	return zodNumber{zodAnyType{name}, n.int, n.nonNegative}
}
