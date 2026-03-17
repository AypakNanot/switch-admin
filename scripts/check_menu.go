package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite3", "data/admin.db")
	if err != nil {
		fmt.Println("Open database failed:", err)
		return
	}
	defer db.Close()

	fmt.Println("=== goadmin_menu 表中的所有菜单 ===")
	rows, err := db.Query("SELECT id, parent_id, type, title, uri, icon FROM goadmin_menu ORDER BY id")
	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, parentId, menuType int
		var title, uri, icon string
		rows.Scan(&id, &parentId, &menuType, &title, &uri, &icon)
		fmt.Printf("ID=%d, ParentID=%d, Type=%d, Title=%s, URI=%s, Icon=%s\n",
			id, parentId, menuType, title, uri, icon)
	}

	fmt.Println("\n=== goadmin_role_menu 表中的所有角色菜单关联 ===")
	rows2, err := db.Query("SELECT role_id, menu_id FROM goadmin_role_menu ORDER BY role_id, menu_id")
	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		var roleId, menuId int
		rows2.Scan(&roleId, &menuId)
		fmt.Printf("RoleID=%d, MenuID=%d\n", roleId, menuId)
	}
}
