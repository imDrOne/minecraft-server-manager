package pagination

// PagePaginationRequestDto todo: add validation
type PagePaginationRequestDto struct {
	Page uint64 `json:"page"`
	Size uint64 `json:"size"`
}

func (d PagePaginationRequestDto) ToValue() PageRequest {
	return PageRequest{
		page: d.Page,
		size: d.Size,
	}
}

type PagePaginationResponseWrapDto struct {
	Data interface{}           `json:"data"`
	Meta PagePaginationMetaDto `json:"meta"`
}

type PagePaginationMetaDto struct {
	Page  uint64 `json:"page"`
	Size  uint64 `json:"size"`
	Total uint64 `json:"total"`
	Pages uint64 `json:"pages"`
}
