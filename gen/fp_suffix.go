package gen

import "regexp"

var (
	opMap     = CreateOpMap()
	suffixRgx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In|Like|NotLike|Contain|NotContain|Start|NotStart|End|NotEnd)$`)
)

type operator struct {
	name string
	sign map[string]string
}

func CreateOpMap() map[string]operator {
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", map[string]string{"mongo": "$gt"}}
	opMap["Lt"] = operator{"Lt", map[string]string{"mongo": "$lt"}}
	return opMap
}
