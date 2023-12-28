package rdb

import (
	"reflect"
	"regexp"
	"strings"
)

type fpSubquery struct {
	field         reflect.StructField
	column, op    string
	select_, from string
}

func (fp *fpSubquery) Process(value reflect.Value) (string, []any) {
	where, args := BuildWhereClause(value.Elem().Interface())
	return fp.buildCondition(where), args
}

func (fp *fpSubquery) buildCondition(where string) string {
	if em := emMap[fp.from]; em != nil {
		fp.from = em.TableName
	}
	return fp.column + fp.op + "(SELECT " + fp.select_ + " FROM " + fp.from + where + ")"
}

var sqRegx = regexp.MustCompile(`(select|from):([\w()]+)`)

func buildFpSubquery(field reflect.StructField) (fp *fpSubquery) {
	subqueryStr := field.Tag.Get("subquery")
	submatch := sqRegx.FindAllStringSubmatch(subqueryStr, -1)
	fp = &fpSubquery{field: field}
	for _, group := range submatch {
		if group[1] == "select" {
			fp.select_ = group[2]
		} else if group[1] == "from" {
			fp.from = group[2]
		}
	}
	fieldName := strings.TrimRightFunc(fp.field.Name, func(r rune) bool {
		return 0x30 < r && r <= 0x39 // remove trailing digits, such as 1 in ScoreGt1
	})
	column, op := suffixMatch(fieldName)
	fp.column, fp.op = column, op.sign
	return
}
