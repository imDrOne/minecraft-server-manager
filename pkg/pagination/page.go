package pagination

type PageMetadata struct {
	page  uint64
	size  uint64
	total uint64
	pages uint64
}

func NewPageMetadata(page, size, total uint64) (PageMetadata, error) {
	if err := validatePage(page); err != nil {
		return PageMetadata{}, err
	}
	if err := validateSize(size); err != nil {
		return PageMetadata{}, err
	}

	var pages uint64
	calcPages(&pages, size, total)
	return PageMetadata{page, size, total, pages}, nil
}

func (p PageMetadata) Page() uint64 {
	return p.page
}

func (p PageMetadata) Size() uint64 {
	return p.size
}

func (p PageMetadata) Total() uint64 {
	return p.total
}

func (p PageMetadata) Pages() uint64 {
	return p.pages
}

func (p PageMetadata) ToDTO() PagePaginationMetaDto {
	return PagePaginationMetaDto{
		Page:  p.page,
		Size:  p.size,
		Total: p.total,
		Pages: p.pages,
	}
}

func calcPages(val *uint64, size, total uint64) {
	*val = total / size
	if total%size != 0 {
		*val++
	}
}

type PageRequest struct {
	page uint64
	size uint64
}

func (p PageRequest) Page() uint64 {
	return p.page
}

func (p PageRequest) Size() uint64 {
	return p.size
}

func (p PageRequest) ToPageMeta(total uint64) PageMetadata {
	var pages uint64
	calcPages(&pages, p.size, total)
	return PageMetadata{
		page:  p.page,
		size:  p.size,
		total: total,
		pages: pages,
	}
}

func NewPageRequest(page uint64, size uint64) (PageRequest, error) {
	if err := validatePage(page); err != nil {
		return PageRequest{}, err
	}
	if err := validateSize(size); err != nil {
		return PageRequest{}, err
	}

	return PageRequest{page: page, size: size}, nil
}

func (p PageRequest) Offset() uint64 {
	return (p.page - 1) * p.size
}

type PaginatedResult[T any] struct {
	Data []T
	Meta PageMetadata
}
