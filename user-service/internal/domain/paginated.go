package domain

type PaginationParam struct {
	Limit  uint64
	Offset uint64
}

type PaginationResult[T any] struct {
	TotalCount uint64
	Data       []T
}
