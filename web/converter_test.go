/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package web

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestConverter(t *testing.T) {
	t.Run("Type Converter", func(t *testing.T) {
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
			{
				"Support *bool", true, func(a any) any { return *a.(*bool) },
				args{typeName: reflect.PointerTo(reflect.TypeOf(true)), params: []string{"true"}},
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
	})

	t.Run("Resolve Query for", func(t *testing.T) {
		type Unit struct {
			Name *string `json:"name,omitempty"`
		}
		type SizeQuery struct {
			HLt  *int  `json:"hLt,omitempty"`
			HGe  *int  `json:"hGe,omitempty"`
			Unit *Unit `json:"unit,omitempty"`
		}
		type qo struct {
			Size *SizeQuery `json:"size,omitempty"`
		}
		type args struct {
			queryMap url.Values
			query    qo
		}
		tests := []struct {
			name   string
			expect string
			args   args
		}{
			{"Nested Parameters", `{"size":{"hLt":20}}`,
				args{url.Values{"Size.HLt": {"20"}}, qo{}}},
			{"Two Nested Parameters", `{"size":{"hLt":20,"hGe":10}}`,
				args{url.Values{"Size.HLt": {"20"}, "Size.HGe": {"10"}}, qo{}}},
			{"Level Three of Nested Parameters", `{"size":{"unit":{"name":"cm"}}}`,
				args{url.Values{"Size.Unit.Name": {"cm"}}, qo{}}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ResolveQuery(tt.args.queryMap, &tt.args.query)
				data, _ := json.Marshal(tt.args.query)
				actual := string(data)
				if tt.expect != actual {
					t.Errorf("\nExpected: %s\nBut got : %s", tt.expect, actual)
				}
			})
		}
	})
}
