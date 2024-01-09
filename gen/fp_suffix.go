package gen

import "regexp"

var (
	opMap     = make(map[string]map[string]operator)
	suffixRgx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|NotNull|Null|NotIn|In|Like|NotLike|Contain|NotContain|Start|NotStart|End|NotEnd)$`)
)

type operator struct {
	name   string
	sign   string
	format string
}
