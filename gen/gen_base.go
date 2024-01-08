package gen

import (
	"bytes"
)

type Generator interface {
	appendBuildMethod()
}

type generator struct {
	*bytes.Buffer
	imports []string
}

func (g *generator) appendPackage(pkg string) {
	g.WriteString("package " + pkg)
	g.WriteString(NewLine)
	g.WriteString(NewLine)
}

func (g *generator) appendImports() {
	for _, s := range g.imports {
		g.WriteString(s)
		g.WriteString(NewLine)
		g.WriteString(NewLine)
	}
}

func (g *generator) appendIfEnd(intent string) {
	g.WriteString(intent)
	g.WriteString("}")
	g.WriteString(NewLine)
}
