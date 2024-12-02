package ts

import (
	"regexp"
)

// Ideally we would build a regexp based on https://262.ecma-international.org/14.0/index.html#prod-IdentifierName,
// e.g. ignoring unicode escapes we might do something like:
//
//	regexp.MustCompile(`^[\p{Other_ID_Start}$_][\p{ID_Continue}$\u200C\u200D]*$`)
//
// But Go 1.22.5's unicode package doesn't support ID_Start and ID_Continue, and even if it did, its regexp package
// [doesn't support Unicode character properties](https://github.com/golang/go/issues/10851#event-435488430).
//
// So we'll go small and simple for now.
var identifier = regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)

func isValidIdentifier(s string) bool {
	return identifier.MatchString(s)
}
