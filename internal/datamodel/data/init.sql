-- GoAdmin SQLite 初始化脚本
-- 创建所有必要的表

-- 用户表
CREATE TABLE IF NOT EXISTS "goadmin_users" (
    "id" integer PRIMARY KEY autoincrement,
    "username" varchar(255) NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL,
    "name" varchar(255) NOT NULL,
    "remember_token" varchar(255) DEFAULT '',
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 角色表
CREATE TABLE IF NOT EXISTS "goadmin_roles" (
    "id" integer PRIMARY KEY autoincrement,
    "name" varchar(255) NOT NULL UNIQUE,
    "slug" varchar(255) NOT NULL UNIQUE,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 权限表
CREATE TABLE IF NOT EXISTS "goadmin_permissions" (
    "id" integer PRIMARY KEY autoincrement,
    "name" varchar(255) NOT NULL UNIQUE,
    "slug" varchar(255) NOT NULL UNIQUE,
    "http_method" varchar(255) DEFAULT '',
    "http_path" text DEFAULT '',
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 菜单表
CREATE TABLE IF NOT EXISTS "goadmin_menu" (
    "id" integer PRIMARY KEY autoincrement,
    "parent_id" integer NOT NULL DEFAULT 0,
    "type" integer NOT NULL DEFAULT 0,
    "order" integer NOT NULL DEFAULT 0,
    "title" varchar(255) NOT NULL,
    "icon" varchar(255) NOT NULL,
    "uri" varchar(255) NOT NULL DEFAULT '',
    "plugin_name" varchar(255) NOT NULL DEFAULT '',
    "header" varchar(255) DEFAULT NULL,
    "uuid" varchar(255) DEFAULT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 角色用户关联表
CREATE TABLE IF NOT EXISTS "goadmin_role_users" (
    "id" integer PRIMARY KEY autoincrement,
    "role_id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS "goadmin_role_permissions" (
    "id" integer PRIMARY KEY autoincrement,
    "role_id" integer NOT NULL,
    "permission_id" integer NOT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 用户权限关联表
CREATE TABLE IF NOT EXISTS "goadmin_user_permissions" (
    "id" integer PRIMARY KEY autoincrement,
    "user_id" integer NOT NULL,
    "permission_id" integer NOT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS "goadmin_role_menu" (
    "id" integer PRIMARY KEY autoincrement,
    "role_id" integer NOT NULL,
    "menu_id" integer NOT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 操作日志表
CREATE TABLE IF NOT EXISTS "goadmin_operation_log" (
    "id" integer PRIMARY KEY autoincrement,
    "user_id" integer NOT NULL,
    "path" varchar(255) NOT NULL,
    "method" varchar(10) NOT NULL,
    "ip" varchar(255) NOT NULL,
    "input" text NOT NULL,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 站点配置表
CREATE TABLE IF NOT EXISTS "goadmin_site" (
    "id" integer PRIMARY KEY autoincrement,
    "key" varchar(255) DEFAULT NULL,
    "value" text DEFAULT NULL,
    "description" varchar(255) DEFAULT NULL,
    "state" integer NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- Session 表（用于 CSRF Token）
CREATE TABLE IF NOT EXISTS "goadmin_session" (
    "id" integer PRIMARY KEY autoincrement,
    "sid" varchar(255) NOT NULL DEFAULT '',
    "values" varchar(255) NOT NULL DEFAULT '',
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP default CURRENT_TIMESTAMP
);

-- 初始化默认数据

-- 默认管理员用户（密码：admin）
INSERT OR IGNORE INTO "goadmin_users" ("id", "username", "password", "name", "remember_token")
VALUES (1, 'admin', '$2a$10$TEDU/aUxLkr2wCxGxI62/.yOtzrzfv426DLLdyha9H2GpWRggB0di', 'Administrator', '');

-- 默认角色
INSERT OR IGNORE INTO "goadmin_roles" ("id", "name", "slug") VALUES (1, 'Administrator', 'administrator');

-- 默认权限
INSERT OR IGNORE INTO "goadmin_permissions" ("id", "name", "slug", "http_method", "http_path")
VALUES (1, 'All permission', '*', '', '*');

-- 角色用户关联
INSERT OR IGNORE INTO "goadmin_role_users" ("role_id", "user_id") VALUES (1, 1);

-- 角色权限关联
INSERT OR IGNORE INTO "goadmin_role_permissions" ("role_id", "permission_id") VALUES (1, 1);

-- Dashboard 菜单（type=1 表示菜单）
INSERT OR IGNORE INTO "goadmin_menu" ("id", "parent_id", "type", "order", "title", "icon", "uri")
VALUES (1, 0, 1, 0, 'Dashboard', 'fa fa-dashboard', '/dashboard');

-- Admin 管理菜单（type=0 表示目录，只在侧边栏显示）
INSERT OR IGNORE INTO "goadmin_menu" ("id", "parent_id", "type", "order", "title", "icon", "uri")
VALUES
(2, 0, 0, 2, 'Admin', 'fa-tasks', ''),
(3, 2, 1, 2, '用户管理', 'fa-users', '/info/manager'),
(4, 2, 1, 3, '角色管理', 'fa-user', '/info/roles'),
(5, 2, 1, 4, '权限管理', 'fa-ban', '/info/permission'),
(6, 2, 1, 5, '菜单管理', 'fa-bars', '/menu'),
(7, 2, 1, 6, '操作日志', 'fa-history', '/info/op');

-- 角色菜单关联
INSERT OR IGNORE INTO "goadmin_role_menu" ("role_id", "menu_id") VALUES
(1, 1),  -- Dashboard
(1, 2),  -- Admin
(1, 3),  -- 用户管理
(1, 4),  -- 角色管理
(1, 5),  -- 权限管理
(1, 6),  -- 菜单管理
(1, 7);  -- 操作日志
