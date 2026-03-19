package maintenance

import (
	"log"

	"github.com/GoAdminGroup/go-admin/modules/db"
)

// InitMenu 初始化维护模块菜单
func InitMenu(conn db.Connection) {
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
