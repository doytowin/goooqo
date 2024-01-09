package rdb

import (
	"reflect"
	"regexp"
	"strings"
)

type fpSubquery struct {
	column, op    string
	select_, from string
}

func (fp *fpSubquery) Process(value reflect.Value) (string, []any) {
	where, args := BuildWhereClause(value.Elem().Interface())
	return fp.buildCondition(where), args
}

func (fp *fpSubquery) Subquery() string {
	if em := emMap[fp.from]; em != nil {
		fp.from = em.TableName
	}
	return fp.column + fp.op + "(SELECT " + fp.select_ + " FROM " + fp.from
}

func (fp *fpSubquery) buildCondition(where string) string {
	return fp.Subquery() + where + ")"
}

var sqRegx = regexp.MustCompile(`(select|from):([\w()]+)`)

func buildFpSubquery(field reflect.StructField) *fpSubquery {
	subqueryStr := field.Tag.Get("subquery")
	return BuildSubquery(subqueryStr, field.Name)
}

func BuildSubquery(subqueryStr string, fieldName string) (fp *fpSubquery) {
	fp = &fpSubquery{}
	submatch := sqRegx.FindAllStringSubmatch(subqueryStr, -1)
	for _, group := range submatch {
		if group[1] == "select" {
			fp.select_ = group[2]
		} else if group[1] == "from" {
			fp.from = group[2]
		}
	}
	fieldName = strings.TrimRightFunc(fieldName, func(r rune) bool {
		return 0x30 < r && r <= 0x39 // remove trailing digits, such as 1 in ScoreGt1
	})
	column, op := suffixMatch(fieldName)
	fp.column, fp.op = column, op.sign
	return
}
