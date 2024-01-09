package gen

import (
	"bytes"
	"fmt"
	"github.com/doytowin/goooqo/core"
	"go/ast"
	"strings"
)

type MongoGenerator struct {
	*generator
}

func NewMongoGenerator() *MongoGenerator {
	return &MongoGenerator{&generator{
		Buffer:     bytes.NewBuffer(make([]byte, 0, 1024)),
		key:        "mongo",
		imports:    []string{`. "go.mongodb.org/mongo-driver/bson/primitive"`},
		bodyFormat: "\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
		ifFormat:   "if q.%s%s {",
	}}
}

func init() {
	mongoOpMap := make(map[string]operator)
	mongoOpMap["Eq"] = operator{name: "Eq", sign: "$eq"}
	mongoOpMap["Ne"] = operator{name: "Ne", sign: "$ne"}
	mongoOpMap["Not"] = operator{name: "Not", sign: "$ne"}
	mongoOpMap["Gt"] = operator{name: "Gt", sign: "$gt"}
	mongoOpMap["Ge"] = operator{name: "Ge", sign: "$gte"}
	mongoOpMap["Lt"] = operator{name: "Lt", sign: "$lt"}
	mongoOpMap["Le"] = operator{name: "Le", sign: "$lte"}
	mongoOpMap["In"] = operator{name: "In", sign: "$in"}
	mongoOpMap["NotIn"] = operator{name: "NotIn", sign: "$nin"}
	mongoOpMap["Null"] = operator{
		name:   "Null",
		sign:   "$type",
		format: "\td = append(d, D{{\"%s\", D{{\"%s\", 10}}}})",
	}
	mongoOpMap["NotNull"] = operator{
		name:   "NotNull",
		sign:   "$type",
		format: "\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", 10}}}}}})",
	}
	mongoOpMap["Contain"] = operator{
		name:   "Contain",
		sign:   "$regex",
		format: "\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	}
	mongoOpMap["NotContain"] = operator{
		name:   "NotContain",
		sign:   "$regex",
		format: "\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", q.%s}}}}}})",
	}
	opMap["mongo"] = mongoOpMap
}

func (g *MongoGenerator) appendBuildMethod(ts *ast.TypeSpec) {
	g.WriteString(fmt.Sprintf("func (q %s) BuildFilter() []D {", ts.Name))
	g.WriteString(NewLine)
	g.WriteString("\td := make([]D, 0, 4)")
	g.WriteString(NewLine)
	g.appendStruct(ts.Type.(*ast.StructType), []string{})
	g.WriteString("\treturn d")
	g.WriteString(NewLine)
	g.WriteString("}")
	g.WriteString(NewLine)
}

func (g *MongoGenerator) appendStruct(stp *ast.StructType, path []string) {
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(toStructPointer(field), path, field.Names[0].Name)
		}
	}
}

func buildNestedFieldName(path []string, fieldName string) string {
	return strings.Join(append(path, fieldName), ".")
}

func buildNestedProperty(path []string, column string) string {
	props := make([]string, 0, len(path)+1)
	for _, fieldName := range path {
		// LATER: determine property by tag
		props = append(props, core.ConvertToColumnCase(fieldName))
	}
	return strings.Join(append(props, column), ".")
}

func (g *MongoGenerator) appendCondition(stp *ast.StructType, path []string, fieldName string) {
	intent := strings.Repeat("\t", len(path)+1)

	column, op := g.suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = buildNestedProperty(path, column)

	if stp != nil {
		g.appendIfStart(intent, structName, " != nil")
		g.appendStruct(stp, append(path, fieldName))
	} else if op.sign == "$type" {
		g.appendIfStart(intent, structName, "")
		g.appendIfBody(intent, op.format, column, op.sign)
	} else {
		g.appendIfStart(intent, structName, " != nil")
		g.appendIfBody(intent, op.format, column, op.sign, structName)
	}
	g.appendIfEnd(intent)
}
