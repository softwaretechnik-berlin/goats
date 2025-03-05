package examples

import (
	"testing"

	"github.com/softwaretechnik-berlin/goats/gotypes/gozod"
)

func TestExample1(t *testing.T) {
	mapper := gozod.NewMapper()

	mapper.ResolveAll(
		gozod.RefFor[Example1](),
	)

	gozod.Generate(mapper, "example_1.ts")
}

func TestExample2(t *testing.T) {
	mapper := gozod.NewMapper()

	mapper.ResolveAll(
		gozod.RefFor[Example2](),
	)

	gozod.Generate(mapper, "example_2.ts")
}

func TestExample3(t *testing.T) {
	mapper := gozod.NewMapper()

	mapper.ResolveAll(
		gozod.RefFor[Example3](),
	)

	gozod.Generate(mapper, "example_3.ts")
}
