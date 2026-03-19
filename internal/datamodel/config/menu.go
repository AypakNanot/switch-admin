package config

import (
	"log"

	"github.com/GoAdminGroup/go-admin/modules/db"
)

// InitMenu 初始化配置模块菜单
func InitMenu(conn db.Connection) {
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
