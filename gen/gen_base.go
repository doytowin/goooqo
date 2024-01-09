package gen

import (
	"bytes"
	"fmt"
	"github.com/doytowin/goooqo/core"
	"go/ast"
	"strings"
)

type Generator interface {
	appendPackage(string2 string)
	appendImports()
	appendBuildMethod(ts *ast.TypeSpec)
	String() string
}

type generator struct {
	*bytes.Buffer
	key        string
	imports    []string
	bodyFormat string
	ifFormat   string
	intent     string
}

func (g *generator) appendPackage(pkg string) {
	g.WriteString("package " + pkg)
	g.WriteString(NewLine)
	g.WriteString(NewLine)
}

func (g *generator) appendImports() {
	for _, s := range g.imports {
		g.WriteString("import " + s)
		g.WriteString(NewLine)
	}
	g.WriteString(NewLine)
}

func (g *generator) appendIfEnd() {
	g.WriteString(g.intent)
	g.WriteString("}")
	g.WriteString(NewLine)
}

func (g *generator) appendIfStart(structName string, cond string) {
	g.WriteString(g.intent)
	g.WriteString(fmt.Sprintf(g.ifFormat, structName, cond))
	g.WriteString(NewLine)
}

func (g *generator) appendIfStartNil(fieldName string) {
	g.appendIfStart(fieldName, " != nil")
}

func (g *generator) appendIfBody(ins string, args ...any) {
	if ins == "" {
		ins = g.bodyFormat
	}
	g.WriteString(g.intent)
	g.WriteString("\t")
	g.WriteString(fmt.Sprintf(ins, args...))
	g.WriteString(NewLine)
}

func (g *generator) writeInstruction(ins string, args ...any) {
	g.WriteString(g.intent)
	g.WriteString(fmt.Sprintf(ins, args...))
	g.WriteString(NewLine)
}

func (g *generator) suffixMatch(fieldName string) (string, operator) {
	if match := suffixRgx.FindStringSubmatch(fieldName); len(match) > 0 {
		op := opMap[g.key][match[1]]
		column := strings.TrimSuffix(fieldName, match[1])
		column = core.ConvertToColumnCase(column)
		return column, op
	}
	return core.ConvertToColumnCase(fieldName), opMap[g.key]["Eq"]
}
