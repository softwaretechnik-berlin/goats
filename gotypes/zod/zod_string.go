package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

type zodString struct {
	zodAnyType
}

var _ ZodString = zodString{}

func (s zodString) UUID() ZodString {
	return zodString{s.chain("uuid")}
}

// TODO reconsider
func (t zodString) DeclaredAs(name ts.Identifier) ZodType {
	return zodString{zodAnyType{name}}
}
