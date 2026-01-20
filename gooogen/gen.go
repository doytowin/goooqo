/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	goFile := os.Getenv("GOFILE")

	generatorType := flag.String("type", "sql", "(Optional) Generator type: sql, mongodb")
	inputFile := flag.String("f", goFile, "(Optional) The Go file containing the query definition")
	outputFile := flag.String("o", "", "(Optional) The Go file to output the query builder")

	flag.Parse()

	if *inputFile == "" {
		log.Fatalf("Input file is not specified")
	}

	log.Infof("Running command %s on %s", os.Args[0], *inputFile)

	var gen Generator
	switch *generatorType {
	case "sql":
		gen = NewSqlGenerator()
	case "mongodb":
		gen = NewMongoGenerator()
	default:
		log.Fatalf("Unsupported generator type: %s", *generatorType)
	}

	if *outputFile == "" {
		*outputFile = strings.ReplaceAll(*inputFile, ".go", "_query_builder.go")
	}

	err := GenerateQueryBuilder(gen, *inputFile, *outputFile)
	if err != nil {
		log.Fatalf("Error generating query builder: %v", err)
	}

	log.Infof("Query builder generated successfully to %s", *outputFile)
}
