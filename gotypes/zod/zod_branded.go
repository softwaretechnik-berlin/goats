package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

type zodBranded struct {
	zodAnyType
	wrapped ZodType
	brand   string
}

func (b zodBranded) Brand(brand string) ZodBranded {
	return chainBrand(b, brand)
}

func (b zodBranded) Unwrap() ZodType {
	return b.wrapped
}

// TODO reconsider
func (b zodBranded) DeclaredAs(name ts.Identifier) ZodType {
	wrapped := b.wrapped.DeclaredAs(name)
	return zodBranded{zodAnyType{wrapped.TypeScript()}, wrapped, b.brand}
}

func chainBrand(t ZodType, brand string) zodBranded {
	return zodBranded{zodAnyType{ts.InvokeMethod(t.TypeScript(), "brand", ts.StringLiteral(brand))}, t, brand}
}
