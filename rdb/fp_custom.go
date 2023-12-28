package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"strings"
)

type fpCustom struct {
	field     *reflect.StructField
	condition string
	phCnt     int
}

func buildFpCustom(field reflect.StructField) FieldProcessor {
	condition := field.Tag.Get("condition")
	phCnt := strings.Count(condition, "?")
	return &fpCustom{&field, condition, phCnt}
}

func (fp *fpCustom) Process(value reflect.Value) (string, []any) {
	var arr []any
	arg := ReadValue(value)
	for j := 0; j < fp.phCnt; j++ {
		arr = append(arr, arg)
	}
	return fp.condition, arr
}
