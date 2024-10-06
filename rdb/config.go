/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"fmt"
	"reflect"
	"strings"
)

var Config = struct {
	TableFormat     string
	JoinIdFormat    string
	JoinTableFormat string
}{
	"t_%s",
	"%s_id",
	"a_%s_and_%s",
}

var m = map[string]string{}

func FormatTable(domain string) string {
	return fmt.Sprintf(Config.TableFormat, domain)
}

func FormatTableByEntity(entity any) string {
	if rdbEntity, ok := entity.(RdbEntity); ok {
		return rdbEntity.GetTableName()
	}
	name := reflect.ValueOf(entity).Type().Name()
	name = strings.ToLower(strings.TrimSuffix(name, "Entity"))
	return fmt.Sprintf(Config.TableFormat, name)
}

func FormatJoinId(domain string) string {
	return fmt.Sprintf(Config.JoinIdFormat, domain)
}

func FormatJoinTable(domain1 string, domain2 string) string {
	if table := m[domain1+"_"+domain2]; table != "" {
		return table
	}
	return fmt.Sprintf(Config.JoinTableFormat, domain1, domain2)
}

func RegisterJoinTable(domain1 string, domain2 string, table string) {
	m[domain1+"_"+domain2] = table
}
