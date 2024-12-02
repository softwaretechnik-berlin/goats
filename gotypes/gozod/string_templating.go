package gozod

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/util"
	"github.com/softwaretechnik-berlin/goats/gotypes/zod"
)

type templateEmbedding interface {
	RegexString() string
	Parse(str ts.Source) ts.Source
}

var _ templateEmbedding = numberEmbedding{}
var _ templateEmbedding = stringEmbedding{}

type numberEmbedding struct{ schema zod.ZodNumber }

func (n numberEmbedding) RegexString() string {
	regex := `\d+`
	if !n.schema.IsInt() {
		regex += `(?:\.\d+)?`
	}
	if !n.schema.IsNonNegative() {
		regex = `-?` + regex
	}
	return regex
}

func (n numberEmbedding) Parse(str ts.Source) ts.Source {
	return n.schema.Parsef("Number(%s)", str)
}

type stringEmbedding struct{ schema zod.ZodString }

func (s stringEmbedding) RegexString() string {
	return ".*"
}

func (s stringEmbedding) Parse(str ts.Source) ts.Source {
	return str
}

func resolveEmbedding(schema zod.ZodType) templateEmbedding {
	switch schema := schema.(type) {
	case zod.ZodBranded:
		return resolveEmbedding(schema.Unwrap())
	case zod.ZodNumber:
		return numberEmbedding{schema}
	case zod.ZodString:
		return stringEmbedding{schema}
	default:
		panic(fmt.Sprintf("%#v", schema))
	}
}

func applyTemplateTransform(schema zod.ZodType, template string) zod.ZodType {
	r, transformMatch := fromTemplatedString(schema, template)
	return zod.String().Transformf(`(s, ctx) => {
    const re = %s;
    const match = re.exec(s);
    if (!match) {
        ctx.addIssue({ code: %s.ZodIssueCode.custom, message: "expected string of the form %s matching " + re });
        return z.NEVER;
    }
    return %s;
}`, ts.RegexLiteral(regexp.MustCompile(`^`+r+`$`)), ts.ImportedName("zod", "z"), ts.StringEscape(ts.StringLiteral(template).String()), transformMatch) // TODO
}

func fromTemplatedString(schema zod.ZodType, template string) (string, ts.Source) {
	if schema, ok := schema.(zod.ZodObject); ok {
		return objectFromTemplatedString(schema, template)
	}
	embedding := resolveEmbedding(schema)
	prefix, suffix, ok := strings.Cut(template, "{}")
	if !ok {
		panic(fmt.Sprintf("can't find placeholder {} in template %#v", template))
	}
	regex := regexp.QuoteMeta(prefix) + "(" + embedding.RegexString() + ")" + regexp.QuoteMeta(suffix)
	return regex, embedding.Parse(ts.Sourcef("match[1]"))
}

func objectFromTemplatedString(schema zod.ZodObject, template string) (string, ts.Source) {
	shape := schema.Shape()
	embeddings := util.Map(shape, func(p zod.ShapeProperty) templateEmbedding { return resolveEmbedding(p.Schema) })
	placeholder := regexp.MustCompile(`\{(` + strings.Join(util.Map(shape, func(p zod.ShapeProperty) string { return regexp.QuoteMeta(p.Name) }), `|`) + `)\}`)
	var regex strings.Builder
	outputProperties := make([]ts.Property, len(shape))
	for matchIndex := 1; ; matchIndex++ {
		loc := placeholder.FindStringIndex(template)
		if loc == nil {
			break
		}
		regex.WriteString(regexp.QuoteMeta(template[:loc[0]]))
		propertyIndex := slices.IndexFunc(shape, func(p zod.ShapeProperty) bool { return p.Name == template[loc[0]+1:loc[1]-1] })
		embedding := embeddings[propertyIndex]
		regex.WriteByte('(')
		regex.WriteString(embedding.RegexString())
		regex.WriteByte(')')
		outputProperties[propertyIndex] = ts.Property{shape[propertyIndex].Name, embedding.Parse(ts.Sourcef("match[%s]", ts.NumberLiteral(matchIndex)))}
		template = template[loc[1]:]
	}
	regex.WriteString(regexp.QuoteMeta(template))
	return regex.String(), ts.Object(outputProperties...)
}
