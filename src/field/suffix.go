package field

import "regexp"

func CreateOpMap() map[string]string {
	opMap := make(map[string]string)
	opMap["Gt"] = ">"
	opMap["Ge"] = ">="
	opMap["Lt"] = "<"
	opMap["Le"] = "<="
	return opMap
}

var opMap = CreateOpMap()
var regx = regexp.MustCompile(`(\w+)(Gt|Ge|Lt|Le)$`)

func Process(fieldName string) string {
	if match := regx.FindStringSubmatch(fieldName); len(match) > 0 {
		return match[1] + " " + opMap[match[2]] + " ?"
	}
	return fieldName + " = ?"
}
