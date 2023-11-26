package field

import (
	"regexp"
	"strings"
)

type operator struct {
	name, sign, placeholder string
}

func CreateOpMap() map[string]operator {
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", " > ", "?"}
	opMap["Ge"] = operator{"Ge", " >= ", "?"}
	opMap["Lt"] = operator{"Lt", " < ", "?"}
	opMap["Le"] = operator{"Le", " <= ", "?"}
	opMap["Not"] = operator{"Not", " != ", "?"}
	opMap["Ne"] = operator{"Ne", " <> ", "?"}
	opMap["Eq"] = operator{"Eq", " == ", "?"}
	opMap["Null"] = operator{"Null", " IS NULL", ""}
	opMap["NotNull"] = operator{"NotNull", " IS NOT NULL", ""}
	opMap["In"] = operator{"In", " IN ", ""}
	opMap["NotIn"] = operator{"NotIn", " NOT IN ", ""}
	return opMap
}

var opMap = CreateOpMap()
var regx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In)$`)

func Process(fieldName string) string {
	if match := regx.FindStringSubmatch(fieldName); len(match) > 0 {
		operator := opMap[match[1]]
		column, _ := strings.CutSuffix(fieldName, match[1])
		return column + operator.sign + operator.placeholder
	}
	return fieldName + " = ?"
}
