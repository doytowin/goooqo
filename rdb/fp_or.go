package rdb

import (
	"reflect"
	"strings"
)

// build multiple conditions and connect them by the Connector
type fpMulti struct {
	Connector string
}

func buildFpOr() FieldProcessor {
	return &fpMulti{Connector: " OR "}
}

func (f *fpMulti) Process(value reflect.Value) (string, []any) {
	conditions, args := buildConditions(value.Elem().Interface())
	condition := strings.Join(conditions, f.Connector)
	if f.Connector == " OR " {
		condition = "(" + condition + ")"
	}
	return condition, args
}
