package rdb

import (
	"fmt"
	"regexp"
	"strings"
)

func BuildPageClause(sql *string, offset int, size int) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", *sql, size, offset)
}

var sortRgx = regexp.MustCompile("(?i)(\\w+)(,(asC|dEsc))?;?")

func BuildSortClause(sort *string) string {
	if sort == nil {
		return ""
	}
	submatch := sortRgx.FindAllStringSubmatch(*sort, -1)
	var orderBy = make([]string, len(submatch))
	for i, group := range submatch {
		orderBy[i] = group[1]
		if group[3] != "" {
			orderBy[i] += " " + strings.ToUpper(group[3])
		}
	}
	return " ORDER BY " + strings.Join(orderBy, ", ")
}
