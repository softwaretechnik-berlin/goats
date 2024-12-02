package zod

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
)

type zodAnyType struct {
	source ts.Source
}

var _ ZodType = zodAnyType{}

func (t zodAnyType) Brand(brand string) ZodBranded {
	return chainBrand(t, brand)
}

func (t zodAnyType) Nullable() ZodNullable {
	return zodNullable{t.chain("nullable"), t}
}

func (t zodAnyType) Optional() ZodOptional {
	return zodOptional{t.chain("optional"), t}
}

func (t zodAnyType) Parse(str ts.Source) ts.Source {
	return ts.InvokeMethod(t.source, "parse", str)
}

func (t zodAnyType) Parsef(format string, a ...ts.Source) ts.Source {
	return t.Parse(ts.Sourcef(format, a...))
}

func (t zodAnyType) Pipe(target ZodType) ZodType {
	return t.chain("pipe", target.TypeScript())
}

func (t zodAnyType) Transform(transform ts.Source) ZodType {
	return t.chain("transform", transform)
}

func (t zodAnyType) Transformf(format string, a ...ts.Source) ZodType {
	return t.Transform(ts.Sourcef(format, a...))
}

// TODO reconsider
func (t zodAnyType) DeclaredAs(name ts.Identifier) ZodType {
	return zodAnyType{name}
}

func (t zodAnyType) TypeScript() ts.Source {
	return t.source
}

func (t zodAnyType) chain(name ts.Identifier, args ...ts.Source) zodAnyType {
	return zodAnyType{ts.InvokeMethod(t.source, name, args...)}
}
