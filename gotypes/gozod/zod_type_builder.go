package gozod

import (
	"cmp"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/util"
	"github.com/softwaretechnik-berlin/goats/gotypes/zod"
)

type zodTypeBuilder struct {
	config
}

func newZodTypeBuilder(config config) zodTypeBuilder {
	return zodTypeBuilder{config}
}

type goToZodMapper = mapper[goinsp.Type, zod.ZodType, ts.Identifier, zod.SchemaAndTypeDeclaration]

func NewMapper(options ...Option) goToZodMapper {
	return newMapper[goinsp.Type, zod.ZodType, ts.Identifier, zod.SchemaAndTypeDeclaration](newZodTypeBuilder(newConfig(options...)))
}

func NewMapperWithSupport(options ...Option) goToZodMapper {
	return newMapper[goinsp.Type, zod.ZodType, ts.Identifier, zod.SchemaAndTypeDeclaration](newZodTypeBuilder(newConfig(options...)))
}

var _ builder[goinsp.Type, zod.ZodType, zod.SchemaAndTypeDeclaration] = zodTypeBuilder{}

func (b zodTypeBuilder) Build(t goinsp.Type, resolver Resolver[goinsp.Type, zod.ZodType]) (schema zod.ZodType, declaration zod.SchemaAndTypeDeclaration, hasDeclaration bool) {
	schema = b.buildRawSchema(t, resolver)
	if transform, ok := lookupConfig(b.transforms, t); ok {
		schema = schema.Transform(transform(resolver))
	}
	schemaBeforeTemplating := schema
	if template, ok := lookupConfig(b.templates, t); ok {
		schema = applyTemplateTransform(schema, template)
	}
	name, ok := b.name(t)
	if !ok {
		return
	}
	if b.shouldBrand(t, schemaBeforeTemplating) {
		schema = schema.Brand(string(name))
	}
	docComment := fmt.Sprintf("%s corresponds to Go type %s (in package %#v).\n", name, t, t.PkgPath())
	if goComment := b.commentsLoader.Load(t); goComment != "" {
		docComment += "The comment on the original Go type follows.\n\n" + goComment
	}
	return schema.DeclaredAs(name), zod.NewSchemaAndTypeDeclaration(docComment, name, schema), true
}

func (b zodTypeBuilder) name(t goinsp.Type) (ts.Identifier, bool) {
	if name, ok := lookupConfig(b.names, t); ok {
		return name, true
	}
	if _, unnamed := lookupConfig(b.unnamedTypes, t); unnamed || t.PkgPath() == "" {
		return "", false
	}
	return ts.Identifier(t.Name().String()), true
}

func (b zodTypeBuilder) shouldBrand(t goinsp.Type, schema zod.ZodType) bool {
	if _, ok := lookupConfig(b.schemas, t); ok {
		return false
	}
	if _, ok := lookupConfig(b.undiscriminatedUnions, t); ok {
		return false
	}
	if _, ok := lookupConfig(b.discriminatedUnions, t); ok {
		return false
	}
	if _, ok := lookupConfig(b.transforms, t); ok {
		return false
	}

	//_, isZodObject := schema.(zod.ZodObject)
	for {
		switch typed := schema.(type) {
		case zod.ZodObject:
			return false
		case zod.ZodBranded:
			schema = typed.Unwrap()
		default:
			return true
		}
	}

	//return !isZodObject
}

func (b zodTypeBuilder) buildRawSchema(t goinsp.Type, resolver Resolver[goinsp.Type, zod.ZodType]) zod.ZodType {
	if schema, ok := lookupConfig(b.schemas, t); ok {
		return schema(resolver)
	}
	if types, ok := lookupConfig(b.undiscriminatedUnions, t); ok {
		return zod.Union(mapSlice(types, resolver.Resolve)...)
	}
	if union, ok := lookupConfig(b.discriminatedUnions, t); ok {
		return zod.DiscriminatedUnion(union.DiscriminatorProperty, mapSlice(union.Types, resolver.Resolve)...)
	}
	if _, ok := lookupConfig(b.templates, t); !ok && t.Implements(reflective.TypeFor[encoding.TextMarshaler]()) {
		var schema zod.ZodType = zod.String()
		//for _, s := range strings.Split(b.commentsLoader.LoadMethod(t, "MarshalText"), "\n") {
		//	s = strings.TrimSpace(s)
		//	if s, ok := strings.CutPrefix(s, "zod:"); ok {
		//		s = strings.TrimSpace(s)
		//		if s, ok := strings.CutPrefix(s, "transform:"); ok {
		//			s = strings.TrimSpace(s)
		//			schema = schema.Transform(ts.AsSource(s))
		//		} else {
		//			panic(s)
		//		}
		//	}
		//}
		return schema
	}
	switch t.Kind() {
	case reflect.Bool:
		return zod.Boolean()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return zod.Number().Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return zod.Number().NonNegative().Int()
	case reflect.Float32, reflect.Float64:
		return zod.Number()
	case reflect.Array:
		return zod.Array(resolver.Resolve(t.Elem())).Length(t.Len())
	case reflect.Interface:
		return zod.Any()
	case reflect.Map:
		schema := zod.Record(resolver.Resolve(t.Key()), resolver.Resolve(t.Elem()))
		// Nil maps are marshalled to JSON null
		// TODO make it possible to configure things such that we assert that we don't emit nil values
		if true {
			schema = schema.Nullable()
			// TODO make it possible to opt out of the homogenizing transformation.
			if true {
				schema = schema.Transformf(`r => r ?? {}`)
			}
		}
		return schema
	case reflect.Pointer:
		return zod.EnsureNullable(resolver.Resolve(t.Elem()))
	case reflect.Slice:
		// Go encodes non-nil byte slices as strings using base64.
		isBase64Encoded := t.Elem().Kind() == reflect.Uint8 && !t.Elem().Implements(reflective.TypeFor[json.Marshaler]()) && !t.Elem().Implements(reflective.TypeFor[encoding.TextMarshaler]())
		var schema zod.ZodType
		if isBase64Encoded {
			schema = zod.String()
		} else {
			schema = zod.Array(resolver.Resolve(t.Elem()))
		}
		// Nil slices are marshalled to JSON null
		// TODO make it possible to configure things such that we assert that we don't emit nil values
		if true {
			schema = schema.Nullable()
			// TODO make it possible to opt out of the homogenizing transformation.
			if true {
				if isBase64Encoded {
					schema = schema.Transformf(`a => a ?? ""`)
				} else {
					schema = schema.Transformf(`a => a ?? []`)
				}
			}
		}
		return schema
	case reflect.String:
		return zod.String()
	case reflect.Struct:

		for i := range t.NumField() {
			field := t.Field(i)
			tsgenTag := field.Tag.Get("gotypes")
			if tagHasFlag(tsgenTag, "value") {
				return b.resolveFieldSchema(field.Type(), field.Tag.Get("json"), tsgenTag, resolver)
			}
		}

		var schema util.Optional[zod.ZodObject]
		var properties []zod.ShapeProperty
		if discriminator, ok := lookupConfig(b.discriminators, t); ok {
			properties = append(properties, zod.ShapeProperty{discriminator.Property, zod.Literal(discriminator.Value)})
		}

		hasFields := false
		embeddedJSONTypes := 0
		var embeddedJSONType goinsp.Type
		forEachTopLevelJSONFieldAndEmbeddedType(t,
			func(name string, field goinsp.StructField, tag string) { hasFields = true },
			func(t goinsp.Type) { embeddedJSONTypes++; embeddedJSONType = t },
		)
		if !hasFields && embeddedJSONTypes == 1 {
			return resolver.Resolve(embeddedJSONType)
		}
		addPropertiesToSchema := func() {
			if schema.HasValue {
				schema.V = schema.V.Extend(properties...)
			} else {
				schema = util.AsOptional(zod.Object(properties...))
			}
		}
		forEachTopLevelJSONFieldAndEmbeddedType(t,
			func(name string, field goinsp.StructField, tag string) {
				// TODO embedded fields with name in json tag or embedded interfaces as object fields
				// TODO embedded object fields inline, subject to complicated visibility rules
				properties = append(properties, zod.ShapeProperty{name, b.resolveFieldSchema(field.Type(), tag, field.Tag.Get("gotypes"), resolver)})
			},
			func(t goinsp.Type) {
				if len(properties) > 0 {
					addPropertiesToSchema()
					properties = nil
				}
				// TODO what if it's not a ZodObject?
				resolved := resolver.Resolve(t)
				embeddedObjectSchema, ok := resolved.(zod.ZodObject)
				if !ok {
					panic(fmt.Sprintf("%#v", resolved))
				}
				schema = util.AsOptional(util.MapOptionalWithDefault(schema, embeddedObjectSchema, func(schema zod.ZodObject) zod.ZodObject { return schema.Merge(embeddedObjectSchema) }))
			},
		)
		if len(properties) > 0 || schema.IsNone() {
			addPropertiesToSchema()
		}
		return schema.MustGet()
	default:
		panic(t.Kind())
	}
}

func (b zodTypeBuilder) resolveFieldSchema(t goinsp.Type, jsonTag string, tsgenTag string, resolver Resolver[goinsp.Type, zod.ZodType]) zod.ZodType {
	schema := resolver.Resolve(t)
	//fromJsonString := false
	//var schema zod.ZodType
	//kindSupportsJSONStringFlag(t, jsonTag, fromJsonString)
	//if schema = nil {
	//}
	if tagHasFlag(jsonTag, "string") && kindSupportsJSONStringFlag(t) {
		needsNullable := false
		schema, needsNullable = zod.StripNullable(schema)
		schema = zod.String().Transformf("s => JSON.parse(s)").Pipe(schema)
		if needsNullable {
			schema = zod.EnsureNullable(schema)
		}
	}
	if tagHasFlag(tsgenTag, "nullable") {
		//if needsNullable {
		schema = zod.EnsureNullable(schema)
	}

	if tagHasFlag(jsonTag, "omitempty") {
		schema = schema.Optional()
	}
	return schema
}

func kindSupportsJSONStringFlag(t goinsp.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map:
		return false
	case reflect.Pointer:
		return kindSupportsJSONStringFlag(t.Elem())
	case reflect.Slice:
		return false
	case reflect.String:
		return true
	case reflect.Struct, reflect.UnsafePointer:
		return false
	default:
		panic(t)
	}
}

func SupportingDeclarations(mapper goToZodMapper) ts.Source {
	type declaration = mappedValue[goinsp.Type, zod.ZodType, ts.Identifier, zod.SchemaAndTypeDeclaration]

	declarationsByGoPackage := make(map[goinsp.ImportPath][]declaration)
	simpleDeclarationsByGoPackage := make(map[goinsp.ImportPath][]declaration)
	for _, decl := range mapper.declarations {
		pkg := decl.in.PkgPath()
		declarationsByGoPackage[pkg] = append(declarationsByGoPackage[pkg], decl)
		if decl.in.WithoutTypeArguments() == decl.in {
			simpleDeclarationsByGoPackage[pkg] = append(simpleDeclarationsByGoPackage[pkg], decl)
		}
	}

	packagesToOutput := maps.Keys(declarationsByGoPackage)
	slices.SortFunc(packagesToOutput, func(a, b goinsp.ImportPath) int {
		group := func(pkg goinsp.ImportPath) uint {
			if strings.ContainsRune(string(pkg), '.') {
				return 1
			} // privilege standard library types
			return 0
		}
		if r := cmp.Compare(group(a), group(b)); r != 0 {
			return r
		}
		return cmp.Compare(a, b)
	})

	var packagesInOutputOrder []goinsp.ImportPath
	ready := func(declarations map[goinsp.ImportPath][]declaration) int {
		return slices.IndexFunc(packagesToOutput, func(path goinsp.ImportPath) bool {
			for _, decl := range declarations[path] {
				for depName := range decl.declaration.info.dependencies {
					depPath := mapper.declarations[depName].in.PkgPath()
					if depPath != path && slices.Contains(packagesToOutput, depPath) {
						return false
					}
				}
			}
			return true
		})
	}
	for len(packagesToOutput) > 0 {
		index := ready(declarationsByGoPackage)
		if index == -1 {
			index = ready(simpleDeclarationsByGoPackage)
			if index == -1 {
				index = 0
			}
		}
		packagesInOutputOrder = append(packagesInOutputOrder, packagesToOutput[index])
		packagesToOutput = slices.Delete(packagesToOutput, index, index+1)
	}

	declarations := make([]declaration, 0, len(mapper.declarations))
	for _, p := range packagesInOutputOrder {
		packageDeclarations := declarationsByGoPackage[p]
		slices.SortFunc(packageDeclarations, func(a, b declaration) int {
			if r := cmp.Compare(a.declaration.info.depth, b.declaration.info.depth); r != 0 {
				return r
			}
			return cmp.Compare(a.declaration.Value.Identifier(), b.declaration.Value.Identifier())
		})
		declarations = append(declarations, packageDeclarations...)
	}

	return ts.StatementGroups(1, util.Map(declarations, func(d declaration) ts.Source {
		return d.declaration.Value.TypeScript()
	})...)
}
