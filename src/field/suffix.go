package field

import "strings"

func Process(fieldName string) string {
	if strings.HasSuffix(fieldName, "Gt") {
		columnName := strings.TrimSuffix(fieldName, "Gt")
		return columnName + " > ?"
	}
	if strings.HasSuffix(fieldName, "Ge") {
		columnName := strings.TrimSuffix(fieldName, "Ge")
		return columnName + " >= ?"
	}
	if strings.HasSuffix(fieldName, "Lt") {
		columnName := strings.TrimSuffix(fieldName, "Lt")
		return columnName + " < ?"
	}
	if strings.HasSuffix(fieldName, "Le") {
		columnName := strings.TrimSuffix(fieldName, "Le")
		return columnName + " <= ?"
	}
	return fieldName + " = ?"
}
