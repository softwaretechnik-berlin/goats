package zod

type zodOptional struct {
	zodAnyType
	wrapped ZodType
}

func (n zodOptional) Unwrap() ZodType {
	return n.wrapped
}
