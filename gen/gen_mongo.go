package gen

import (
	"bytes"
	"fmt"
	"github.com/doytowin/goooqo/core"
	"go/ast"
	"strings"
)

func appendImports(buffer *bytes.Buffer) {
	buffer.WriteString(`import . "go.mongodb.org/mongo-driver/bson/primitive"`)
	buffer.WriteString(NewLine)
	buffer.WriteString(NewLine)
}

func appendBuildMethod(buffer *bytes.Buffer, stp *ast.StructType) {
	buffer.WriteString("func (q InventoryQuery) BuildFilter() []D {")
	buffer.WriteString(NewLine)
	buffer.WriteString("\td := make([]D, 0, 10)")
	buffer.WriteString(NewLine)
	appendStruct(buffer, stp, []string{})
	buffer.WriteString("\treturn d")
	buffer.WriteString(NewLine)
	buffer.WriteString("}")
	buffer.WriteString(NewLine)
}

func appendStruct(buffer *bytes.Buffer, stp *ast.StructType, path []string) {
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			appendCondition(buffer, toStructPointer(field), path, field.Names[0].Name)
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

func appendCondition(buffer *bytes.Buffer, stp *ast.StructType, path []string, fieldName string) {
	intent := strings.Repeat("\t", len(path)+1)

	column, op := suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = buildNestedProperty(path, column)

	if op.name == "Null" {
		appendIfStart(buffer, intent, structName, "")
		buffer.WriteString(intent)
		buffer.WriteString(fmt.Sprintf("\td = append(d, D{{\"%s\", D{{\"%s\", 10}}}})", column, op.sign["mongo"]))
		buffer.WriteString(NewLine)
	} else if op.name == "NotNull" {
		appendIfStart(buffer, intent, structName, "")
		buffer.WriteString(intent)
		buffer.WriteString(fmt.Sprintf("\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", 10}}}}}})", column, op.sign["mongo"]))
		buffer.WriteString(NewLine)
	} else {
		appendIfStart(buffer, intent, structName, " != nil")
		if stp != nil {
			appendStruct(buffer, stp, append(path, fieldName))
		} else {
			appendIfBody(buffer, intent, column, op, structName)
		}
	}
	appendIfEnd(buffer, intent)
}

func appendIfStart(buffer *bytes.Buffer, intent string, structName string, cond string) {
	buffer.WriteString(intent)
	buffer.WriteString(fmt.Sprintf("if q.%s%s {", structName, cond))
	buffer.WriteString(NewLine)
}

func appendIfBody(buffer *bytes.Buffer, intent string, column string, op operator, structName string) {
	buffer.WriteString(intent)
	buffer.WriteString(fmt.Sprintf("\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})", column, op.sign["mongo"], structName))
	buffer.WriteString(NewLine)
}

func appendIfEnd(buffer *bytes.Buffer, intent string) {
	buffer.WriteString(intent)
	buffer.WriteString("}")
	buffer.WriteString(NewLine)
}
