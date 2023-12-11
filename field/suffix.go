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

var escapePtn = regexp.MustCompile("[\\\\_%]")

func ReadLikeValue(value reflect.Value) string {
	s := ReadValue(value).(string)
	return escapePtn.ReplaceAllString(s, "\\$0")
}

func CreateOpMap() map[string]operator {
	const Like = " LIKE "
	const NotLike = " NOT LIKE "
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
	opMap["Like"] = operator{"Like", Like, func(value reflect.Value) (string, []any) {
		s := ReadValue(value).(string)
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}}
	opMap["NotLike"] = operator{"NotLike", NotLike, func(value reflect.Value) (string, []any) {
		s := ReadValue(value).(string)
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}}
	opMap["Contain"] = operator{"Contain", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}}
	opMap["NotContain"] = operator{"NotContain", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}}
	opMap["Start"] = operator{"Start", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}}
	opMap["NotStart"] = operator{"NotStart", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}}
	opMap["End"] = operator{"End", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}}
	opMap["NotEnd"] = operator{"End", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}}
	return opMap
}

func resolvePlaceHolder(arg string) string {
	ph := "?"
	if strings.Contains(arg, "\\") {
		ph = ph + " ESCAPE '\\'"
	}
	return ph
}

var opMap = CreateOpMap()
var regx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In|Like|NotLike|Contain|NotContain|Start|NotStart|End|NotEnd)$`)

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
