package zod

import (
	"slices"

	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/util"
)

type zodObject struct {
	zodAnyType

	shape []ShapeProperty
}

var _ ZodObject = zodObject{}

func (o zodObject) Brand(brand string) ZodBranded {
	return chainBrand(o, brand)
}

func (o zodObject) Extend(shape ...ShapeProperty) ZodObject {
	return zodObject{o.chain("extend", shapeTypeScript(shape)), append(slices.Clip(o.shape), shape...)}
}

func (o zodObject) Merge(schema ZodObject) ZodObject {
	return zodObject{o.chain("merge", schema.TypeScript()), append(slices.Clip(o.shape), schema.Shape()...)}
}

func (o zodObject) Shape() []ShapeProperty {
	return o.shape
}

// TODO reconsider
func (o zodObject) DeclaredAs(name ts.Identifier) ZodType {
	return zodObject{zodAnyType{name}, o.shape}
}

func shapeTypeScript(shape []ShapeProperty) ts.Source {
	return ts.Object(util.Map(shape, func(p ShapeProperty) ts.Property { return ts.Property{p.Name, p.Schema.TypeScript()} })...)
}
