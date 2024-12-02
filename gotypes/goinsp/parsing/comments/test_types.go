package comments

type uncommentedType struct{}

type (
	solitaryUncommentedTypeInGroupWithoutComment struct{}
)

// typeWithMultilineComment is a test case for this package.
// You'll find that it's used to check on comment parsing.
//
// Hopefully it works perfectly
type typeWithMultilineComment struct{}

// solitaryUncommentedTypeInGroupWithComment is a test case for this package.
// It's a type that's declared in a declaration group of which it is the sole member.
// And it doesn't have its own comment, whereas its group does.
type (
	solitaryUncommentedTypeInGroupWithComment struct{}
)

type (
	// solitaryCommentedTypeInGroupWithoutComment is a test case for this package.
	// It's a type that's declared in a declaration group of which it is the sole member.
	// And it doesn't has its own comment. It's group has none.
	solitaryCommentedTypeInGroupWithoutComment struct{}
)

// Oh, look at this group comment!
// This isn't the best comment for the type declared inside it.
type (
	// solitaryCommentedTypeInGroupWithComment is a test case for this package.
	// It's a type that's declared in a declaration group of which it is the sole member.
	// And it doesn't has its own comment. It's group has a comment that isn't specific to this type.
	solitaryCommentedTypeInGroupWithComment struct{}
)
