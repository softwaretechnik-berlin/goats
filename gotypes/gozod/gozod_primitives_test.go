package gozod_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/parsing/comments"
	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
	"github.com/softwaretechnik-berlin/goats/gotypes/gozod"
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/zod"
)

const z = `import { z } from "zod";`

func TestGenerateForPrimitiveTypes(t *testing.T) {
	assertSimpleSchemaFor[any](t, z,
		`z.any()`,
		examples[any]{
			simpleExample[any](true, `true`),
			simpleExample[any](false, `false`),
			simpleExample[any](float64(1), `1`),
			simpleExample[any]("foo", `"foo"`),
			simpleExample[any](nil, `null`),
		},
		rejects{},
	)
	assertSimpleSchemaFor[bool](t, z,
		`z.boolean()`,
		examples[bool]{
			simpleExample(true, `true`),
			simpleExample(false, `false`),
		},
		rejects{`null`, `undefined`},
	)
	assertSimpleSchemaFor[float32](t, z,
		`z.number()`,
		examples[float32]{
			simpleExample[float32](0, `0`),
			simpleExample[float32](-0.5, `-0.5`),
			simpleExample[float32](12345.678, `12345.678`),
		},
		rejects{`null`, `undefined`},
	)
	assertSimpleSchemaFor[float64](t, z,
		`z.number()`,
		examples[float64]{
			simpleExample(float64(0), `0`),
			simpleExample(-0.5, `-0.5`),
			simpleExample(12345.6789, `12345.6789`),
		},
		rejects{`null`, `undefined`},
	)
	assertSimpleSchemaFor[int](t, z,
		`z.number().int()`,
		examples[int]{
			simpleExample(0, `0`),
			simpleExample(-12345, `-12345`),
			simpleExample(12345, `12345`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSimpleSchemaFor[int8](t, z,
		`z.number().int()`,
		examples[int8]{
			simpleExample[int8](0, `0`),
			simpleExample[int8](-128, `-128`),
			simpleExample[int8](127, `127`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSimpleSchemaFor[int16](t, z,
		`z.number().int()`,
		examples[int16]{
			simpleExample[int16](0, `0`),
			simpleExample[int16](-12345, `-12345`),
			simpleExample[int16](12345, `12345`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSimpleSchemaFor[int32](t, z,
		`z.number().int()`,
		examples[int32]{
			simpleExample[int32](0, `0`),
			simpleExample[int32](-12345, `-12345`),
			simpleExample[int32](12345, `12345`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSimpleSchemaFor[int64](t, z,
		`z.number().int()`,
		examples[int64]{
			simpleExample[int64](0, `0`),
			simpleExample[int64](-12345, `-12345`),
			simpleExample[int64](12345, `12345`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSimpleSchemaFor[uint](t, z,
		`z.number().nonnegative().int()`,
		examples[uint]{
			simpleExample[uint](0, `0`),
			simpleExample[uint](12345, `12345`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[uint8](t, z,
		`z.number().nonnegative().int()`,
		examples[uint8]{
			simpleExample[uint8](0, `0`),
			simpleExample[uint8](127, `127`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[uint16](t, z,
		`z.number().nonnegative().int()`,
		examples[uint16]{
			simpleExample[uint16](0, `0`),
			simpleExample[uint16](12345, `12345`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[uint32](t, z,
		`z.number().nonnegative().int()`,
		examples[uint32]{
			simpleExample[uint32](0, `0`),
			simpleExample[uint32](12345, `12345`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[uint64](t, z,
		`z.number().nonnegative().int()`,
		examples[uint64]{
			simpleExample[uint64](0, `0`),
			simpleExample[uint64](12345, `12345`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[uintptr](t, z,
		`z.number().nonnegative().int()`,
		examples[uintptr]{
			simpleExample[uintptr](0, `0`),
			simpleExample[uintptr](12345, `12345`),
		},
		rejects{`null`, `undefined`, `-1`, `0.5`},
	)
	assertSimpleSchemaFor[string](t, z,
		`z.string()`,
		examples[string]{
			simpleExample("", `""`),
			simpleExample("foo", `"foo"`),
		},
		rejects{`null`, `undefined`, `1`},
	)
}

func TestGenerateForSliceTypes(t *testing.T) {
	// Go encodes nil slices to JSON null and empty non-nil slices to empty JSON arrays/strings.
	// So in general, the most accurate representation for of the type of a JSON-marshalled slice is a nullable type.
	// However just as we can usually treat a nil slice in Go the same as an empty slice, on the JavaScript/TypeScript
	// side we usually want to treat a null as semantically equivalent to an empty array/string. So by default, slices
	// schemas will accept null but transform it into an empty value.
	//
	// TODO Untransformed nullable types that preserve the null vs. empty distinction can be enabled with
	// `WithNullableSlices()`.
	//
	// TODO Finally, if you're willing to assert that you will always populate these types with a non-nil value, you can
	// use `WithNonNullSlices()` to suppress the nullability completely and reject null values.
	assertSimpleSchemaFor[[]string](t, z,
		`z.array(z.string()).nullable().transform(a => a ?? [])`,
		examples[[]string]{
			simpleExample[[]string](nil, `null`),
			simpleExample([]string{}, `[]`),
			simpleExample([]string{"foo"}, `["foo"]`),
			simpleExample([]string{"foo", "bar"}, `["foo","bar"]`),
		},
		rejects{`undefined`, `"foo"`, `[0]`},
	)
	assertSimpleSchemaFor[[][]string](t, z,
		`z.array(z.array(z.string()).nullable().transform(a => a ?? [])).nullable().transform(a => a ?? [])`,
		examples[[][]string]{
			simpleExample[[][]string](nil, `null`),
			simpleExample([][]string{}, `[]`),
			simpleExample([][]string{nil, {}, {"foo"}}, `[null,[],["foo"]]`),
			simpleExample([][]string{{"foo", "bar"}, {"baz"}}, `[["foo","bar"],["baz"]]`),
		},
		rejects{`undefined`, `"foo"`, `["foo"]`, `[[0]]`},
	)
	assertSimpleSchemaFor[[]int](t, z,
		`z.array(z.number().int()).nullable().transform(a => a ?? [])`,
		examples[[]int]{
			simpleExample[[]int](nil, `null`),
			simpleExample([]int{}, `[]`),
			simpleExample([]int{0}, `[0]`),
			simpleExample([]int{0, 1}, `[0,1]`),
		},
		rejects{`undefined`, `0`, `[0.5]`},
	)
	assertSimpleSchemaFor[[][]int](t, z,
		`z.array(z.array(z.number().int()).nullable().transform(a => a ?? [])).nullable().transform(a => a ?? [])`,
		examples[[][]int]{
			simpleExample[[][]int](nil, `null`),
			simpleExample([][]int{}, `[]`),
			simpleExample([][]int{nil, {}, {0}}, `[null,[],[0]]`),
			simpleExample([][]int{{0, 1}, {2}}, `[[0,1],[2]]`),
		},
		rejects{`undefined`, `0`, `[0]`, `[[0.5]]`},
	)
	// Go encodes non-nil slices which have an element type of kind uint8 that implements neither json.Marshaler nor encoding.TextMarshaler as base64-encoded strings.
	assertSimpleSchemaFor[[]byte](t, z,
		`z.string().nullable().transform(a => a ?? "")`,
		examples[[]byte]{
			simpleExample[[]byte](nil, `null`),
			simpleExample([]byte{}, `""`),
			simpleExample([]byte{0}, `"AA=="`),
			simpleExample([]byte{0, 1}, `"AAE="`),
			simpleExample([]byte{0, 1, 2}, `"AAEC"`),
			simpleExample([]byte{0, 1, 2, 3}, `"AAECAw=="`),
		},
		rejects{`undefined`, `0`, `[0.5]`},
	)
	// TODO should have an option to opt for a base64 decode instead of getting the string (or even make that the default)?
	assertSimpleSchemaFor[[][]byte](t, z,
		`z.array(z.string().nullable().transform(a => a ?? "")).nullable().transform(a => a ?? [])`,
		examples[[][]byte]{
			simpleExample[[][]byte](nil, `null`),
			simpleExample([][]byte{}, `[]`),
			simpleExample([][]byte{nil, {}, {0}}, `[null,"","AA=="]`),
			simpleExample([][]byte{{0, 1}, {2}}, `["AAE=","Ag=="]`),
		},
		rejects{`undefined`, `0`, `[0]`, `[[0.5]]`},
	)
}

func TestGenerateForPointerTypes(t *testing.T) {
	// Pointer types are nullable, because they include the nil value, which becomes null in JSON.
	// TODO should we have a mechanism to make it possible to declare away the nullability of pointers?
	// TODO should we have a mechanism to make it possible to transform null values to a zero value?
	assertSimpleSchemaFor[*string](t, z,
		`z.string().nullable()`,
		examples[*string]{
			simpleExample[*string](nil, `null`),
			simpleExample(ptr(""), `""`),
			simpleExample(ptr("foo"), `"foo"`),
		},
		rejects{`undefined`, `1`},
	)
	// avoid double `.nullable()`
	assertSimpleSchemaFor[**string](t, z,
		`z.string().nullable()`,
		examples[**string]{
			{nil, []**string{ptr[*string](nil)}, `null`},
			simpleExample(ptr(ptr("")), `""`),
			simpleExample(ptr(ptr("foo")), `"foo"`),
		},
		rejects{`undefined`, `1`},
	)
	assertSimpleSchemaFor[*[]byte](t, z,
		`z.string().nullable().transform(a => a ?? "").nullable()`,
		examples[*[]byte]{
			{nil, []*[]byte{ptr[[]byte](nil)}, `null`},
			simpleExample(ptr([]byte{}), `""`),
			simpleExample(ptr([]byte{0}), `"AA=="`),
			simpleExample(ptr([]byte{0, 1}), `"AAE="`),
			simpleExample(ptr([]byte{0, 1, 2}), `"AAEC"`),
			simpleExample(ptr([]byte{0, 1, 2, 3}), `"AAECAw=="`),
		},
		rejects{`undefined`, `0`, `[0.5]`},
	)
	assertSimpleSchemaFor[*[]string](t, z,
		`z.array(z.string()).nullable().transform(a => a ?? []).nullable()`,
		examples[*[]string]{
			{nil, []*[]string{ptr[[]string](nil)}, `null`},
			simpleExample(ptr([]string{}), `[]`),
			simpleExample(ptr([]string{"foo"}), `["foo"]`),
			simpleExample(ptr([]string{"foo", "bar"}), `["foo","bar"]`),
		},
		rejects{`undefined`, `"foo"`, `[0]`},
	)
}

type (
	newTypeS      string
	newTypeB      bool
	newTypeI      int
	newTypeF      float32
	newTypeStruct struct{ A string }
)

func TestGenerateForNamedTypes(t *testing.T) {
	// TODO options to configure branding

	// named simple types are branded
	assertSchemaWithSupportFor[newTypeS](t,
		``, `newTypeS`,
		z, `/**
 * newTypeS corresponds to Go type gozod_test.newTypeS (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const newTypeS = z.string().brand("newTypeS");
export type newTypeS = z.infer<typeof newTypeS>;
`,
		examples[newTypeS]{
			simpleExample[newTypeS]("", `""`),
			simpleExample[newTypeS]("foo", `"foo"`),
		},
		rejects{`null`, `undefined`, `0`},
	)
	assertSchemaWithSupportFor[newTypeB](t,
		"", `newTypeB`,
		z, `/**
 * newTypeB corresponds to Go type gozod_test.newTypeB (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const newTypeB = z.boolean().brand("newTypeB");
export type newTypeB = z.infer<typeof newTypeB>;
`,
		examples[newTypeB]{
			simpleExample[newTypeB](true, `true`),
			simpleExample[newTypeB](false, `false`),
		},
		rejects{`null`, `undefined`, `0`},
	)
	assertSchemaWithSupportFor[newTypeI](t,
		"", `newTypeI`,
		z, `/**
 * newTypeI corresponds to Go type gozod_test.newTypeI (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const newTypeI = z.number().int().brand("newTypeI");
export type newTypeI = z.infer<typeof newTypeI>;
`,
		examples[newTypeI]{
			simpleExample[newTypeI](0, `0`),
			simpleExample[newTypeI](1, `1`),
		},
		rejects{`null`, `undefined`, `0.5`},
	)
	assertSchemaWithSupportFor[newTypeF](t,
		"", `newTypeF`,
		z, `/**
 * newTypeF corresponds to Go type gozod_test.newTypeF (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const newTypeF = z.number().brand("newTypeF");
export type newTypeF = z.infer<typeof newTypeF>;
`,
		examples[newTypeF]{
			simpleExample[newTypeF](0, `0`),
			simpleExample[newTypeF](1, `1`),
		},
		rejects{`null`, `undefined`, `""`},
	)

	// named structs are not branded
	assertSchemaWithSupportFor[newTypeStruct](t,
		"", `newTypeStruct`,
		z, `/**
 * newTypeStruct corresponds to Go type gozod_test.newTypeStruct (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const newTypeStruct = z.object({ A: z.string() });
export type newTypeStruct = z.infer<typeof newTypeStruct>;
`,
		examples[newTypeStruct]{
			simpleExample(newTypeStruct{A: ""}, `{"A":""}`),
			simpleExample(newTypeStruct{A: "foo"}, `{"A":"foo"}`),
		},
		rejects{`null`, `undefined`, `{}`},
	)
}

type (
	notSerialized string
	asUnderscore  string
	demoStruct    struct {
		nonExported                string
		Exported                   string
		ExportedHyphenTag          notSerialized `json:"-"`
		ExportedHyphenCommaTag     asUnderscore  `json:"-,"`
		ThisWillBeRenamed          string        `json:"renamed"`
		ThisWillBeRenamed2IMeanToo string        `json:"renamed2,"`
		StrStr                     string        `json:",string"`
		IntStr                     int           `json:",string"`
		FloatStr                   float64       `json:",string"`
		BoolStr                    bool          `json:",string"`
		NullableBoolStr            *bool         `json:",string"`
		SliceNonStr                []float64     `json:",string"`
	}
	omittablesStruct struct {
		Omittable       string `json:",omitempty"`
		Omittable2      string `json:"renamedOmittable,omitempty"`
		OptionalBoolStr bool   `json:",string,omitempty"`
		NullishBoolStr  *bool  `json:",string,omitempty"`
	}
)

func TestGenerateForStructTypes(t *testing.T) {
	assertSchemaWithSupportFor[demoStruct](t,
		``, `demoStruct`,
		z, `/**
 * asUnderscore corresponds to Go type gozod_test.asUnderscore (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const asUnderscore = z.string().brand("asUnderscore");
export type asUnderscore = z.infer<typeof asUnderscore>;

/**
 * demoStruct corresponds to Go type gozod_test.demoStruct (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const demoStruct = z.object({
    Exported: z.string(),
    "-": asUnderscore,
    renamed: z.string(),
    renamed2: z.string(),
    StrStr: z.string().transform(s => JSON.parse(s)).pipe(z.string()),
    IntStr: z.string().transform(s => JSON.parse(s)).pipe(z.number().int()),
    FloatStr: z.string().transform(s => JSON.parse(s)).pipe(z.number()),
    BoolStr: z.string().transform(s => JSON.parse(s)).pipe(z.boolean()),
    NullableBoolStr: z.string().transform(s => JSON.parse(s)).pipe(z.boolean()).nullable(),
    SliceNonStr: z.array(z.number()).nullable().transform(a => a ?? []),
});
export type demoStruct = z.infer<typeof demoStruct>;
`,
		examples[demoStruct]{
			example[demoStruct]{
				demoStruct{
					Exported:                   "A",
					ExportedHyphenCommaTag:     "B",
					ThisWillBeRenamed:          "C",
					ThisWillBeRenamed2IMeanToo: "D",
					StrStr:                     "E",
					IntStr:                     42,
					FloatStr:                   1.23456789,
					BoolStr:                    true,
					NullableBoolStr:            ptr(true),
					SliceNonStr:                []float64{1.5},
				},
				[]demoStruct{
					{
						nonExported:                "NOPE",
						Exported:                   "A",
						ExportedHyphenTag:          "NOPE",
						ExportedHyphenCommaTag:     "B",
						ThisWillBeRenamed:          "C",
						ThisWillBeRenamed2IMeanToo: "D",
						StrStr:                     "E",
						IntStr:                     42,
						FloatStr:                   1.23456789,
						BoolStr:                    true,
						NullableBoolStr:            ptr(true),
						SliceNonStr:                []float64{1.5},
					},
				},
				`{"Exported":"A","-":"B","renamed":"C","renamed2":"D","StrStr":"\"E\"","IntStr":"42","FloatStr":"1.23456789","BoolStr":"true","NullableBoolStr":"true","SliceNonStr":[1.5]}`,
			},
		},
		rejects{`null`, `undefined`, `{}`},
	)
	assertSchemaWithSupportFor[omittablesStruct](t,
		``, `omittablesStruct`,
		z, `/**
 * omittablesStruct corresponds to Go type gozod_test.omittablesStruct (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
 */
export const omittablesStruct = z.object({
    Omittable: z.string().optional(),
    renamedOmittable: z.string().optional(),
    OptionalBoolStr: z.string().transform(s => JSON.parse(s)).pipe(z.boolean()).optional(),
    NullishBoolStr: z.string().transform(s => JSON.parse(s)).pipe(z.boolean()).nullable().optional(),
});
export type omittablesStruct = z.infer<typeof omittablesStruct>;
`,
		examples[omittablesStruct]{
			simpleExample(omittablesStruct{}, `{}`),
			simpleExample(omittablesStruct{Omittable: "A"}, `{"Omittable":"A"}`),
			simpleExample(omittablesStruct{Omittable2: "B"}, `{"renamedOmittable":"B"}`),
			simpleExample(omittablesStruct{OptionalBoolStr: true}, `{"OptionalBoolStr":"true"}`),
			simpleExample(omittablesStruct{NullishBoolStr: ptr(false)}, `{"NullishBoolStr":"false"}`),
		},
		rejects{`null`, `undefined`, `""`},
	)
}

// TODO handle Map types

//type (
//	embeddedableStringA struct {
//		A string
//	}
//	embeddableStringBWithInnerA struct {
//		B string
//		embeddedableStringA
//	}
//	embeddedableBoolA struct {
//		A bool
//	}
//	embeddableBoolBWithInnerA struct {
//		B bool
//		embeddedableStringA
//	}
//)
//
//func TestEmbbeddedTypes(t *testing.T) {
//	// TODO improve behaviour around non-object embedded types that don't define their own behaviour
//	//assertSimpleSchemaFor[struct{ string }](t,
//	//	z, `z.object({})`,
//	//	examples[struct{ string }]{
//	//		example[struct{ string }]{struct{ string }{}, []struct{ string }{struct{ string }{"foo"}}, `{}`},
//	//	},
//	//	rejects{`null`, `undefined`, `"foo"`},
//	//)
//	//assertSimpleSchemaFor[struct{ string }](t,
//	//	z, `z.string()`,
//	//	examples[struct{ string }]{
//	//		simpleExample(struct{ string }{"foo"}, `"foo"`),
//	//	},
//	//	rejects{`null`, `undefined`, `0`},
//	//)
//}

// TODO behaviour around competing property names in embedded structs
//func TestGenerate(t *testing.T) {
//	assertSchemaWithSupportFor[struct {
//		embeddable[string]
//	}](t, "", `embeddable`, z, `/**
// * embeddable corresponds to Go type gozod_test.embeddable (in package "github.com/softwaretechnik-berlin/goats/gotypes/gozod_test").
// */
//z.object({
//    Embedded: z.string(),
//    DoubleEmbedded: z.string(),
//})`)
//
//	assertSimpleSchemaFor[struct {
//		TopLevel string
//		embeddable[string]
//	}](t, z, `z.object({
//    TopLevel: z.string(),
//    Embedded: z.string(),
//    DoubleEmbedded: z.string(),
//})`)
//
//	assertSimpleSchemaFor[struct {
//		embeddable[string]
//		embeddable2[string]
//	}](t, z, `z.object({})`)
//
//	assertSimpleSchemaFor[struct {
//		embeddable[string]
//		extraEmbeddable[int]
//	}](t, z, `z.object({
//    Embedded: z.string(),
//    DoubleEmbedded: z.string(),
//})`)
//
//	// TODO generics (alt 1: instantiate with _ in name; alt 2: generate schema functions)
//	// TODO self-recursive types
//	// TODO generally recursive types that recurse in only 1 type
//	// TODO alternate method for optional properties that excludes the ability to assign undefined
//}
//
//type jsonMarshalledByte byte
//
//var _ encoding.TextMarshaler = textMarshalledByte(0)
//
//func (f jsonMarshalledByte) MarshalJSON() ([]byte, error) {
//	return ([]byte)(fmt.Sprintf(`{"type":"byte",value:%d}`, f)), nil
//}
//
//type textMarshalledByte byte
//
//var _ encoding.TextMarshaler = textMarshalledByte(0)
//
//func (f textMarshalledByte) MarshalText() ([]byte, error) {
//	return ([]byte)(fmt.Sprintf(`{%d}`, f)), nil
//}

func assertTypeScript(t *testing.T, source ts.Source, expectedImports string, expectedCode string) {
	expected := expectedCode
	if expectedImports != "" {
		expected = expectedImports + "\n\n" + expected
	}
	assert.Equal(t, expected, source.String())
}

func assertTypeScriptRepresentationOf(t *testing.T, schema zod.ZodType, expectedImports string, expectedCode string) {
	assertTypeScript(t, schema.TypeScript(), expectedImports, expectedCode)
}

var sharedCommentsLoader = comments.NewLoader()
var testMapper = gozod.NewMapper(gozod.WithCommentsLoader(sharedCommentsLoader))

func negativeExample(expression string) {

}

func assertSimpleSchemaFor[T any](t *testing.T,
	expectedImports string, expectedCode string,
	examples examples[T], rejects rejects,
) {
	m := gozod.NewMapper(gozod.WithCommentsLoader(sharedCommentsLoader))
	assertTypeScriptRepresentationOf(t, m.Resolve(reflective.TypeFor[T]()).Value, expectedImports, expectedCode)
	assertTypeScript(t, gozod.SupportingDeclarations(m), "", "")

	assertExamplesAndRejects(t, examples, rejects)
}

func assertSchemaWithSupportFor[T any](t *testing.T,
	expectedResolutionImports string, expectedResolutionCode string,
	expectedSupportImports string, expectedSupportCode string,
	examples examples[T], rejects rejects,
) {
	m := gozod.NewMapper(gozod.WithCommentsLoader(sharedCommentsLoader))
	assertTypeScriptRepresentationOf(t, m.Resolve(reflective.TypeFor[T]()).Value, expectedResolutionImports, expectedResolutionCode)
	assertTypeScript(t, gozod.SupportingDeclarations(m), expectedSupportImports, expectedSupportCode)

	assertExamplesAndRejects(t, examples, rejects)
}

func assertExamplesAndRejects[T any](t *testing.T, examples examples[T], rejects rejects) {
	for _, e := range examples {
		assertMarshalledJSON := func(value T) {
			marshalled, err := json.Marshal(value)
			if assert.NoError(t, err) {
				assert.Equal(t, e.json, string(marshalled), `%v example: %#v should marshall to JSON %#v`, reflect.TypeFor[T](), value, e.json)
			}
		}
		assertMarshalledJSON(e.value)
		for _, value := range e.identicallyMarshalled {
			assertMarshalledJSON(value)
		}

		var unmarshalled T
		if assert.NoError(t, json.Unmarshal(([]byte)(e.json), &unmarshalled)) {
			assert.True(t, reflect.DeepEqual(e.value, unmarshalled), `%v example: JSON %#v should unmarshall to %#v but got %#v`, reflect.TypeFor[T](), e.json, e.value, unmarshalled)
		}

		// TODO check that TypeScript accepts the value's type
		// TODO check that the schema accepts the value
		// TODO check the schema's output value
		// TODO check that the expected output value is accepted by inferred schema output type.
	}

	for _, r := range rejects {
		var unmarshalled T
		if json.Unmarshal(([]byte)(r), &unmarshalled) == nil {
			assert.Zero(t, unmarshalled)
			marshalledZero, err := json.Marshal(unmarshalled)
			if assert.NoError(t, err) {
				assert.NotEqual(t, r, string(marshalledZero))
			}
		}
		// TODO check that TypeScript rejects the value's type
		// TODO check that the schema rejects the value
	}
}

type examples[T any] []example[T]
type rejects []string

type example[T any] struct {
	value                 T
	identicallyMarshalled []T
	json                  string
}

func simpleExample[T any](value T, json string) example[T] {
	return example[T]{value, nil, json}
}

func ptr[A any](value A) *A {
	return &value
}
