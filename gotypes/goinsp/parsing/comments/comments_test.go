package comments

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
)

func TestLoadTypeComment(t *testing.T) {
	loader := NewLoader()

	assert.Equal(t, ``, loader.Load(reflective.TypeFor[uncommentedType]()))

	assert.Equal(t, ``, loader.Load(reflective.TypeFor[solitaryUncommentedTypeInGroupWithoutComment]()))

	assert.Equal(t, `solitaryUncommentedTypeInGroupWithComment is a test case for this package.
It's a type that's declared in a declaration group of which it is the sole member.
And it doesn't have its own comment, whereas its group does.
`, loader.Load(reflective.TypeFor[solitaryUncommentedTypeInGroupWithComment]()))

	assert.Equal(t, `solitaryCommentedTypeInGroupWithoutComment is a test case for this package.
It's a type that's declared in a declaration group of which it is the sole member.
And it doesn't has its own comment. It's group has none.
`, loader.Load(reflective.TypeFor[solitaryCommentedTypeInGroupWithoutComment]()))
}
