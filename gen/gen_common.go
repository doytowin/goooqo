package gen

import (
	"bytes"
	"fmt"
	"github.com/doytowin/goooqo/core"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"runtime"
	"strings"
)

var NewLine = func() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}()

func Generate(input, output string) error {
	return WriteFile(output, GenerateCode(input))
}

func WriteFile(filename string, code string) error {
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	_, err := io.WriteString(f, code)
	return err
}

func GenerateCode(filename string) string {
	// Create the AST by parsing src.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}
	stpList := lookupQueryStruct(f)

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	appendPackage(buffer, f.Name.String())
	appendImports(buffer)
	for _, stp := range stpList {
		appendBuildMethod(buffer, stp)
	}
	return buffer.String()
}

func lookupQueryStruct(f *ast.File) (result []*ast.StructType) {
	for _, v := range f.Decls {
		if stc, ok := v.(*ast.GenDecl); ok && stc.Tok == token.TYPE {
			for _, spec := range stc.Specs {
				if tp, ok := spec.(*ast.TypeSpec); ok {
					if stp, ok := tp.Type.(*ast.StructType); ok && stp.Struct.IsValid() {
						fields := stp.Fields.List
						if fmt.Sprint(fields[0].Type) == "&{goooqo PageQuery}" {
							result = append(result, stp)
						}
					}
				}
			}
		}
	}
	return
}

func appendPackage(buffer *bytes.Buffer, pkg string) {
	buffer.WriteString("package " + pkg)
	buffer.WriteString(NewLine)
	buffer.WriteString(NewLine)
}

func suffixMatch(fieldName string) (string, operator) {
	if match := suffixRgx.FindStringSubmatch(fieldName); len(match) > 0 {
		op := opMap[match[1]]
		column := strings.TrimSuffix(fieldName, match[1])
		column = core.ConvertToColumnCase(column)
		return column, op
	}
	return core.ConvertToColumnCase(fieldName), opMap["Eq"]
}
