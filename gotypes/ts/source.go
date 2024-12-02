package ts

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type tsImport struct {
	module string
	name   Identifier
}

type imports struct {
	byName map[tsImport]struct{}
}

func (i *imports) Add(imp tsImport) {
	if i.byName == nil {
		i.byName = make(map[tsImport]struct{})
	}
	i.byName[imp] = struct{}{}
}

type sourceText string

func (s sourceText) String() string               { return string(s) }
func (s sourceText) addToImports(_ *imports)      {}
func (s sourceText) writeSourceTo(w sourceWriter) { w.WriteString(string(s)) }

type sourceWithImport struct {
	tsImport tsImport
	text     string
}

func (s sourceWithImport) String() string               { return s.text }
func (s sourceWithImport) addToImports(imps *imports)   { imps.Add(s.tsImport) }
func (s sourceWithImport) writeSourceTo(w sourceWriter) { w.WriteString(s.text) }

type sourceGroup struct {
	style    groupStyle
	elements []Source
}

func (s sourceGroup) String() string {
	return toString(s)
}

func (s sourceGroup) addToImports(imps *imports) {
	for _, e := range s.elements {
		e.addToImports(imps)
	}
}

func (s sourceGroup) writeSourceTo(w sourceWriter) {
	s.style.writeGroupTo(w, s.elements)
}

type groupStyle interface {
	writeGroupTo(w sourceWriter, elements []Source)
}

var _ groupStyle = (*bracedStyle)(nil)
var _ groupStyle = statementsStyle(0)
var _ groupStyle = sourcef{}

type bracedStyle struct {
	open, singleLinePadding, closed sourceText
	multilineThreshold              int
}

var (
	invocation = bracedStyle{"(", "", ")", 5}
	array      = bracedStyle{"[", "", "]", 2}
	object     = bracedStyle{"{", " ", "}", 2}
)

func (s bracedStyle) writeGroupTo(w sourceWriter, elements []Source) {
	s.open.writeSourceTo(w)
	if len(elements) > 0 {
		if len(elements) < s.multilineThreshold {
			s.singleLinePadding.writeSourceTo(w)
			for i, element := range elements {
				if i != 0 {
					w.WriteString(", ")
				}
				element.writeSourceTo(w)
			}
			s.singleLinePadding.writeSourceTo(w)
		} else {
			inner := w.indent()
			for _, e := range elements {
				inner.WriteNewline()
				e.writeSourceTo(inner)
				inner.WriteString(",")
			}
			w.WriteNewline()
		}
	}
	s.closed.writeSourceTo(w)
}

type statementsStyle int

func (s statementsStyle) writeGroupTo(w sourceWriter, elements []Source) {
	for i, e := range elements {
		if i != 0 {
			for range s {
				w.WriteNewline()
			}
		}
		e.writeSourceTo(w)
		w.EnsureNewline()
	}
}

type sourcef struct {
	format string
}

func (s sourcef) writeGroupTo(w sourceWriter, elements []Source) {
	lines := strings.Split(s.format, "\n")
	remaining := lines[0]
	for _, arg := range elements {
		for {
			prefix, suffix, ok := strings.Cut(remaining, "%s")
			if ok {
				w.WriteString(prefix)
				arg.writeSourceTo(w.indentBy(lines[0][:len(lines[0])-len(strings.TrimLeft(lines[0], " \t"))]))
				remaining = suffix
				break
			}
			w.WriteString(remaining)
			w.WriteNewline()
			lines = lines[1:]
			remaining = lines[0]
		}
	}
	for {
		w.WriteString(remaining)
		if len(lines) == 1 {
			return
		}
		w.WriteNewline()
		lines = lines[1:]
		remaining = lines[0]
	}
}

type indentationAwareWriter struct {
	delegate         io.StringWriter
	justWroteNewline bool
}

func (w *indentationAwareWriter) WriteStringAtIndentation(indentation string, s string) {
	if w.justWroteNewline {
		w.writeString(indentation)
		w.justWroteNewline = false
	}
	w.writeString(s)
}

func (w *indentationAwareWriter) writeString(s string) {
	if _, err := w.delegate.WriteString(s); err != nil {
		panic(err)
	}
	if len(s) > 0 {
		w.justWroteNewline = s[len(s)-1] == '\n'
	}
}

func (w *indentationAwareWriter) WriteNewline() {
	w.writeString("\n")
}

func (w *indentationAwareWriter) EnsureNewline() {
	if !w.justWroteNewline {
		w.WriteNewline()
	}
}

type sourceWriter struct {
	*indentationAwareWriter
	indentation string
}

func (w sourceWriter) WriteString(s string) {
	w.indentationAwareWriter.WriteStringAtIndentation(w.indentation, s)
}

func (w sourceWriter) indent() sourceWriter {
	return w.indentBy("    ")
}

func (w sourceWriter) indentBy(additional string) sourceWriter {
	return sourceWriter{w.indentationAwareWriter, w.indentation + additional}
}

func toString(s Source) string {
	var buf bytes.Buffer
	var w io.StringWriter = &buf
	var imps imports
	s.addToImports(&imps)
	sortedImports := maps.Keys(imps.byName)
	slices.SortFunc(sortedImports, func(a, b tsImport) int {
		if r := cmp.Compare(a.module, b.module); r != 0 {
			return r
		}
		return cmp.Compare(a.name, b.name)
	})
	sw := sourceWriter{&indentationAwareWriter{w, false}, ""}
	if len(sortedImports) > 0 {
		// TODO group by module
		for _, imp := range sortedImports {
			sw.WriteString(fmt.Sprintf("import { %s } from %s;\n", imp.name, StringLiteral(imp.module)))
		}
		sw.WriteString("\n")
	}
	s.writeSourceTo(sw)
	return buf.String()
}
