package rdb

import (
	"fmt"
	"github.com/doytowin/goooqo/core"
	"reflect"
	"testing"
)

type mapping struct {
	field, expect string
	expectValue   any
	value         reflect.Value
}

func TestProcess(t *testing.T) {
	useCases := []mapping{
		{"id", "id = ?", []int64{5}, reflect.ValueOf(5)},
		{"idGt", "id > ?", []int64{5}, reflect.ValueOf(5)},
		{"idGe", "id >= ?", []int64{5}, reflect.ValueOf(5)},
		{"idLt", "id < ?", []int64{5}, reflect.ValueOf(5)},
		{"idLe", "id <= ?", []int64{5}, reflect.ValueOf(5)},
		{"idNot", "id != ?", []int64{5}, reflect.ValueOf(5)},
		{"idNe", "id <> ?", []int64{5}, reflect.ValueOf(5)},
		{"idEq", "id = ?", []int64{5}, reflect.ValueOf(5)},
		{"idNull", "id IS NULL", nil, reflect.ValueOf(true)},
		{"idNotNull", "id IS NOT NULL", nil, reflect.ValueOf(true)},
		{"idIn", "id IN (?, ?, ?)", []int64{5, 6, 7}, reflect.ValueOf([]int{5, 6, 7})},
		{"memoIn", "memo IN (?, ?)", []string{"Good", "Bad"}, reflect.ValueOf([]string{"Good", "Bad"})},
		{"idNotIn", "id NOT IN (?, ?, ?)", []int{5, 6, 7}, reflect.ValueOf([]int{5, 6, 7})},
		{"MemoContain", "memo LIKE ?", "[%at%]", reflect.ValueOf("at")},
		{"MemoContain", "memo LIKE ? ESCAPE '\\'", "[%a\\_\\%t%]", reflect.ValueOf("a_%t")},
		{"MemoNotContain", "memo NOT LIKE ?", "[%at%]", reflect.ValueOf("at")},
		{"MemoStart", "memo LIKE ?", "[at%]", reflect.ValueOf("at")},
		{"MemoNotStart", "memo NOT LIKE ?", "[at%]", reflect.ValueOf("at")},
		{"MemoEnd", "memo LIKE ?", "[%at]", reflect.ValueOf("at")},
		{"MemoNotEnd", "memo NOT LIKE ?", "[%at]", reflect.ValueOf("at")},
		{"MemoLike", "memo LIKE ? ESCAPE '\\'", "[%\\_at%]", reflect.ValueOf("%\\_at%")},
		{"MemoNotLike", "memo NOT LIKE ?", "[%at%]", reflect.ValueOf("%at%")},
		{"memoNull", "memo IS NULL", nil, reflect.ValueOf(core.PBool(true))},
		{"memoNull", "memo IS NOT NULL", nil, reflect.ValueOf(core.PBool(false))},
	}
	for _, useCase := range useCases {
		t.Run(useCase.field, func(t *testing.T) {
			actual, arg := Process(useCase.field, useCase.value)
			if actual != useCase.expect {
				t.Errorf("Expected: %s, but got %s", useCase.expect, actual)
			}
			if !((len(arg) == 0 && useCase.expectValue == nil) || fmt.Sprint(arg) == fmt.Sprint(useCase.expectValue)) {
				t.Errorf("Expected: %s, but got %s", useCase.expectValue, arg)
			}
		})
	}

}
