package dto

type PagePaginationDto struct {
	Page uint64 `json:"page"`
	Size uint64 `json:"size"`
}
