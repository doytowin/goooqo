package field

import "strings"

func ProcessOr(or any) (string, []any) {
	conditions, args := buildConditions(or)
	return "(" + strings.Join(conditions, " OR ") + ")", args
}
