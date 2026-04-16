package dto

type Pagination struct {
	Limit  uint64 `form:"limit" json:"limit"`
	Offset uint64 `form:"offset" json:"offset"`
}

type ListResponse[T any] struct {
	Items      []T    `json:"items"`
	TotalCount uint64 `json:"total_count"`
	Limit      uint64 `json:"limit"`
	Offset     uint64 `json:"offset"`
}

func (p Pagination) Normalize(defaultLimit, maxLimit uint64) Pagination {
	if p.Limit == 0 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	return p
}
