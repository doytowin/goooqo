package gen

import (
	"bytes"
	"fmt"
	"go/ast"
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
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			fieldName := field.Names[0].Name
			appendCondition(buffer, fieldName)
		}
	}
	buffer.WriteString("\treturn d")
	buffer.WriteString(NewLine)
	buffer.WriteString("}")
	buffer.WriteString(NewLine)
}

func appendCondition(buffer *bytes.Buffer, fieldName string) {
	column, op := suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	buffer.WriteString(fmt.Sprintf("\tif q.%s != nil {", fieldName))
	buffer.WriteString(NewLine)
	buffer.WriteString(fmt.Sprintf("\t\td = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})", column, op.sign["mongo"], fieldName))
	buffer.WriteString(NewLine)
	buffer.WriteString("\t}")
	buffer.WriteString(NewLine)
}
