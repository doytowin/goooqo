package gen

import (
	"bytes"
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
		bodyFormat: "d = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
		ifFormat:   "if q.%s%s {",
	}}
}

const regexSign = "$regex"

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
	mongoOpMap["Null"] = operator{name: "Null", sign: "$type"}
	mongoOpMap["Contain"] = operator{
		name:   "Contain",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	}
	mongoOpMap["NotContain"] = operator{
		name:   "NotContain",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", q.%s}}}}}})",
	}
	mongoOpMap["Start"] = operator{
		name:   "Start",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", \"^\" + *q.%s}}}})",
	}
	mongoOpMap["NotStart"] = operator{
		name:   "NotStart",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", \"^\" + *q.%s}}}}}})",
	}
	mongoOpMap["End"] = operator{
		name:   "End",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", *q.%s + \"$\"}}}})",
	}
	mongoOpMap["NotEnd"] = operator{
		name:   "NotEnd",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", *q.%s + \"$\"}}}}}})",
	}
	opMap["mongo"] = mongoOpMap
}

func (g *MongoGenerator) appendBuildMethod(ts *ast.TypeSpec) {
	g.writeInstruction("func (q %s) BuildFilter() A {", ts.Name)
	g.appendFuncBody(ts)
	g.writeInstruction("}")
}

func (g *MongoGenerator) appendFuncBody(ts *ast.TypeSpec) {
	g.appendIfBody("d := make(A, 0, 4)")
	g.appendStruct(ts.Type.(*ast.StructType), []string{})
	g.appendIfBody("return d")
}

func (g *MongoGenerator) appendStruct(stp *ast.StructType, path []string) {
	g.intent = strings.Repeat("\t", len(path)+1)
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(toStructPointer(field), path, field.Names[0].Name)
		}
	}
	g.intent = strings.Repeat("\t", len(path))
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
	column, op := g.suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = buildNestedProperty(path, column)

	if stp != nil {
		g.appendIfStartNil(structName)
		g.appendStruct(stp, append(path, fieldName))
	} else if op.sign == "$type" {
		g.appendIfStartNil(structName)
		g.writeInstruction("\tif *q.%s {", structName)
		g.writeInstruction("\t\td = append(d, D{{\"%s\", D{{\"$type\", 10}}}})", column)
		g.writeInstruction("\t} else {")
		g.writeInstruction("\t\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"$type\", 10}}}}}})", column)
		g.writeInstruction("\t}")
	} else if op.sign == regexSign {
		g.writeInstruction("if q.%s != nil && *q.%s != \"\" {", structName, structName)
		g.appendIfBody(op.format, column, op.sign, structName)
	} else {
		g.appendIfStartNil(structName)
		g.appendIfBody(op.format, column, op.sign, structName)
	}
	g.appendIfEnd()
}
