-- =====================================================
-- switch-admin 数据库迁移脚本 v2
-- 双模式适配架构支持
-- =====================================================

-- 1. 系统配置表
-- 存储运行模式、适配器配置等系统级配置
CREATE TABLE IF NOT EXISTS sys_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key VARCHAR(64) NOT NULL UNIQUE,
    config_value TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认运行模式
INSERT OR IGNORE INTO sys_config (config_key, config_value, description)
VALUES ('run_mode', 'mock', '运行模式：mock=离线测试，switch=交换机');

-- 插入其他默认配置
INSERT OR IGNORE INTO sys_config (config_key, config_value, description)
VALUES
    ('app_name', 'switch-admin', '应用名称'),
    ('app_version', '1.0.0', '应用版本');


-- =====================================================
-- 2. 适配器配置表
-- 支持一功能多适配器，通过优先级控制选择
-- =====================================================
CREATE TABLE IF NOT EXISTS adapter_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    function_name VARCHAR(64) NOT NULL,      -- 功能名称：port, vlan, system, lacp...
    adapter_type VARCHAR(32) NOT NULL,       -- 适配器类型：cli, netconf, rest
    priority INTEGER DEFAULT 0,              -- 优先级：数字越大优先级越高
    enabled INTEGER DEFAULT 1,               -- 是否启用 (1=启用，0=禁用)
    config TEXT,                             -- JSON 格式配置参数
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(function_name, adapter_type)
);

-- 插入默认适配器配置（CLI 方式）
INSERT OR IGNORE INTO adapter_config (function_name, adapter_type, priority, enabled, config)
VALUES
    ('port', 'cli', 1, 1, '{"protocol":"ssh","host":"127.0.0.1","port":22,"timeout":30}'),
    ('system', 'cli', 1, 1, '{"protocol":"ssh","host":"127.0.0.1","port":22,"timeout":30}'),
    ('vlan', 'cli', 1, 1, '{"protocol":"ssh","host":"127.0.0.1","port":22,"timeout":30}'),
    ('lacp', 'cli', 1, 1, '{"protocol":"ssh","host":"127.0.0.1","port":22,"timeout":30}');


-- =====================================================
-- 3. 端口模拟数据表
-- 离线测试模式下使用
-- =====================================================
CREATE TABLE IF NOT EXISTS mock_port (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    port_name VARCHAR(32) NOT NULL UNIQUE,
    admin_status INTEGER DEFAULT 1,          -- 管理状态：1=enable, 0=disable
    link_status INTEGER DEFAULT 0,           -- 链路状态：1=up, 0=down
    speed VARCHAR(16) DEFAULT '-',           -- 速率：10M/100M/1000M/10G/-
    duplex VARCHAR(16) DEFAULT '-',          -- 双工：Full/Half/-
    description VARCHAR(256) DEFAULT '',     -- 端口描述
    rx_bytes BIGINT DEFAULT 0,               -- 接收字节数
    tx_bytes BIGINT DEFAULT 0,               -- 发送字节数
    rx_packets BIGINT DEFAULT 0,             -- 接收包数
    tx_packets BIGINT DEFAULT 0,             -- 发送包数
    rx_errors BIGINT DEFAULT 0,              -- 接收错误数
    tx_errors BIGINT DEFAULT 0,              -- 发送错误数
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入模拟数据（48 端口交换机示例）
INSERT OR IGNORE INTO mock_port (port_name, admin_status, link_status, speed, duplex, description)
VALUES
    ('GE1/0/1', 1, 1, '1000M', 'Full', 'Server-A'),
    ('GE1/0/2', 1, 0, '-', '-', ''),
    ('GE1/0/3', 1, 1, '1000M', 'Full', 'AP-Floor1'),
    ('GE1/0/4', 0, 0, '-', '-', 'Unused'),
    ('GE1/0/5', 1, 1, '1000M', 'Full', 'AP-Floor2'),
    ('GE1/0/6', 1, 1, '1000M', 'Full', 'Switch-Uplink'),
    ('GE1/0/7', 1, 0, '-', '-', ''),
    ('GE1/0/8', 1, 1, '1000M', 'Full', 'Router-WAN');


-- =====================================================
-- 4. 系统信息模拟数据表
-- 离线测试模式下使用
-- =====================================================
CREATE TABLE IF NOT EXISTS mock_system_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model VARCHAR(64) DEFAULT 'BroadEdge-S3652',
    serial_number VARCHAR(64) DEFAULT 'E605MT252088',
    mac_address VARCHAR(32) DEFAULT '00:07:30:D2:35:67',
    software_version VARCHAR(64) DEFAULT 'OPTEL v7.0.5.15',
    hardware_version VARCHAR(32) DEFAULT '3.0',
    uptime_seconds INTEGER DEFAULT 0,
    boot_time DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认模拟数据（运行时间从启动时计算）
INSERT OR IGNORE INTO mock_system_info (id, model, serial_number, mac_address, software_version, hardware_version, boot_time)
VALUES (1, 'BroadEdge-S3652', 'E605MT252088', '00:07:30:D2:35:67', 'OPTEL v7.0.5.15', '3.0', datetime('now'));


-- =====================================================
-- 5. 适配器连接状态表（可选，用于调试和监控）
-- =====================================================
CREATE TABLE IF NOT EXISTS adapter_connection_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    adapter_type VARCHAR(32) NOT NULL,
    function_name VARCHAR(64) NOT NULL,
    operation VARCHAR(64) NOT NULL,        -- connect, disconnect, command
    status VARCHAR(32) NOT NULL,           -- success, failed
    message TEXT,
    duration_ms INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- =====================================================
-- 6. 创建索引（优化查询性能）
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_adapter_config_function ON adapter_config(function_name);
CREATE INDEX IF NOT EXISTS idx_adapter_config_enabled ON adapter_config(enabled);
CREATE INDEX IF NOT EXISTS idx_sys_config_key ON sys_config(config_key);

-- =====================================================
-- 迁移完成提示
-- =====================================================
SELECT 'Migration v2 completed successfully!' AS status;
