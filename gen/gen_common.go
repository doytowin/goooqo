/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package gen

import (
	"fmt"
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

func GenerateQueryBuilder(gen Generator, inputFilepath string, outputFilepath string) error {
	code := GenerateCode(inputFilepath, gen)
	return WriteFile(outputFilepath, code)
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
	for _, ts := range tsList {
		gen.addStruct("", ts)
	}

	gen.appendPackage(f.Name.String())
	gen.appendImports()
	for ts := gen.nextStruct(); ts != nil; ts = gen.nextStruct() {
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
						fieldType := fmt.Sprint(fields[0].Type)
						if strings.Contains(fieldType, "PageQuery") {
							result = append(result, ts)
						}
					}
				}
			}
		}
	}
	return
}

func toTypePointer(field *ast.Field) *ast.TypeSpec {
	if expr, ok := field.Type.(*ast.StarExpr); ok {
		if ident, ok := expr.X.(*ast.Ident); ok && ident.Obj != nil {
			if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
				return ts
			}
		}
	}
	return nil
}
