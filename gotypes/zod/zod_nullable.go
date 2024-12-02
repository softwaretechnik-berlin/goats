package zod

type zodNullable struct {
	zodAnyType
	wrapped ZodType
}

func (n zodNullable) Unwrap() ZodType {
	return n.wrapped
}
