package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"testing"
)

func TestBuildWhereClause(t *testing.T) {

	tests := []struct {
		name       string
		query      any
		expect     string
		expectArgs []any
	}{
		{
			name:       "Support custom condition",
			query:      TestQuery{Account: PStr("f0rb"), Deleted: PBool(true)},
			expect:     " WHERE (username = ? OR email = ?) AND deleted = ?",
			expectArgs: []any{"f0rb", "f0rb", true},
		},
		{
			name:       "Given field with type *bool and suffix Null, when assigned true, then map to IS NULL",
			query:      TestQuery{EmailNull: PBool(true)},
			expect:     " WHERE email IS NULL",
			expectArgs: []any{},
		},
		{
			name:       "Given field with type *bool and suffix Null, when assigned false, then map to IS NOT NULL",
			query:      TestQuery{EmailNull: PBool(false)},
			expect:     " WHERE email IS NOT NULL",
			expectArgs: []any{},
		},
		{
			name:       "Given field with type *bool and suffix Null, when not assigned, then map nothing",
			query:      TestQuery{},
			expect:     "",
			expectArgs: []any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, args := BuildWhereClause(tt.query)
			if actual != tt.expect {
				t.Errorf("\nExpected: %s\nBut got : %s", tt.expect, actual)
			}
			if !reflect.DeepEqual(args, tt.expectArgs) {
				t.Errorf("BuildWhereClause() args = %v, expect %v", args, tt.expectArgs)
			}
		})
	}

}
