package core

import "fmt"

type PageQuery struct {
	PageNumber *int
	PageSize   *int
}

func (pageQuery PageQuery) BuildPageClause() string {
	size := 10
	page := 0

	if pageQuery.PageNumber != nil && *pageQuery.PageNumber > 1 {
		page = *pageQuery.PageNumber - 1
	}
	if pageQuery.PageSize != nil && *pageQuery.PageSize > 0 {
		size = *pageQuery.PageSize
	}
	offset := page * size

	return fmt.Sprintf(" LIMIT %d OFFSET %d", size, offset)
}

func (pageQuery PageQuery) NeedPaging() bool {
	return pageQuery.PageSize != nil || pageQuery.PageNumber != nil
}
