package rdb

import (
	"reflect"
	"strings"
)

type fpBasicArrayByOr struct {
	fpSuffix FieldProcessor
}

func buildFpBasicArrayByOr(fieldName string) FieldProcessor {
	return &fpBasicArrayByOr{fpSuffix: buildFpSuffix(strings.TrimSuffix(fieldName, "Or"))}
}

func (fp *fpBasicArrayByOr) Process(value reflect.Value) (string, []any) {
	value = value.Elem()
	var args, arr []any
	conditions := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		conditions[i], arr = fp.fpSuffix.Process(value.Index(i))
		args = append(args, arr...)
	}
	return fpForOr.connect(conditions), args
}
