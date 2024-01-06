package web

import (
	"github.com/doytowin/goooqo/core"
	"reflect"
	"strconv"
	"strings"
)

var converterMap = map[reflect.Type]func(v []string) (any, error){}

func RegisterConverter(typeName reflect.Type, converter func(v []string) (any, error)) {
	converterMap[typeName] = converter
}

func init() {
	RegisterConverter(reflect.TypeOf(true), func(v []string) (any, error) {
		return strings.EqualFold(v[0], "TRue"), nil
	})

	RegisterConverter(reflect.TypeOf(0), func(v []string) (any, error) {
		return strconv.Atoi(v[0])
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf(0)), func(v []string) (any, error) {
		v0, err := strconv.Atoi(v[0])
		return &v0, err
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf([]int{0})), func(params []string) (any, error) {
		var v []int
		num, err := strconv.Atoi(params[0])
		strArr := strings.Split(params[0], ",")
		for _, s := range strArr {
			num, err = strconv.Atoi(s)
			if core.NoError(err) {
				v = append(v, num)
			}
		}
		return &v, err
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf("")), func(v []string) (any, error) {
		return &v[0], nil
	})
}
