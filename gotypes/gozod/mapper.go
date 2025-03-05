package gozod

import (
	"fmt"
)

// mapper is a tool for accumulating Declarations created from a B Value with a ID,
// where we know how to build B values and potentially a ID for each A Value.
type mapper[A comparable, B any, ID comparable, Declaration withIdentifier[ID]] struct {
	namesByInput map[A]ID
	declarations map[ID]mappedValue[A, B, ID, Declaration]
	builder      builder[A, B, Declaration]
}

type mappedValue[A comparable, B any, Name comparable, Declaration withIdentifier[Name]] struct {
	in          A
	declaration withAccounting[Declaration, Name]
	reference   withAccounting[B, Name]
}

type withIdentifier[I any] interface {
	Identifier() I
}

func newMapper[A comparable, B any, Identifier comparable, Declaration withIdentifier[Identifier]](builder builder[A, B, Declaration]) mapper[A, B, Identifier, Declaration] {
	return mapper[A, B, Identifier, Declaration]{
		make(map[A]Identifier),
		make(map[Identifier]mappedValue[A, B, Identifier, Declaration]),
		builder,
	}
}

func (m mapper[A, B, ID, Declaration]) Resolve(a A) withAccounting[B, ID] {
	if name, ok := m.namesByInput[a]; ok {
		decl := m.declarations[name]
		return decl.reference
	}

	r := newAccountingResolver(m)
	b, declaration, hasDeclaration := m.builder.Build(a, &r)
	if !hasDeclaration {
		return withAccounting[B, ID]{b, r.Observed}
	}

	name := declaration.Identifier()
	if _, ok := m.declarations[name]; ok {
		panic(fmt.Sprintf("would declare %v as %v, but there is already another declaration with that name", a, name))
	}
	decl := mappedValue[A, B, ID, Declaration]{
		a,
		withAccounting[Declaration, ID]{declaration, r.Observed},
		withAccounting[B, ID]{b, accountingInfo[ID]{map[ID]struct{}{name: {}}, r.Observed.depth + 1}},
	}
	m.declarations[name] = decl
	m.namesByInput[a] = name
	return decl.reference
}

func (m mapper[A, B, ID, Declaration]) ResolveAll(inputs ...A) {
	for _, a := range inputs {
		m.Resolve(a)
	}
}
