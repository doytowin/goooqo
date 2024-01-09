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
		g.WriteString(NewLine)
	}
}

func (g *generator) appendIfEnd(intent string) {
	g.WriteString(intent)
	g.WriteString("}")
	g.WriteString(NewLine)
}

func (g *generator) appendIfStart(intent string, structName string, cond string) {
	g.WriteString(intent)
	g.WriteString(fmt.Sprintf(g.ifFormat, structName, cond))
	g.WriteString(NewLine)
}

func (g *generator) appendIfBody(intent string, format string, args ...any) {
	if format == "" {
		format = g.bodyFormat
	}
	g.WriteString(intent)
	g.WriteString(fmt.Sprintf(format, args...))
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
