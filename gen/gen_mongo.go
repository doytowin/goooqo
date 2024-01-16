package gen

import (
	"github.com/doytowin/goooqo/core"
	"go/ast"
	"strings"
)

type MongoGenerator struct {
	*generator
}

func NewMongoGenerator() *MongoGenerator {
	return &MongoGenerator{newGenerator("mongo",
		[]string{`. "go.mongodb.org/mongo-driver/bson/primitive"`},
		"d = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	)}
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
			g.appendCondition(field, path, field.Names[0].Name)
		} else if tn := resolveTypeName(field.Type); strings.HasSuffix(tn, "Or") {
			g.appendCondition(field, path, strings.TrimPrefix(tn, "*"))
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
		if strings.HasSuffix(fieldName, "Or") {
			continue
		}
		// LATER: determine property by tag
		props = append(props, core.ConvertToColumnCase(fieldName))
	}
	return strings.Join(append(props, column), ".")
}

func (g *MongoGenerator) appendCondition(field *ast.Field, path []string, fieldName string) {
	column, op := g.suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = buildNestedProperty(path, column)

	stp := toStructPointer(field)
	if stp != nil {
		g.appendIfStartNil(structName)
		if lenP := len(path); lenP > 0 && strings.HasSuffix(path[lenP-1], "Or") {
			g.appendAndBody(stp, path, fieldName)
		} else if strings.HasSuffix(fieldName, "Or") {
			g.appendOrBody(stp, path, fieldName)
		} else {
			g.appendStruct(stp, append(path, fieldName))
		}
	} else if op.sign == "$type" {
		g.appendIfStartNil(structName)
		g.writeInstruction("\tif *q.%s {", structName)
		g.writeInstruction(g.replaceIns("\t\td = append(d, D{{\"%s\", D{{\"$type\", 10}}}})"), column)
		g.writeInstruction("\t} else {")
		g.writeInstruction(g.replaceIns("\t\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"$type\", 10}}}}}})"), column)
		g.writeInstruction("\t}")
	} else if op.sign == regexSign {
		g.writeInstruction("if q.%s != nil && *q.%s != \"\" {", structName, structName)
		g.appendIfBody(op.format, column, op.sign, structName)
	} else {
		g.appendIfStartNil(structName)
		if resolveTypeName(field.Type) == "*M" {
			g.appendIfBody("d = append(d, *q.%s)", structName)
		} else {
			g.appendIfBody(op.format, column, op.sign, structName)
		}
	}
	g.appendIfEnd()
}

func (g *MongoGenerator) appendOrBody(stp *ast.StructType, path []string, fieldName string) {
	g.writeInstruction("\tor := make(A, 0, 4)")
	g.replaceIns = func(ins string) string {
		return strings.ReplaceAll(ins, "d = append(d", "or = append(or")
	}
	g.appendStruct(stp, append(path, fieldName))
	g.replaceIns = keep
	g.writeInstruction("\tif len(or) > 1 {")
	g.writeInstruction("\t\td = append(d, D{{\"$or\", or}})")
	g.writeInstruction("\t} else if len(or) == 1 {")
	g.writeInstruction("\t\td = append(d, or[0])")
	g.writeInstruction("\t}")
}

func (g *MongoGenerator) appendAndBody(stp *ast.StructType, path []string, fieldName string) {
	g.writeInstruction("\tand := make(A, 0, 4)")
	g.replaceIns = func(ins string) string {
		return strings.ReplaceAll(ins, "d = append(d", "and = append(and")
	}
	g.appendStruct(stp, append(path, fieldName))
	g.replaceIns = keep
	g.writeInstruction("\tif len(and) > 1 {")
	g.writeInstruction("\t\tor = append(or, D{{\"$and\", and}})")
	g.writeInstruction("\t} else if len(and) == 1 {")
	g.writeInstruction("\t\tor = append(or, and[0])")
	g.writeInstruction("\t}")
}

func resolveTypeName(expr ast.Expr) string {
	var stack []string
loop:
	for {
		switch x := expr.(type) {
		case *ast.StarExpr:
			expr = x.X
			stack = append(stack, "*")
		case *ast.ArrayType:
			expr = x.Elt
			stack = append(stack, "[]")
		case *ast.Ident:
			stack = append(stack, x.Name)
			break loop
		case *ast.SelectorExpr:
			stack = append(stack, x.Sel.Name)
			break loop
		default:
			break loop
		}
	}
	return strings.Join(stack, "")
}
