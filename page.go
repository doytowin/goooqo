package goquery

import "fmt"

type PageQuery struct {
	PageNumber *int
	PageSize   *int
}

func (pageQuery *PageQuery) buildPageClause() string {
	size := *pageQuery.PageSize
	offset := (*pageQuery.PageNumber - 1) * *pageQuery.PageSize

	return fmt.Sprintf(" LIMIT %d OFFSET %d", size, offset)
}

func (pageQuery *PageQuery) needPaging() bool {
	return pageQuery.PageSize != nil || pageQuery.PageNumber != nil
}
