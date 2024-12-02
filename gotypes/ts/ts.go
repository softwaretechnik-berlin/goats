// Package ts provides utilities for generating TypeScript code in Go.
package ts

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"

	"github.com/softwaretechnik-berlin/goats/gotypes/util"
)

// Source represents TypeScript source code.
//
// Source tracks imports separately for the rest of the source text so that imports of combined Source values can be
// deduplicated and hoisted to the top of the source. It is indentation-aware, so that nested Sources produce
// human-friendly output.
type Source interface {
	// String renders the source as TypeScript, complete with import statements.
	String() string

	addToImports(imps *imports)
	writeSourceTo(w sourceWriter)
}

var _ Source = Identifier("")
var _ Source = sourceText("")
var _ Source = sourceWithImport{}
var _ Source = sourceGroup{}

// Array outputs the given elements surrounded by `[` and `]` and interspersed with `,`.
// It gives reasonable line-breaking, whitespace and indentation.
func Array(elements ...Source) Source {
	return sourceGroup{&array, elements}
}

func AsSource(str string) Source {
	return sourceText(str)
}

// DocComment renders the given string as a multiline `/** … */`-style doc-comment.
// It is indentation-aware.
func DocComment(comment string) Source {
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return sourceText("")
	}
	return Statements(
		sourceText("/**"),
		Statements(util.Map(strings.Split(comment, "\n"), func(line string) Source {
			if len(line) == 0 {
				return sourceText(" *")
			}
			return sourceText(" * " + line)
		})...),
		sourceText(" */"),
	)
}

// Identifier is a TypeScript identifier string.
//
// It can be used as Source, and since is of string kind it can also be useful as an identifier in code working with TypeScript.
type Identifier string

func (i Identifier) String() string               { return sourceText(i).String() }
func (i Identifier) addToImports(imps *imports)   { sourceText(i).addToImports(imps) }
func (i Identifier) writeSourceTo(w sourceWriter) { sourceText(i).writeSourceTo(w) }

// ImportedName returns a Source representing a name that has been imported from a module.
func ImportedName(module string, name Identifier) Source {
	return sourceWithImport{tsImport{module, name}, string(name)}
}

// InvokeFunction follows the function by a parenthesized comma-separated list of arguments.
// It gives reasonable line-breaking, whitespace and indentation.
func InvokeFunction(function Source, arguments ...Source) Source {
	return Sourcef("%s%s", function, sourceGroup{&invocation, arguments})
}

// InvokeMethod follows the receiver by a `.`, the method name and parenthesized comma-separated list of arguments.
// It gives reasonable line-breaking, whitespace and indentation.
func InvokeMethod(receiver Source, name Identifier, arguments ...Source) Source {
	return InvokeFunction(Sourcef("%s.%s", receiver, name), arguments...)
}

// NumberLiteral returns literal Source for the given value.
func NumberLiteral[N constraints.Integer | constraints.Float](value N) Source {
	// TODO logic
	return sourceText(fmt.Sprintf(`%v`, value))
}

// Object outputs the given properties as `name: value`-pairs surrounded by `{` and `}` and interspersed with `,`.
// It gives reasonable line-breaking, whitespace and indentation.
func Object(properties ...Property) Source {
	return sourceGroup{&object, util.Map(properties, Property.AsSource)}
}

// Property is a named value for use with Object.
type Property struct {
	Name  string
	Value Source
}

// AsSource represents the property as a `name: value` pair.
func (p Property) AsSource() Source {
	var name Source
	if isValidIdentifier(p.Name) {
		name = sourceText(p.Name)
	} else {
		name = StringLiteral(p.Name)
	}
	return Sourcef(`%s: %s`, name, p.Value)
}

// RegexLiteral represents the given regexp as a TypeScript regex literal.
// TODO there is missing escaping.
func RegexLiteral(re *regexp.Regexp) Source {
	// TODO logic
	return sourceText(fmt.Sprintf(`/%s/`, re))
}

// Sourcef turns format into indentation-aware source,
// replacing `%s` placeholders in it with the given arguments in an indentation-aware way.
func Sourcef(format string, a ...Source) Source {
	if n := strings.Count(format, "%s"); n != len(a) {
		panic(format)
	}
	if len(a) == 0 {
		return sourceText(format)
	}
	return sourceGroup{sourcef{format}, a}
}

// Statements chains the given statements together with indentation-aware newlines.
func Statements(statements ...Source) Source {
	return StatementGroups(0, statements...)
}

// StatementGroups chains groups of statements with the given number of blank lines in between them.
func StatementGroups(blankLinesBetweenGroups int, groups ...Source) Source {
	return sourceGroup{statementsStyle(blankLinesBetweenGroups), groups}
}

var stringEscaper = strings.NewReplacer(`"`, `\"`)

// StringEscape escapes the given string for inclusion in a TypeScript string literal.
// TODO there is missing escaping.
func StringEscape(str string) Source {
	return sourceText(stringEscaper.Replace(str))
}

// StringLiteral represents the given string as a TypeScript string literal.
// TODO there is missing escaping.
func StringLiteral(str string) Source {
	return sourceText(fmt.Sprintf(`"%s"`, StringEscape(str)))
}
