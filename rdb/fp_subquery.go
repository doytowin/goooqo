package rdb

import (
	"reflect"
	"regexp"
	"strings"
)

type subquery struct {
	select_ string
	from    string
}

func processSubquery(field reflect.StructField, value reflect.Value) (string, []any) {
	subq := buildSubquery(field)
	fieldName := strings.TrimRightFunc(field.Name, func(r rune) bool {
		return 0x30 < r && r <= 0x39 // remove trailing digits, such as ScoreGt1
	})
	columnOp := processField(fieldName)
	where, args := BuildWhereClause(value.Elem().Interface())
	condition := columnOp + "(SELECT " + subq.select_ + " FROM " + subq.from + where + ")"
	return condition, args
}

var sqRegx = regexp.MustCompile(`(select|from):([\w()]+)`)

func buildSubquery(field reflect.StructField) (subq *subquery) {
	subqueryStr := field.Tag.Get("subquery")
	submatch := sqRegx.FindAllStringSubmatch(subqueryStr, -1)
	subq = &subquery{}
	for _, group := range submatch {
		if group[1] == "select" {
			subq.select_ = group[2]
		} else if group[1] == "from" {
			if from := emMap[group[2]]; from != nil {
				subq.from = from.TableName
			} else {
				subq.from = group[2]
			}
		}
	}
	return
}

func processField(fieldName string) string {
	column, op := suffixMatch(fieldName)
	return column + op.sign
}
