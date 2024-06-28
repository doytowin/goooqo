package main

import (
	"github.com/doytowin/goooqo/gen"
	"testing"
)

func Test_GenerateFile(t *testing.T) {
	tests := []struct {
		input, output string
		generator     gen.Generator
	}{
		{input: "./inventory.go", output: "./inventory_query_builder.go", generator: gen.NewMongoGenerator()},
		{input: "./user.go", output: "./user_query_builder.go", generator: gen.NewSqlGenerator()},
	}
	for _, tt := range tests {
		t.Run("Generate for "+tt.input, func(t *testing.T) {
			err := gen.GenerateQueryBuilder(tt.generator, tt.input, tt.output)
			if err != nil {
				t.Fatalf("%s", err)
			}
		})
	}
}
