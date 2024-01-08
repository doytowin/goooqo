package web

import (
	"reflect"
	"testing"
)

func TestConverter(t *testing.T) {
	type args struct {
		typeName reflect.Type
		params   []string
	}
	tests := []struct {
		name   string
		expect any
		val    func(any) any
		args   args
	}{
		{
			"Support int", 32, func(a any) any { return a },
			args{typeName: reflect.TypeOf(32), params: []string{"32"}},
		},
		{
			"Support *int", 32, func(a any) any { return *a.(*int) },
			args{typeName: reflect.PointerTo(reflect.TypeOf(1)), params: []string{"32"}},
		},
		{
			"Support *float64", 22.5, func(a any) any { return *a.(*float64) },
			args{typeName: reflect.PointerTo(reflect.TypeOf(0.1)), params: []string{"22.5"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := converterMap[tt.args.typeName]
			if f == nil {
				t.Fatal("Not support: ", tt.args.typeName)
			}
			actual, err := f(tt.args.params)
			if err != nil {
				t.Error(err)
			} else if !(tt.expect == tt.val(actual)) {
				t.Errorf("Expected: %v, but got %v", tt.expect, tt.val(actual))
			}
		})
	}
}
