package field

import "strings"

func ProcessOr(or interface{}) (string, []any) {
	conditions, args := buildConditions(or)
	return strings.Join(conditions, " OR "), args
}
