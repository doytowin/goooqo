/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"fmt"
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

func RegisterVirtualEntity(ve string, target string) {
	m[ve] = target
}

func MapVirtualEntity(ve string) string {
	if target := m[ve]; target != "" {
		return target
	}
	return ve
}
