package field

import "regexp"

type operator struct {
	name, sign, placeholder string
}

func CreateOpMap() map[string]operator {
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", ">", " ?"}
	opMap["Ge"] = operator{"Ge", ">=", " ?"}
	opMap["Lt"] = operator{"Lt", "<", " ?"}
	opMap["Le"] = operator{"Le", "<=", " ?"}
	opMap["Not"] = operator{"Not", "!=", " ?"}
	opMap["Ne"] = operator{"Ne", "<>", " ?"}
	opMap["Eq"] = operator{"Eq", "==", " ?"}
	opMap["Null"] = operator{"Null", "IS NULL", ""}
	return opMap
}

var opMap = CreateOpMap()
var regx = regexp.MustCompile(`(\w+)(Gt|Ge|Lt|Le|Not|Ne|Eq|Null)$`)

func Process(fieldName string) string {
	if match := regx.FindStringSubmatch(fieldName); len(match) > 0 {
		operator := opMap[match[2]]
		return match[1] + " " + operator.sign + operator.placeholder
	}
	return fieldName + " = ?"
}
