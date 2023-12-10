package field

import (
	. "github.com/doytowin/goquery/util"
	"reflect"
	"regexp"
	"strings"
)

type operator struct {
	name, sign string
	process    func(value reflect.Value) (string, []any)
}

func ReadValueToArray(value reflect.Value) (string, []any) {
	return "?", []any{ReadValue(value)}
}

func ReadValueForIn(value reflect.Value) (string, []any) {
	var args []any
	arg := reflect.Indirect(value)
	ph := "("
	for i := 0; i < arg.Len(); i++ {
		args = append(args, arg.Index(i).Int())
		ph += "?"
		if i < arg.Len()-1 {
			ph += ", "
		}
	}
	ph += ")"
	return ph, args
}
func EmptyValue(reflect.Value) (string, []any) {
	return "", []any{}
}

func CreateOpMap() map[string]operator {
	const LIKE = " LIKE "
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", " > ", ReadValueToArray}
	opMap["Ge"] = operator{"Ge", " >= ", ReadValueToArray}
	opMap["Lt"] = operator{"Lt", " < ", ReadValueToArray}
	opMap["Le"] = operator{"Le", " <= ", ReadValueToArray}
	opMap["Not"] = operator{"Not", " != ", ReadValueToArray}
	opMap["Ne"] = operator{"Ne", " <> ", ReadValueToArray}
	opMap["Eq"] = operator{"Eq", " == ", ReadValueToArray}
	opMap["Null"] = operator{"Null", " IS NULL", EmptyValue}
	opMap["NotNull"] = operator{"NotNull", " IS NOT NULL", EmptyValue}
	opMap["In"] = operator{"In", " IN ", ReadValueForIn}
	opMap["NotIn"] = operator{"NotIn", " NOT IN ", ReadValueForIn}
	opMap["Like"] = operator{"Like", LIKE, ReadValueToArray}
	opMap["Contain"] = operator{"Contain", LIKE, func(value reflect.Value) (string, []any) {
		return "?", []any{"%" + ReadValue(value).(string) + "%"}
	}}
	opMap["NotContain"] = operator{"NotContain", " NOT LIKE ", func(value reflect.Value) (string, []any) {
		return "?", []any{"%" + ReadValue(value).(string) + "%"}
	}}
	opMap["Start"] = operator{"Start", LIKE, func(value reflect.Value) (string, []any) {
		return "?", []any{ReadValue(value).(string) + "%"}
	}}
	opMap["NotStart"] = operator{"NotStart", " NOT LIKE ", func(value reflect.Value) (string, []any) {
		return "?", []any{ReadValue(value).(string) + "%"}
	}}
	opMap["End"] = operator{"End", LIKE, func(value reflect.Value) (string, []any) {
		return "?", []any{"%" + ReadValue(value).(string)}
	}}
	return opMap
}

var opMap = CreateOpMap()
var regx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In|Like|Contain|NotContain|Start|NotStart|End)$`)

func Process(fieldName string, value reflect.Value) (string, []any) {
	if match := regx.FindStringSubmatch(fieldName); len(match) > 0 {
		operator := opMap[match[1]]
		column := strings.TrimSuffix(fieldName, match[1])
		column = UnCapitalize(column)
		placeholder, args := operator.process(value)
		return column + operator.sign + placeholder, args
	}
	_, args := ReadValueToArray(value)
	return UnCapitalize(fieldName) + " = ?", args
}
