package gen

import "regexp"

var (
	opMap     = CreateOpMap()
	suffixRgx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In|Like|NotLike|Contain|NotContain|Start|NotStart|End|NotEnd)$`)
)

type operator struct {
	name   string
	sign   map[string]string
	format string
}

func CreateOpMap() map[string]operator {
	opMap := make(map[string]operator)
	opMap["Eq"] = operator{name: "Eq", sign: map[string]string{"mongo": "$eq"}}
	opMap["Ne"] = operator{name: "Ne", sign: map[string]string{"mongo": "$ne"}}
	opMap["Not"] = operator{name: "Not", sign: map[string]string{"mongo": "$ne"}}
	opMap["Gt"] = operator{name: "Gt", sign: map[string]string{"mongo": "$gt"}}
	opMap["Ge"] = operator{name: "Ge", sign: map[string]string{"mongo": "$gte"}}
	opMap["Lt"] = operator{name: "Lt", sign: map[string]string{"mongo": "$lt"}}
	opMap["Le"] = operator{name: "Le", sign: map[string]string{"mongo": "$lte"}}
	opMap["In"] = operator{name: "In", sign: map[string]string{"mongo": "$in"}}
	opMap["NotIn"] = operator{name: "NotIn", sign: map[string]string{"mongo": "$nin"}}
	opMap["Null"] = operator{
		name:   "Null",
		sign:   map[string]string{"mongo": "$type"},
		format: "\td = append(d, D{{\"%s\", D{{\"%s\", 10}}}})",
	}
	opMap["NotNull"] = operator{
		name:   "NotNull",
		sign:   map[string]string{"mongo": "$type"},
		format: "\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", 10}}}}}})",
	}
	opMap["Contain"] = operator{
		name:   "Contain",
		sign:   map[string]string{"mongo": "$regex"},
		format: "\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	}
	opMap["NotContain"] = operator{
		name:   "NotContain",
		sign:   map[string]string{"mongo": "$regex"},
		format: "\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", q.%s}}}}}})",
	}
	return opMap
}
