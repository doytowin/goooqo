/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test

import . "github.com/doytowin/goooqo/core"

type MenuEntity struct {
	IntId
	ParentId *int    `json:"parentId,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type MenuQuery struct {
	PageQuery
	Id *int

	// Query the submenus of a specific parent menu:
	// parent_id IN (SELECT id FROM t_menu WHERE [conditions])
	Parent *MenuQuery `entitypath:"menu->ParentId,menu"`

	// Query the parent menu of a specific submenu:
	// id IN (SELECT parent_id FROM t_menu WHERE [conditions])
	Children *MenuQuery `entitypath:"menu,menu<-ParentId"`

	/**
	Query the menus accessible to a specific user:
	id IN (
		SELECT menu_id FROM a_perm_and_menu WHERE perm_id IN (
			SELECT perm_id FROM a_role_and_perm WHERE role_id IN (
				SELECT role_id FROM a_user_and_role WHERE user_id IN (
					SELECT id FROM t_user WHERE [conditions]
				)
			)
		)
	)*/
	User *UserQuery `entitypath:"menu,perm,role,user"`
}
