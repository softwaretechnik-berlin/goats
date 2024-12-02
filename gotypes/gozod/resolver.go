package gozod

import (
	"maps"
)

type Resolver[A, B any] interface {
	Resolve(A) B
}

var _ Resolver[string, int] = (*accountingResolver[string, int, float32])(nil)

type accountingResolver[A comparable, B any, Name comparable] struct {
	delegate Resolver[A, withAccounting[B, Name]]
	Observed accountingInfo[Name]
}

func newAccountingResolver[A comparable, B any, Name comparable](delegate Resolver[A, withAccounting[B, Name]]) accountingResolver[A, B, Name] {
	return accountingResolver[A, B, Name]{delegate, accountingInfo[Name]{}}
}

func (r *accountingResolver[A, B, Name]) Resolve(a A) B {
	resolved := r.delegate.Resolve(a)
	r.Observed.depth = max(r.Observed.depth, resolved.info.depth)
	if len(resolved.info.dependencies) > 0 && len(r.Observed.dependencies) == 0 {
		r.Observed.dependencies = make(map[Name]struct{})
	}
	maps.Copy(r.Observed.dependencies, resolved.info.dependencies)
	return resolved.Value
}

type accountingInfo[Name comparable] struct {
	dependencies map[Name]struct{}
	depth        uint
}

type withAccounting[B any, Name comparable] struct {
	Value B
	info  accountingInfo[Name]
}
