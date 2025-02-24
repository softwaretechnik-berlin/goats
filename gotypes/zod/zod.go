package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/util"
)

type ZodType interface {
	Brand(brand string) ZodBranded
	Nullable() ZodNullable
	Optional() ZodOptional
	Parse(str ts.Source) ts.Source
	Parsef(format string, a ...ts.Source) ts.Source
	Pipe(target ZodType) ZodType
	Transform(transform ts.Source) ZodType
	Transformf(format string, a ...ts.Source) ZodType

	DeclaredAs(name ts.Identifier) ZodType
	TypeScript() ts.Source
}

var _ ZodType = ZodArray(nil)
var _ ZodType = ZodBranded(nil)
var _ ZodType = ZodNumber(nil)
var _ ZodType = ZodNullable(nil)
var _ ZodType = ZodObject(nil)
var _ ZodType = ZodString(nil)

type ZodArray interface {
	ZodType
	Length(len uint) ZodArray
}

type ZodBranded interface {
	ZodType

	Unwrap() ZodType
}

type ZodNumber interface {
	ZodType

	Int() ZodNumber
	NonNegative() ZodNumber

	IsInt() bool
	IsNonNegative() bool
}

type ZodNullable interface {
	ZodType

	Unwrap() ZodType
}

type ZodOptional interface {
	ZodType

	Unwrap() ZodType
}

type ShapeProperty = struct {
	Name   string
	Schema ZodType
}

type ZodObject interface {
	ZodType
	Extend(shape ...ShapeProperty) ZodObject
	Merge(schema ZodObject) ZodObject

	Shape() []ShapeProperty
}

type ZodString interface {
	ZodType
	UUID() ZodString
}

func Any() ZodType {
	return zTypeFunc("any")
}

func Array(schema ZodType) ZodArray {
	return zodArray{zTypeFunc("array", schema.TypeScript())}
}

func Boolean() ZodType {
	return zTypeFunc("boolean")
}

func Literal(value string) ZodType {
	return zTypeFunc("literal", ts.StringLiteral(value))
}

func Nullable(t ZodType) ZodNullable {
	return zodNullable{zTypeFunc("nullable", t.TypeScript()), t}
}

// EnsureNullable is a convenience method that calls Nullable on the given schema unless it is sure that doing so will
// produce a schema with identical parse semantics.
//
// E.g. suppose we have a property nullable property of type A, and the representation of type is already `z.string().nullable()`;
// EnsureNullable leaves the schema unchanged in this case (instead of making it `z.string().nullable().nullable()`),
// while e.g. making `z.number()` into `z.number().nullable()`.
//
// Note that this method's concept of "nullable" is the same as that of Zod's `.nullable()`, which is not the same as
// Zod's `.isNullable()`'s concept of "nullable". The former ensures that both the input and output types are nullable,
// that a null input will produce a null output and that the wrapped schema will not be used when parsing null; the
// latter only tests whether the schema accepts null as a valid input value.
func EnsureNullable(t ZodType) ZodType {
	if _, isNullable := t.(zodNullable); isNullable {
		return t
	}
	return t.Nullable()
}

// Enum type with the given permissible values
func Enum(values ...string) ZodType {
	return zTypeFunc("enum", ts.Array(util.Map(values, ts.StringLiteral)...))
}

// StripNullable strips away any known nullable wrappers and returns a bool indicating whether nullability was stripped away.
func StripNullable(t ZodType) (ZodType, bool) {
	if t, isNullable := t.(zodNullable); isNullable {
		t, _ := StripNullable(t.Unwrap())
		return t, true
	}
	return t, false
}

func Number() ZodNumber {
	return zodNumber{zTypeFunc("number"), false, false}
}

func Object(shape ...ShapeProperty) ZodObject {
	return zodObject{zTypeFunc("object", shapeTypeScript(shape)), shape}
}

func Record(keySchema, valueType ZodType) ZodType {
	return zodArray{zTypeFunc("record", keySchema.TypeScript(), valueType.TypeScript())}
}

func String() ZodString {
	return zodString{zTypeFunc("string")}
}

func Union(types ...ZodType) ZodType {
	return zTypeFunc("union", ts.Array(util.Map(types, ZodType.TypeScript)...))
}

func DiscriminatedUnion(discriminator string, types ...ZodType) ZodType {
	return zTypeFunc("discriminatedUnion", ts.StringLiteral(discriminator), ts.Array(util.Map(types, ZodType.TypeScript)...))
}

// ZodZtypeExpr is an escape hatch to create a ZodType from an arbitrary ts.Source.
func ZodTypeExpr(expr ts.Source) ZodType {
	return zodAnyType{expr}
}

// TODO reconsider
type SchemaAndTypeDeclaration struct {
	comment    string
	identifier ts.Identifier
	schema     ZodType
}

func (d SchemaAndTypeDeclaration) Identifier() ts.Identifier { return d.identifier }

func (d SchemaAndTypeDeclaration) TypeScript() ts.Source {
	return ts.Statements(
		ts.DocComment(d.comment),
		ts.Sourcef(`export const %s = %s;`, d.identifier, d.schema.TypeScript()),
		ts.Sourcef(`export type %s = %s.infer<typeof %s>;`, d.identifier, z, d.identifier),
	)
}

// NewSchemaAndTypeDeclaration TODO
func NewSchemaAndTypeDeclaration(comment string, name ts.Identifier, schema ZodType) SchemaAndTypeDeclaration {
	return SchemaAndTypeDeclaration{comment, name, schema}
}
