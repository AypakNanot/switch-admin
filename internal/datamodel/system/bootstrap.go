package system

import (
	"log"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

// BootstrapFilePath 指定 Bootstrap 文件路径
const BootstrapFilePath = "internal/datamodel/system/bootstrap.go"

// Bootstrap 初始化函数
// 初始化数据库表和菜单数据
func Bootstrap(e *engine.Engine) {
	connection := e.SqliteConnection()
	InitDatabaseTables(connection)
	InitMenu(connection)
	InitDashboard(connection)
}
func InitDatabaseTables(conn db.Connection) {
	// 检查表是否已存在
	exists, _ := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='goadmin_users'")
	if exists != nil && len(exists) > 0 {
		log.Println("数据库表已存在，跳过初始化")
		return
	}

	// 优先使用外部 data/init.sql 文件（如果存在）
	var sqlContent string
	externalPath := filepath.Join("data", "init.sql")
	sqlBytes, err := os.ReadFile(externalPath)
	if err != nil {
		// 外部文件不存在，使用嵌入的 SQL
		log.Printf("外部初始化文件不存在 (%v)，使用内置 SQL", err)
		sqlContent = GetEmbeddedInitSQL()
	} else {
		// 使用外部文件
		log.Printf("使用外部初始化文件：%s", externalPath)
		sqlContent = string(sqlBytes)
	}

	// 执行 SQL 脚本
	_, err = conn.Exec(sqlContent)
	if err != nil {
		log.Printf("执行初始化 SQL 失败：%v", err)
		return
	}

	log.Println("数据库表初始化完成")
}

// InitMenu 初始化网络管理主菜单
func InitMenu(conn db.Connection) {
	// 强制清除旧的网络管理菜单数据（以便重新初始化）
	conn.Exec("DELETE FROM goadmin_role_menu WHERE menu_id IN (SELECT id FROM goadmin_menu WHERE title = '网络管理')")
	conn.Exec("DELETE FROM goadmin_menu WHERE parent_id IN (SELECT id FROM goadmin_menu WHERE title = '网络管理')")
	conn.Exec("DELETE FROM goadmin_menu WHERE title = '网络管理'")

	// 插入网络管理主菜单
	networkResult, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (0, 0, '网络管理', '', 'fa fa-sitemap', 1)")
	if err != nil {
		log.Printf("插入网络管理菜单失败：%v", err)
		return
	}
	networkId, _ := networkResult.LastInsertId()
	log.Println("插入网络管理主菜单，图标：fa fa-sitemap")

	// 插入网络管理子菜单
	subMenus := []struct {
		name  string
		uri   string
		icon  string
		order int
	}{
		{"IPv4 路由表", "/network/route-table", "fa fa-list-alt", 1},
		{"IPv4 静态路由", "/network/static-route", "fa fa-road", 2},
		{"Ping 诊断", "/network/ping", "fa fa-exchange", 3},
		{"Traceroute 诊断", "/network/traceroute", "fa fa-random", 4},
		{"虚拟电缆检测", "/network/cable-test", "fa fa-wrench", 5},
	}

	for _, menu := range subMenus {
		// 检查子菜单是否已存在
		exists, _ := conn.Query("SELECT id FROM goadmin_menu WHERE title = ? AND parent_id = ? LIMIT 1", menu.name, networkId)
		if exists != nil && len(exists) > 0 {
			continue // 已存在，跳过
		}
		_, err := conn.Exec("INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`) VALUES (?, 1, ?, ?, ?, ?)",
			networkId, menu.name, menu.uri, menu.icon, menu.order)
		if err != nil {
			log.Printf("插入菜单 %s 失败：%v", menu.name, err)
		} else {
			log.Printf("插入菜单：%s (parent_id=%d)", menu.name, networkId)
		}
	}

	// 添加角色菜单关联
	adminRoleId := int64(1)
	conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, networkId)

	// 添加子菜单关联
	subMenusResult, _ := conn.Query("SELECT id FROM goadmin_menu WHERE parent_id = ?", networkId)
	for _, row := range subMenusResult {
		menuId := int(row["id"].(int64))
		conn.Exec("INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id) VALUES (?, ?)", adminRoleId, menuId)
	}

	log.Println("网络管理菜单初始化完成")
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
	conn.Exec("DELETE FROM goadmin_role_menu WHERE menu_id IN (SELECT id FROM goadmin_menu WHERE title LIKE '%配置%' OR parent_id IN (SELECT id FROM goadmin_menu WHERE title LIKE '%配置%'))")

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
		{"生成树", "/config/stp", "fa fa-tree", 9},
		{"ERPS", "/config/erps", "fa fa-retweet", 10},
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
