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
	opMap["Eq"] = operator{"Eq", map[string]string{"mongo": "$eq"}}
	opMap["Gt"] = operator{"Gt", map[string]string{"mongo": "$gt"}}
	opMap["Ge"] = operator{"Ge", map[string]string{"mongo": "$gte"}}
	opMap["Lt"] = operator{"Lt", map[string]string{"mongo": "$lt"}}
	opMap["Le"] = operator{"Le", map[string]string{"mongo": "$lte"}}
	return opMap
}
