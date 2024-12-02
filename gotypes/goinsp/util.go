package goinsp

type NoneWhenZero[T any] struct{ Value T }

type PotentiallyUnavailable[T any] struct {
	Available bool
	Value     T
}
