package zod_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/softwaretechnik-berlin/goats/gotypes/zod"
)

var zImport = `import { z } from "zod";`

func TestZodTypeScript(t *testing.T) {
	assertTypeScriptRepresentationOf(t, zod.Any(), zImport, `z.any()`)
	assertTypeScriptRepresentationOf(t, zod.Array(zod.String()), zImport, `z.array(z.string())`)
	assertTypeScriptRepresentationOf(t, zod.Boolean(), zImport, `z.boolean()`)
	assertTypeScriptRepresentationOf(t, zod.Literal("foo"), zImport, `z.literal("foo")`)
	assertTypeScriptRepresentationOf(t, zod.Number(), zImport, `z.number()`)
	assertTypeScriptRepresentationOf(t, zod.Number().Int(), zImport, `z.number().int()`)
	assertTypeScriptRepresentationOf(t, zod.Object(), zImport, `z.object({})`)
	assertTypeScriptRepresentationOf(t, zod.Object(zod.ShapeProperty{"foo", zod.String()}), zImport, `z.object({ foo: z.string() })`)
	assertTypeScriptRepresentationOf(t, zod.Object(
		zod.ShapeProperty{"foo", zod.String()},
		zod.ShapeProperty{"bar", zod.Number()},
	), zImport, `z.object({
    foo: z.string(),
    bar: z.number(),
})`)
	assertTypeScriptRepresentationOf(t, zod.String(), zImport, `z.string()`)
	assertTypeScriptRepresentationOf(t, zod.String().UUID(), zImport, `z.string().uuid()`)
	assertTypeScriptRepresentationOf(t, zod.Union(), zImport, `z.union([])`)
	assertTypeScriptRepresentationOf(t, zod.Union(zod.String()), zImport, `z.union([z.string()])`)
	assertTypeScriptRepresentationOf(t, zod.Union(zod.String(), zod.Number()), zImport, `z.union([
    z.string(),
    z.number(),
])`)
	assertTypeScriptRepresentationOf(t, zod.DiscriminatedUnion("foo",
		zod.Object(zod.ShapeProperty{"foo", zod.String()}),
		zod.Object(zod.ShapeProperty{"foo", zod.Number()}),
	), zImport, `z.discriminatedUnion("foo", [
    z.object({ foo: z.string() }),
    z.object({ foo: z.number() }),
])`)

	assertTypeScriptRepresentationOf(t, zod.String().Nullable(), zImport, `z.string().nullable()`)
	assertTypeScriptRepresentationOf(t, zod.Nullable(zod.String()), zImport, `z.nullable(z.string())`)
	assertTypeScriptRepresentationOf(t, zod.EnsureNullable(zod.String()), zImport, `z.string().nullable()`)

	assertTypeScriptRepresentationOf(t, zod.String().Nullable().Nullable(), zImport, `z.string().nullable().nullable()`)
	assertTypeScriptRepresentationOf(t, zod.Nullable(zod.String()).Nullable(), zImport, `z.nullable(z.string()).nullable()`)
	assertTypeScriptRepresentationOf(t, zod.EnsureNullable(zod.String()).Nullable(), zImport, `z.string().nullable().nullable()`)

	assertTypeScriptRepresentationOf(t, zod.Nullable(zod.String().Nullable()), zImport, `z.nullable(z.string().nullable())`)
	assertTypeScriptRepresentationOf(t, zod.Nullable(zod.Nullable(zod.String())), zImport, `z.nullable(z.nullable(z.string()))`)
	assertTypeScriptRepresentationOf(t, zod.Nullable(zod.EnsureNullable(zod.String())), zImport, `z.nullable(z.string().nullable())`)

	assertTypeScriptRepresentationOf(t, zod.EnsureNullable(zod.String().Nullable()), zImport, `z.string().nullable()`)
	assertTypeScriptRepresentationOf(t, zod.EnsureNullable(zod.Nullable(zod.String())), zImport, `z.nullable(z.string())`)
	assertTypeScriptRepresentationOf(t, zod.EnsureNullable(zod.EnsureNullable(zod.String())), zImport, `z.string().nullable()`)
}

func assertTypeScriptRepresentationOf(t *testing.T, schema zod.ZodType, expectedImports string, expectedCode string) {
	code := schema.TypeScript()
	//assert.Equal(t, expectedCode, code.WithoutImports())
	assert.Equal(t, expectedImports+"\n\n"+expectedCode, code.String())
}
