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
	InitMaintenanceMenu(connection) // 单独添加维护菜单
	InitConfigMenu(connection)      // 添加配置模块菜单
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

// InitMenu 初始化网络管理和维护菜单
func InitMenu(conn db.Connection) {
	// 清除旧的菜单数据（以便重新初始化）
	conn.Exec("DELETE FROM goadmin_menu WHERE parent_id IN (SELECT id FROM goadmin_menu WHERE title = '网络管理')")
	conn.Exec("DELETE FROM goadmin_menu WHERE title = '网络管理'")
	conn.Exec("DELETE FROM goadmin_menu WHERE parent_id IN (SELECT id FROM goadmin_menu WHERE title = '系统配置')")
	conn.Exec("DELETE FROM goadmin_menu WHERE title = '系统配置'")
	conn.Exec("DELETE FROM goadmin_role_menu WHERE menu_id IN (SELECT id FROM goadmin_menu WHERE title = '网络管理' OR title = '系统配置')")

	// 分别检查各个主菜单是否已初始化
	networkExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '网络管理' LIMIT 1")
	systemExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '系统配置' LIMIT 1")
	maintenanceExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1")

	// 如果所有菜单都已初始化，直接返回
	if (networkExists != nil && len(networkExists) > 0) &&
		(systemExists != nil && len(systemExists) > 0) &&
		(maintenanceExists != nil && len(maintenanceExists) > 0) {
		log.Println("菜单已初始化，跳过")
		return
	}

	var networkId int64
	var systemId int64
	var maintenanceId int64

	// 插入网络管理主菜单（如果不存在）
	if networkExists == nil || len(networkExists) == 0 {
		networkResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '网络管理', '', 'fa fa-network-wired', 1)")
		if err != nil {
			log.Printf("插入网络管理菜单失败：%v", err)
		} else {
			networkId, _ = networkResult.LastInsertId()
			log.Println("插入网络管理主菜单")
		}
	} else {
		networkId = networkExists[0]["id"].(int64)
	}

	// 插入系统配置主菜单（如果不存在）
	if systemExists == nil || len(systemExists) == 0 {
		systemResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '系统配置', '', 'fa fa-cog', 2)")
		if err != nil {
			log.Printf("插入系统配置菜单失败：%v", err)
		} else {
			systemId, _ = systemResult.LastInsertId()
			log.Println("插入系统配置主菜单")
		}
	} else {
		systemId = systemExists[0]["id"].(int64)
	}

	// 插入维护主菜单（如果不存在）
	if maintenanceExists == nil || len(maintenanceExists) == 0 {
		maintenanceResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '维护', '', 'fa fa-wrench', 3)")
		if err != nil {
			log.Printf("插入维护菜单失败：%v", err)
		} else {
			maintenanceId, _ = maintenanceResult.LastInsertId()
			log.Println("插入维护主菜单")
		}
	} else {
		maintenanceId = maintenanceExists[0]["id"].(int64)
	}

	// 插入子菜单（只插入不存在的）
	subMenus := []struct {
		parentId int64
		name     string
		uri      string
		icon     string
		order    int
	}{
		// 网络管理子菜单
		{networkId, "IPv4 路由表", "/network/route-table", "fa fa-list-alt", 1},
		{networkId, "IPv4 静态路由", "/network/static-route", "fa fa-road", 2},
		{networkId, "Ping 诊断", "/network/ping", "fa fa-exchange", 3},
		{networkId, "Traceroute 诊断", "/network/traceroute", "fa fa-random", 4},
		{networkId, "虚拟电缆检测", "/network/cable-test", "fa fa-wrench", 5},
		// 系统配置子菜单
		{systemId, "系统配置", "/system/config", "fa fa-wrench", 1},
		// 维护子菜单
		{maintenanceId, "重启/保存", "/maintenance/reboot-save", "fa fa-power-off", 1},
		{maintenanceId, "用户管理", "/maintenance/users", "fa fa-users", 2},
		{maintenanceId, "系统配置", "/maintenance/system-config", "fa fa-cogs", 3},
		{maintenanceId, "加载配置", "/maintenance/load-config", "fa fa-upload", 4},
		{maintenanceId, "文件管理", "/maintenance/files", "fa fa-file", 5},
		{maintenanceId, "日志管理", "/maintenance/logs", "fa fa-history", 6},
		{maintenanceId, "SNMP 配置", "/maintenance/snmp", "fa fa-bell", 7},
		{maintenanceId, "SNMP Trap 配置", "/maintenance/snmp-trap", "fa fa-exclamation-triangle", 8},
		{maintenanceId, "蠕虫攻击防护", "/maintenance/worm-protection", "fa fa-bug", 9},
		{maintenanceId, "DDoS 攻击防护", "/maintenance/ddos-protection", "fa fa-shield", 10},
		{maintenanceId, "ARP 攻击防护", "/maintenance/arp-protection", "fa fa-lock", 11},
		{maintenanceId, "当前会话", "/maintenance/sessions", "fa fa-clock-o", 12},
	}

	for _, menu := range subMenus {
		// 检查子菜单是否已存在
		exists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = ? AND parent_id = ? LIMIT 1", menu.name, menu.parentId)
		if exists != nil && len(exists) > 0 {
			continue // 已存在，跳过
		}
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

	// 获取维护主菜单 ID
	maintenanceMenuExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1")
	if len(maintenanceMenuExists) > 0 {
		maintenanceMenuId := int(maintenanceMenuExists[0]["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, maintenanceMenuId)

		// 添加子菜单关联
		maintenanceSubMenus, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", maintenanceMenuId)
		for _, row := range maintenanceSubMenus {
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

// InitMaintenanceMenu 单独初始化维护菜单（即使其他菜单已存在）
func InitMaintenanceMenu(conn db.Connection) {
	// 检查维护主菜单是否已存在
	maintenanceExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1")
	if maintenanceExists != nil && len(maintenanceExists) > 0 {
		// 维护菜单已存在，检查是否需要添加子菜单关联到角色菜单
		maintenanceMenuId := int(maintenanceExists[0]["id"].(int64))
		adminRoleId := int64(1)

		// 检查主菜单是否已关联到角色
		roleMenuExists, _ := conn.Query("SELECT id FROM goadmin_role_menu WHERE role_id = ? AND menu_id = ? LIMIT 1", adminRoleId, maintenanceMenuId)
		if roleMenuExists == nil || len(roleMenuExists) == 0 {
			conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, maintenanceMenuId)
			log.Println("添加维护主菜单到角色菜单关联")
		}

		// 检查子菜单是否已关联到角色
		subMenus, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", maintenanceMenuId)
		for _, row := range subMenus {
			menuId := int(row["id"].(int64))
			roleMenuExists, _ := conn.Query("SELECT id FROM goadmin_role_menu WHERE role_id = ? AND menu_id = ? LIMIT 1", adminRoleId, menuId)
			if roleMenuExists == nil || len(roleMenuExists) == 0 {
				conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
			}
		}
		return // 已存在，直接返回
	}

	// 插入维护主菜单
	maintenanceResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '维护', '', 'fa fa-wrench', 3)")
	if err != nil {
		log.Printf("插入维护菜单失败：%v", err)
		return
	}
	maintenanceId, _ := maintenanceResult.LastInsertId()
	log.Println("插入维护主菜单")

	// 插入维护子菜单
	subMenus := []struct {
		name  string
		uri   string
		icon  string
		order int
	}{
		{"重启/保存", "/maintenance/reboot-save", "fa fa-power-off", 1},
		{"用户管理", "/maintenance/users", "fa fa-users", 2},
		{"系统配置", "/maintenance/system-config", "fa fa-cogs", 3},
		{"加载配置", "/maintenance/load-config", "fa fa-upload", 4},
		{"文件管理", "/maintenance/files", "fa fa-file", 5},
		{"日志管理", "/maintenance/logs", "fa fa-history", 6},
		{"SNMP 配置", "/maintenance/snmp", "fa fa-bell", 7},
		{"SNMP Trap 配置", "/maintenance/snmp-trap", "fa fa-exclamation-triangle", 8},
		{"蠕虫攻击防护", "/maintenance/worm-protection", "fa fa-bug", 9},
		{"DDoS 攻击防护", "/maintenance/ddos-protection", "fa fa-shield", 10},
		{"ARP 攻击防护", "/maintenance/arp-protection", "fa fa-lock", 11},
		{"当前会话", "/maintenance/sessions", "fa fa-clock-o", 12},
	}

	for _, menu := range subMenus {
		_, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (?, 1, ?, ?, ?, ?)",
			maintenanceId, menu.name, menu.uri, menu.icon, menu.order)
		if err != nil {
			log.Printf("插入菜单 %s 失败：%v", menu.name, err)
		} else {
			log.Printf("插入菜单：%s (parent_id=%d)", menu.name, maintenanceId)
		}
	}

	// 添加角色菜单关联
	adminRoleId := int64(1)
	conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, maintenanceId)

	// 添加子菜单关联
	subMenusResult, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", maintenanceId)
	for _, row := range subMenusResult {
		menuId := int(row["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
	}

	log.Println("维护菜单初始化完成")
}

// InitConfigMenu 初始化配置模块菜单
func InitConfigMenu(conn db.Connection) {
	// 清除旧的配置菜单数据（以便重新初始化）
	conn.Exec("DELETE FROM goadmin_menu WHERE parent_id IN (SELECT id FROM goadmin_menu WHERE title = '配置')")
	conn.Exec("DELETE FROM goadmin_menu WHERE title = '配置'")
	conn.Exec("DELETE FROM goadmin_role_menu WHERE menu_id IN (SELECT id FROM goadmin_menu WHERE title = '配置' OR parent_id IN (SELECT id FROM goadmin_menu WHERE title = '配置'))")

	// 检查配置主菜单是否已存在
	configExists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = '配置' LIMIT 1")

	var configId int64
	if configExists == nil || len(configExists) == 0 {
		// 插入配置主菜单
		configResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '配置', '', 'fa fa-cogs', 4)")
		if err != nil {
			log.Printf("插入配置菜单失败：%v", err)
			return
		}
		configId, _ = configResult.LastInsertId()
		log.Println("插入配置主菜单")
	} else {
		configId = configExists[0]["id"].(int64)
	}

	// 配置模块子菜单 - 第一阶段核心功能
	subMenus := []struct {
		name  string
		uri   string
		icon  string
		order int
	}{
		// 第一阶段：核心基础功能
		{"端口状态", "/config/ports", "fa fa-list", 1},
		{"链路聚合", "/config/link-aggregation", "fa fa-chain", 2},
		{"风暴控制", "/config/storm-control", "fa fa-bolt", 3},
		{"流量控制", "/config/flow-control", "fa fa-tachometer", 4},
		{"端口隔离", "/config/port-isolation", "fa fa-th", 5},
		{"端口监测", "/config/port-monitor", "fa fa-desktop", 6},
		{"VLAN", "/config/vlan", "fa fa-sitemap", 7},
		{"MAC 地址表", "/config/mac-table", "fa fa-table", 8},
		{"生成树", "/config/stp", "fa fa-project-diagram", 9},
		{"ERPS", "/config/erps", "fa fa-circle-notch", 10},
		{"PoE", "/config/poe", "fa fa-plug", 11},
		{"端口镜像", "/config/port-mirror", "fa fa-clone", 12},
		{"组播", "/config/multicast", "fa fa-rss", 13},
		{"资源", "/config/resource", "fa fa-microchip", 14},
		{"堆叠", "/config/stack", "fa fa-cubes", 15},
	}

	for _, menu := range subMenus {
		// 检查子菜单是否已存在
		exists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = ? AND parent_id = ? LIMIT 1", menu.name, configId)
		if exists != nil && len(exists) > 0 {
			continue
		}
		_, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (?, 1, ?, ?, ?, ?)",
			configId, menu.name, menu.uri, menu.icon, menu.order)
		if err != nil {
			log.Printf("插入菜单 %s 失败：%v", menu.name, err)
		} else {
			log.Printf("插入菜单：%s (parent_id=%d)", menu.name, configId)
		}
	}

	// 添加角色菜单关联
	adminRoleId := int64(1)
	conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, configId)

	// 添加子菜单关联
	subMenusResult, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", configId)
	for _, row := range subMenusResult {
		menuId := int(row["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
	}

	log.Println("配置菜单初始化完成")
}
