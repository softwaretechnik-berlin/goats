package gozod

import (
	"os"

	"github.com/samber/lo"
)

func Generate(mapper goToZodMapper, outputFileName string) {
	declarations := SupportingDeclarations(mapper)

	w := lo.Must(os.Create(outputFileName))
	defer func() { lo.Must0(w.Close()) }()

	lo.Must(w.WriteString(declarations.String()))
}

func GenerateString(mapper goToZodMapper, outputFileName string) string {
	declarations := SupportingDeclarations(mapper)
	return declarations.String()
}
