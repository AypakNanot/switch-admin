package datamodel

import (
	"log"
	"os"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

// BootstrapFilePath 指定 Bootstrap 文件路径
const BootstrapFilePath = "internal/datamodel/bootstrap.go"

// Bootstrap 初始化函数
// 初始化数据库表和菜单数据
func Bootstrap(e *engine.Engine) {
	connection := e.SqliteConnection()
	InitDatabaseTables(connection)
	InitMenu(connection)
	InitDashboard(connection)
}

// InitDatabaseTables 初始化 GoAdmin 所需的数据库表
func InitDatabaseTables(conn db.Connection) {
	// 检查表是否已存在
	exists, _ := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='goadmin_users'")
	if exists != nil && len(exists) > 0 {
		log.Println("数据库表已存在，跳过初始化")
		return
	}

	// 读取并执行初始化 SQL 脚本
	sqlBytes, err := os.ReadFile("data/init.sql")
	if err != nil {
		log.Printf("读取初始化 SQL 脚本失败：%v", err)
		return
	}

	// 执行 SQL 脚本（按分号分割）
	statements := string(sqlBytes)
	_, err = conn.Exec(statements)
	if err != nil {
		log.Printf("执行初始化 SQL 失败：%v", err)
		return
	}

	log.Println("数据库表初始化完成")
}

// InitMenu 初始化网络管理菜单
func InitMenu(conn db.Connection) {
	// 检查菜单是否已初始化
	exists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '网络管理' LIMIT 1")
	if exists != nil && len(exists) > 0 {
		log.Println("菜单已初始化，跳过")
		return
	}

	// 先插入主菜单获取 ID（type=0 表示目录，只在侧边栏显示不在顶部显示；order 是保留字，需要用反引号括起来）
	networkResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '网络管理', '', 'fa fa-network-wired', 1)")
	if err != nil {
		log.Printf("插入网络管理菜单失败：%v", err)
		return
	}
	networkId, _ := networkResult.LastInsertId()

	systemResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '系统配置', '', 'fa fa-cog', 2)")
	if err != nil {
		log.Printf("插入系统配置菜单失败：%v", err)
		return
	}
	systemId, _ := systemResult.LastInsertId()

	// 插入子菜单
	subMenus := []struct {
		parentId int64
		name     string
		uri      string
		icon     string
		order    int
	}{
		{networkId, "IPv4 路由表", "/network/route-table", "fa fa-table", 1},
		{networkId, "IPv4 静态路由", "/network/static-route", "fa fa-route", 2},
		{networkId, "Ping 诊断", "/network/ping", "fa fa-pingpong-paddle", 3},
		{networkId, "Traceroute 诊断", "/network/traceroute", "fa fa-share-alt", 4},
		{networkId, "虚拟电缆检测", "/network/cable-test", "fa fa-plug", 5},
		{systemId, "系统配置", "/system/config", "fa fa-wrench", 1},
	}

	for _, menu := range subMenus {
		_, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (?, 1, ?, ?, ?, ?)",
			menu.parentId, menu.name, menu.uri, menu.icon, menu.order)
		if err != nil {
			log.Printf("插入菜单 %s 失败：%v", menu.name, err)
		} else {
			log.Printf("插入菜单：%s (parent_id=%d)", menu.name, menu.parentId)
		}
	}

	// 添加角色菜单关联（让管理员角色可以看到这些菜单）
	// 首先获取管理员角色 ID（通常为 1）
	adminRoleId := int64(1)

	// 获取网络管理主菜单 ID
	networkMenuExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '网络管理' LIMIT 1")
	if len(networkMenuExists) > 0 {
		networkMenuId := int(networkMenuExists[0]["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, networkMenuId)

		// 添加子菜单关联
		networkSubMenus, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", networkMenuId)
		for _, row := range networkSubMenus {
			menuId := int(row["id"].(int64))
			conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
		}
	}

	// 获取系统配置主菜单 ID
	systemMenuExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '系统配置' LIMIT 1")
	if len(systemMenuExists) > 0 {
		systemMenuId := int(systemMenuExists[0]["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, systemMenuId)

		// 添加子菜单关联
		systemSubMenus, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", systemMenuId)
		for _, row := range systemSubMenus {
			menuId := int(row["id"].(int64))
			conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
		}
	}

	log.Println("菜单初始化完成")
}

// InitDashboard 初始化 Dashboard 菜单
func InitDashboard(conn db.Connection) {
	// 检查是否已有 Dashboard 菜单
	exists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = 'Dashboard' LIMIT 1")
	if exists == nil || len(exists) == 0 {
		// 插入 Dashboard 菜单（作为根菜单，order 为 0）
		_, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, title, uri, icon, `order`) VALUES (0, 'Dashboard', '/dashboard', 'fa fa-dashboard', 0)")
		if err != nil {
			log.Printf("插入 Dashboard 菜单失败：%v", err)
		} else {
			log.Println("插入 Dashboard 菜单")
		}
	}
}
