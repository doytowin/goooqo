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
		Buffer:  bytes.NewBuffer(make([]byte, 0, 1024)),
		imports: []string{`import . "go.mongodb.org/mongo-driver/bson/primitive"`},
	}}
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

	column, op := suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = buildNestedProperty(path, column)

	if stp != nil {
		g.appendIfStart(intent, structName, " != nil")
		g.appendStruct(stp, append(path, fieldName))
	} else if op.sign["mongo"] == "$type" {
		g.appendIfStart(intent, structName, "")
		g.WriteString(intent)
		g.WriteString(fmt.Sprintf(op.format, column, op.sign["mongo"]))
		g.WriteString(NewLine)
	} else {
		g.appendIfStart(intent, structName, " != nil")
		g.appendIfBody(intent, column, op, structName)
	}
	g.appendIfEnd(intent)
}

func (g *MongoGenerator) appendIfStart(intent string, structName string, cond string) {
	g.WriteString(intent)
	g.WriteString(fmt.Sprintf("if q.%s%s {", structName, cond))
	g.WriteString(NewLine)
}

func (g *MongoGenerator) appendIfBody(intent string, column string, op operator, structName string) {
	format := op.format
	if format == "" {
		format = "\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})"
	}
	g.WriteString(intent)
	g.WriteString(fmt.Sprintf(format, column, op.sign["mongo"], structName))
	g.WriteString(NewLine)
}
