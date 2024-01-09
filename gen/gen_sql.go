package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

type SqlGenerator struct {
	*generator
}

func init() {
	sqlOpMap := make(map[string]operator)
	sqlOpMap["Eq"] = operator{name: "Eq", sign: "="}
	sqlOpMap["Ne"] = operator{name: "Ne", sign: "<>"}
	sqlOpMap["Not"] = operator{name: "Not", sign: "!="}
	sqlOpMap["Gt"] = operator{name: "Gt", sign: ">"}
	sqlOpMap["Ge"] = operator{name: "Ge", sign: ">="}
	sqlOpMap["Lt"] = operator{name: "Lt", sign: "<"}
	sqlOpMap["Le"] = operator{name: "Le", sign: "<="}
	sqlOpMap["In"] = operator{name: "In", sign: "IN", format: "\tconditions = append(conditions, \"%s%s\"+strings.Repeat(\"?\", len(*q.%s)))"}
	sqlOpMap["NotIn"] = operator{name: "NotIn", sign: "NOT IN", format: "\tconditions = append(conditions, \"%s%s\" + strings.Repeat(\"?\", len(*q.%s)))"}
	sqlOpMap["Null"] = operator{name: "Null", sign: "IS NULL", format: "\tconditions = append(conditions, \"%s %s\")"}
	sqlOpMap["NotNull"] = operator{name: "NotNull", sign: "IS NOT NULL", format: "\tconditions = append(conditions, \"%s %s\")"}
	sqlOpMap["Like"] = operator{name: "Like", sign: "LIKE", format: "\tconditions = append(conditions, \"%s %s ?\")"}
	opMap["sql"] = sqlOpMap
}

func NewSqlGenerator() *SqlGenerator {
	return &SqlGenerator{&generator{
		Buffer:     bytes.NewBuffer(make([]byte, 0, 1024)),
		key:        "sql",
		imports:    []string{`"strings"`},
		bodyFormat: "\tconditions = append(conditions, \"%s %s ?\")",
		ifFormat:   "if q.%s%s {",
	}}
}

func (g *SqlGenerator) appendBuildMethod(ts *ast.TypeSpec) {
	g.WriteString(fmt.Sprintf("func (q %s) BuildConditions() ([]string, []any) {", ts.Name))
	g.WriteString(NewLine)
	g.WriteString("\tconditions := make([]string, 0, 4)")
	g.WriteString(NewLine)
	g.WriteString("\targs := make([]any, 0, 4)")
	g.WriteString(NewLine)
	g.appendStruct(ts.Type.(*ast.StructType), []string{})
	g.WriteString("\treturn conditions, args")
	g.WriteString(NewLine)
	g.WriteString("}")
	g.WriteString(NewLine)
}

func (g *SqlGenerator) appendStruct(stp *ast.StructType, path []string) {
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(toStructPointer(field), path, field.Names[0].Name)
		}
	}
}

func (g *SqlGenerator) appendCondition(stp *ast.StructType, path []string, fieldName string) {
	g.intent = strings.Repeat("\t", len(path)+1)

	column, op := g.suffixMatch(fieldName)

	if stp != nil {
		g.appendIfStartNil(fieldName)
	} else if strings.Contains(op.sign, "NULL") {
		g.appendIfStart(fieldName, "")
		g.appendIfBody(op.format, column, op.sign)
	} else if strings.Contains(op.sign, "IN") {
		g.appendIfStartNil(fieldName)
		g.appendIfBody(op.format, column, op.sign, fieldName)
		g.appendArgs(fieldName)
	} else {
		g.appendIfStartNil(fieldName)
		g.appendIfBody(op.format, column, op.sign)
		g.appendArgs(fieldName)
	}
	g.appendIfEnd()
}

func (g *SqlGenerator) appendArgs(name string) {
	g.WriteString(g.intent)
	g.WriteString(fmt.Sprintf("\targs = append(args, q.%s)", name))
	g.WriteString(NewLine)
}
