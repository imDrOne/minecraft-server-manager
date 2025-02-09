package dto

type PagePaginationMetaDto struct {
	Page  uint64 `json:"page"`
	Size  uint64 `json:"size"`
	Total uint64 `json:"total"`
	Pages uint64 `json:"pages"`
}

type PageResponseWrapDto struct {
	Data interface{}           `json:"data"`
	Meta PagePaginationMetaDto `json:"meta"`
}
