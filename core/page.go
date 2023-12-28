package core

type PageQuery struct {
	PageNumber *int    `json:"page,omitempty"`
	PageSize   *int    `json:"size,omitempty"`
	Sort       *string `json:"sort,omitempty"`
}

func (pageQuery PageQuery) GetPageNumber() int {
	page := 0
	if pageQuery.PageNumber != nil && *pageQuery.PageNumber > 1 {
		page = *pageQuery.PageNumber - 1
	}
	return page
}

func (pageQuery PageQuery) GetPageSize() int {
	size := 10
	if pageQuery.PageSize != nil && *pageQuery.PageSize > 0 {
		size = *pageQuery.PageSize
	}
	return size
}

func (pageQuery PageQuery) CalcOffset() int {
	return pageQuery.GetPageNumber() * pageQuery.GetPageSize()
}

func (pageQuery PageQuery) GetSort() *string {
	return pageQuery.Sort
}

func (pageQuery PageQuery) NeedPaging() bool {
	return pageQuery.PageSize != nil || pageQuery.PageNumber != nil
}
