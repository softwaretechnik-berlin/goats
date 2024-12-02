package gozod

// A builder knows how to build a B for a given A given a resolver for other As it depends on.
// It may also return a Declaration as a byproduct.
type builder[A, B, Declaration any] interface {
	Build(A, Resolver[A, B]) (B, Declaration, bool)
}
