package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"runtime"
)

var NewLine = func() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}()

func Generate(code, output string) error {
	return WriteFile(output, code)
}

func WriteFile(filename string, code string) error {
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	_, err := io.WriteString(f, code)
	return err
}

func GenerateCode(filename string, gen Generator) string {
	// Create the AST by parsing src.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}
	tsList := lookupQueryStruct(f)

	gen.appendPackage(f.Name.String())
	gen.appendImports()
	for _, ts := range tsList {
		gen.appendBuildMethod(ts)
	}
	return gen.String()
}

func lookupQueryStruct(f *ast.File) (result []*ast.TypeSpec) {
	for _, v := range f.Decls {
		if stc, ok := v.(*ast.GenDecl); ok && stc.Tok == token.TYPE {
			for _, spec := range stc.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if stp, ok := ts.Type.(*ast.StructType); ok && stp.Struct.IsValid() {
						fields := stp.Fields.List
						if fmt.Sprint(fields[0].Type) == "&{goooqo PageQuery}" {
							result = append(result, ts)
						}
					}
				}
			}
		}
	}
	return
}

func toStructPointer(field *ast.Field) *ast.StructType {
	if expr, ok := field.Type.(*ast.StarExpr); ok {
		if ident, ok := expr.X.(*ast.Ident); ok && ident.Obj != nil {
			if tp, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
				if stp, ok := tp.Type.(*ast.StructType); ok && stp.Struct.IsValid() {
					return stp
				}
			}
		}
	}
	return nil
}
