-- 添加维护主菜单
INSERT OR IGNORE INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
VALUES (0, 0, '维护', '', 'fa fa-wrench', 3);

-- 获取维护菜单 ID 并插入子菜单（假设维护菜单 ID 为 10）
-- 先查询维护菜单 ID
-- SELECT id FROM goadmin_menu WHERE title = '维护';

-- 插入维护子菜单（使用父 ID=10，如果实际 ID 不同需要调整）
INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '重启/保存', '/maintenance/reboot-save', 'fa fa-power-off', 1
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '重启/保存');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '用户管理', '/maintenance/users', 'fa fa-users', 2
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '用户管理');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '系统配置', '/maintenance/system-config', 'fa fa-cogs', 3
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '系统配置' AND parent_id = (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1));

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '加载配置', '/maintenance/load-config', 'fa fa-upload', 4
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '加载配置');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '文件管理', '/maintenance/files', 'fa fa-file', 5
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '文件管理');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '日志管理', '/maintenance/logs', 'fa fa-history', 6
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '日志管理');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, 'SNMP 配置', '/maintenance/snmp', 'fa fa-bell', 7
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = 'SNMP 配置');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, 'SNMP Trap 配置', '/maintenance/snmp-trap', 'fa fa-exclamation-triangle', 8
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = 'SNMP Trap 配置');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '蠕虫攻击防护', '/maintenance/worm-protection', 'fa fa-bug', 9
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '蠕虫攻击防护');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, 'DDoS 攻击防护', '/maintenance/ddos-protection', 'fa fa-shield', 10
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = 'DDoS 攻击防护');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, 'ARP 攻击防护', '/maintenance/arp-protection', 'fa fa-lock', 11
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = 'ARP 攻击防护');

INSERT INTO goadmin_menu (parent_id, type, title, uri, icon, `order`)
SELECT
    (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1),
    1, '当前会话', '/maintenance/sessions', 'fa fa-clock-o', 12
WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu WHERE title = '当前会话');

-- 添加角色菜单关联（角色 ID=1 是管理员）
INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id)
SELECT 1, (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1);

-- 添加子菜单的角色关联
INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id)
SELECT 1, id FROM goadmin_menu
WHERE parent_id = (SELECT id FROM goadmin_menu WHERE title = '维护' LIMIT 1);
