package core

import (
	"fmt"
	"regexp"
	"strings"
)

type PageQuery struct {
	PageNumber *int    `json:"page,omitempty"`
	PageSize   *int    `json:"size,omitempty"`
	Sort       *string `json:"sort,omitempty"`
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

var sortRgx = regexp.MustCompile("(?i)(\\w+)(,(asC|dEsc))?;?")

func (pageQuery PageQuery) BuildSortClause() string {
	if pageQuery.Sort == nil {
		return ""
	}
	submatch := sortRgx.FindAllStringSubmatch(*pageQuery.Sort, -1)
	var sort = make([]string, len(submatch))
	for i, group := range submatch {
		sort[i] = group[1]
		if group[3] != "" {
			sort[i] += " " + strings.ToUpper(group[3])
		}
	}
	return " ORDER BY " + strings.Join(sort, ", ")
}
