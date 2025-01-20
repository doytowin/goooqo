/*
 * The Clear BSD License
 *
 * Copyright (c) 2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	inputFile := os.Getenv("GOFILE")
	log.Infof("Running command %s on %s", os.Args[0], inputFile)

	generatorType := flag.String("type", "sql", "Generator type: sql, mongodb")
	outputFile := flag.String("o", "", "The Go file to output the query builder")

	flag.Parse()

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
		*outputFile = strings.ReplaceAll(inputFile, ".go", "_query_builder.go")
	}

	err := GenerateQueryBuilder(gen, inputFile, *outputFile)
	if err != nil {
		log.Fatalf("Error generating query builder: %v", err)
	}

	log.Infof("Query builder generated successfully for %s", inputFile)
}
