package examples

import (
	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
	"github.com/softwaretechnik-berlin/goats/gotypes/gozod"
	"testing"
)

func TestExample1(t *testing.T) {
	mapper := gozod.NewMapper()

	mapper.ResolveAll(
		reflective.TypeFor[Example1](),
	)

	gozod.Generate(mapper, "example_1.ts")
}

func TestExample2(t *testing.T) {
	mapper := gozod.NewMapper()

	mapper.ResolveAll(
		reflective.TypeFor[Example2](),
	)

	gozod.Generate(mapper, "example_2.ts")
}
