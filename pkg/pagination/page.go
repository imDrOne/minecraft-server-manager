package pagination

type PageMetadata struct {
	Page  uint64
	Size  uint64
	Total uint64
	Pages uint64
}

func NewPageMetadata(page, size, total uint64) *PageMetadata {
	pages := total / size
	if total%size != 0 {
		pages++ // округляем вверх, если есть остаток
	}
	return &PageMetadata{Page: page, Size: size, Total: total, Pages: pages}
}

type PageRequest struct {
	Page uint64
	Size uint64
}

func (p PageRequest) Offset() uint64 {
	return (p.Page - 1) * p.Size
}
