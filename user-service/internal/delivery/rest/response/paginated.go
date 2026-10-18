package response

type Paginated[T any] struct {
	TotalCount uint64 `json:"total_count"`
	Data       []T    `json:"data"`
}
