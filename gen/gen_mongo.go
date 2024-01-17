package gen

import (
	"github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"go/ast"
	"reflect"
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
	g.WriteString(NewLine)
	g.writeInstruction("func (q %s) BuildFilter() A {", ts.Name)
	g.appendFuncBody(ts)
	g.writeInstruction("}")
}

func (g *MongoGenerator) appendFuncBody(ts *ast.TypeSpec) {
	g.appendIfBody("d := make(A, 0, 4)")
	g.appendStruct(ts, []string{})
	g.appendIfBody("return d")
}

func (g *MongoGenerator) appendStruct(ts *ast.TypeSpec, path []string) {
	stp := ts.Type.(*ast.StructType)
	g.intent = strings.Repeat("\t", len(path)+1)
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(field, path, field.Names[0].Name)
		} else if tn := resolveTypeName(field.Type); strings.HasSuffix(tn, "Or") {
			g.appendCondition(field, path, strings.TrimPrefix(tn, "*"))
		} else {
			log.Info("[MongoGenerator#appendStruct] Unresolved TypeName: ", tn)
		}
	}
	g.intent = strings.Repeat("\t", len(path))
}

func buildNestedFieldName(path []string, fieldName string) string {
	return strings.Join(append(path, fieldName), ".")
}

func (g *MongoGenerator) buildNestedProperty(path []string, column string) string {
	props := make([]string, 0, len(path)+1)
	currentPrefix := g.prefix[g.structIdx-1]
	if currentPrefix != "" {
		props = append(props, currentPrefix)
	}
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
	column = g.buildNestedProperty(path, column)

	if field.Tag != nil {
		if columnTag, ok := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Lookup("column"); ok {
			column = columnTag
		}
	}

	if ts := toTypePointer(field); ts != nil {
		g.appendSubStruct(ts, structName, fieldName, path, column)
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

func (g *MongoGenerator) appendSubStruct(ts *ast.TypeSpec, structName string, fieldName string, path []string, column string) {
	g.appendIfStartNil(structName)
	if strings.HasSuffix(fieldName, "Or") {
		if lenP := len(path); lenP > 0 && strings.HasSuffix(path[lenP-1], "Or") {
			g.appendOrOrBody(path, fieldName)
		} else {
			g.appendOrBody(ts, path, fieldName)
		}
	} else {
		if lenP := len(path); lenP > 0 && strings.HasSuffix(path[lenP-1], "Or") {
			g.appendAndBody(path, fieldName)
		} else {
			g.addStruct(column, ts) // generate for type recursively
			g.writeInstruction("\td = append(d, q.%s.BuildFilter()...)", structName)
		}
	}
}

func (g *MongoGenerator) appendOrBody(ts *ast.TypeSpec, path []string, fieldName string) {
	g.writeInstruction("\tor := make(A, 0, 4)")
	g.replaceIns = func(ins string) string {
		return strings.ReplaceAll(ins, "d = append(d", "or = append(or")
	}
	g.appendStruct(ts, append(path, fieldName))
	g.replaceIns = keep
	g.writeInstruction("\tif len(or) > 1 {")
	g.writeInstruction("\t\td = append(d, D{{\"$or\", or}})")
	g.writeInstruction("\t} else if len(or) == 1 {")
	g.writeInstruction("\t\td = append(d, or[0])")
	g.writeInstruction("\t}")
}

func (g *MongoGenerator) appendOrOrBody(path []string, fieldName string) {
	structName := buildNestedFieldName(path, fieldName)
	g.writeInstruction("\tor = append(or, q.%s.BuildFilter()...)", structName)
}

func (g *MongoGenerator) appendAndBody(path []string, fieldName string) {
	structName := buildNestedFieldName(path, fieldName)
	g.writeInstruction("\tand := q.%s.BuildFilter()", structName)
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
