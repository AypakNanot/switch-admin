# 双模式适配架构变更提案

**变更 ID**: `dual-mode-adapter-architecture`
**创建日期**: 2026-03-15
**状态**: `draft`
**作者**: AI Assistant

---

## 1. 变更概述

### 1.1 变更目标

为 switch-admin 项目引入**双模式适配架构**，支持：
1. **离线测试模式** - 在虚拟机/本地开发环境运行，使用数据库模拟数据
2. **交换机模式** - 在真实交换机上运行，直接操作交换机硬件

### 1.2 变更原因

| 场景 | 当前问题 | 变更后收益 |
|------|----------|------------|
| **开发测试** | 需要在真实交换机上测试，设备资源受限 | 可在任意环境离线测试，不依赖硬件 |
| **功能验证** | 无法快速验证功能正确性 | 静态模拟数据，快速验证业务逻辑 |
| **多适配器** | 硬编码单一 CLI 方式 | 支持 CLI/Netconf/REST API，可配置切换 |
| **运行时切换** | 需要重启服务切换模式 | API 动态切换，无需重启 |

### 1.3 影响范围

| 模块 | 影响程度 | 说明 |
|------|----------|------|
| `internal/service/` | 高 | 新增模式判断层、适配器层 |
| `internal/datamodel/` | 中 | 部分页面需适配双模式 |
| `cmd/main.go` | 中 | 启动时初始化模式配置 |
| 数据库 Schema | 中 | 新增模式配置表、适配器配置表 |

---

## 2. 架构设计

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           HTTP 请求                                      │
│                        (GoAdmin Handler)                                 │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Service 层                                      │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                    【模式判断层】ModeResolver                    │    │
│  │  • 读取当前运行模式 (离线测试/交换机)                             │    │
│  │  • 根据模式路由到不同实现                                        │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
                    ▼                           ▼
    ┌───────────────────────────┐   ┌───────────────────────────┐
    │   离线测试模式             │   │   交换机模式               │
    │   (Mock Mode)             │   │   (Switch Mode)           │
    │                           │   │                           │
    │  ┌─────────────────────┐  │   │  ┌─────────────────────┐  │
    │  │   数据模拟层          │  │   │  │   交换机适配层       │  │
    │  │   MockDataProvider  │  │   │  │   SwitchAdapter     │  │
    │  └─────────────────────┘  │   │  └─────────────────────┘  │
    │           │                │   │           │                │
    │           ▼                │   │           ▼                │
    │  ┌─────────────────────┐  │   │  ┌─────────────────────┐  │
    │  │   SQLite 数据库      │  │   │  │   适配器实现         │  │
    │  │   (模拟数据)        │  │   │  │  ┌────────────────┐ │  │
    │  └─────────────────────┘  │   │  │  │ CLI Adapter    │ │  │
    │                           │   │  │  ├────────────────┤ │  │
    │                           │   │  │  │ Netconf Adapter│ │  │
    │                           │   │  │  ├────────────────┤ │  │
    │                           │   │  │  │ REST Adapter   │ │  │
    │                           │   │  │  └────────────────┘ │  │
    │                           │   │  └─────────────────────┘  │
    │                           │   │           │                │
    │                           │   │           ▼                │
    │                           │   │  ┌─────────────────────┐  │
    │                           │   │  │   交换机硬件         │  │
    │                           │   │  └─────────────────────┘  │
    └───────────────────────────┘   └───────────────────────────┘
```

### 2.2 核心接口定义

#### 2.2.1 模式提供者接口

```go
// DataProvider 数据提供者接口
type DataProvider interface {
    // GetPortStatus 获取端口状态
    GetPortStatus(portID string) (*PortStatus, error)
    // GetAllPorts 获取所有端口
    GetAllPorts() ([]*PortStatus, error)
    // SetPortAdminStatus 设置端口管理状态
    SetPortAdminStatus(portID string, enabled bool) error
    // GetSystemInfo 获取系统信息
    GetSystemInfo() (*SystemInfo, error)
    // ClearPortStats 清零端口统计
    ClearPortStats(portID string) error
    // ClearAllPortStats 清零所有端口统计
    ClearAllPortStats() error
}
```

#### 2.2.2 交换机适配器接口

```go
// SwitchAdapter 交换机适配器接口
type SwitchAdapter interface {
    // Name 返回适配器名称
    Name() string
    // Connect 连接到交换机
    Connect() error
    // Close 关闭连接
    Close() error
    // IsConnected 检查连接状态
    IsConnected() bool

    // 端口操作
    GetPortStatus(portID string) (*PortStatus, error)
    GetAllPorts() ([]*PortStatus, error)
    SetPortAdminStatus(portID string, enabled bool) error

    // 系统操作
    GetSystemInfo() (*SystemInfo, error)

    // 统计操作
    ClearPortStats(portID string) error
    ClearAllPortStats() error
}
```

### 2.3 运行模式设计

| 模式 | 标识 | 数据源 | 用户模块 | 其他功能 |
|------|------|--------|----------|----------|
| **离线测试模式** | `mock` | 数据库模拟数据 | 数据库 | 数据库 |
| **交换机模式** | `switch` | 交换机硬件 | 数据库 + 交换机 | 交换机 |

### 2.4 适配器配置设计

```go
// AdapterConfig 适配器配置
type AdapterConfig struct {
    ID           int64  `json:"id" xorm:"pk autoincr"`
    FunctionName string `json:"function_name"`  // 功能名称：port, vlan, lacp...
    AdapterType  string `json:"adapter_type"`   // 适配器类型：cli, netconf, rest
    Priority     int    `json:"priority"`       // 优先级（同功能多适配器时）
    Enabled      bool   `json:"enabled"`        // 是否启用
    Config       string `json:"config"`         // JSON 格式配置参数
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

---

## 3. 数据库设计

### 3.1 新增表结构

#### 3.1.1 系统配置表 `sys_config`

```sql
CREATE TABLE sys_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key VARCHAR(64) NOT NULL UNIQUE,
    config_value TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认运行模式
INSERT INTO sys_config (config_key, config_value, description)
VALUES ('run_mode', 'mock', '运行模式：mock=离线测试，switch=交换机');
```

#### 3.1.2 适配器配置表 `adapter_config`

```sql
CREATE TABLE adapter_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    function_name VARCHAR(64) NOT NULL,
    adapter_type VARCHAR(32) NOT NULL,
    priority INTEGER DEFAULT 0,
    enabled INTEGER DEFAULT 1,
    config TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(function_name, adapter_type)
);

-- 插入默认适配器配置
INSERT INTO adapter_config (function_name, adapter_type, priority, enabled, config)
VALUES
    ('port', 'cli', 1, 1, '{"protocol":"ssh","host":"localhost","port":22}'),
    ('system', 'cli', 1, 1, '{"protocol":"ssh","host":"localhost","port":22}');
```

#### 3.1.3 端口模拟数据表 `mock_port`

```sql
CREATE TABLE mock_port (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    port_name VARCHAR(32) NOT NULL UNIQUE,
    admin_status INTEGER DEFAULT 1,      -- 1=enable, 0=disable
    link_status INTEGER DEFAULT 0,       -- 1=up, 0=down
    speed VARCHAR(16) DEFAULT '-',
    duplex VARCHAR(16) DEFAULT '-',
    description VARCHAR(256) DEFAULT '',
    rx_bytes BIGINT DEFAULT 0,
    tx_bytes BIGINT DEFAULT 0,
    rx_packets BIGINT DEFAULT 0,
    tx_packets BIGINT DEFAULT 0,
    rx_errors BIGINT DEFAULT 0,
    tx_errors BIGINT DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入模拟数据（48 端口交换机示例）
INSERT INTO mock_port (port_name, admin_status, link_status, speed, duplex, description)
VALUES
    ('GE1/0/1', 1, 1, '1000M', 'Full', 'Server-A'),
    ('GE1/0/2', 1, 0, '-', '-', ''),
    ('GE1/0/3', 1, 1, '1000M', 'Full', 'AP-Floor1'),
    ('GE1/0/4', 0, 0, '-', '-', 'Unused');
-- ... 更多端口
```

#### 3.1.4 系统信息模拟数据表 `mock_system_info`

```sql
CREATE TABLE mock_system_info (
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

-- 插入默认模拟数据
INSERT INTO mock_system_info (boot_time, uptime_seconds)
VALUES (datetime('now', '-18 hours'), 64920);
```

### 3.2 用户模块数据同步

```go
// SyncUsersFromSwitch 在交换机模式下，启动时同步用户到数据库
func SyncUsersFromSwitch(ctx context.Context) error {
    // 1. 从交换机读取用户列表
    users, err := switchAdapter.GetUsers()
    if err != nil {
        return err
    }

    // 2. 持久化到数据库
    for _, user := range users {
        err = db.InsertOrUpdate(&User{
            Username: user.Username,
            Role: user.Role,
            // ...
        })
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

## 4. 运行时切换 API 设计

### 4.1 切换运行模式

```
POST /api/v1/system/mode
Content-Type: application/json

{
    "mode": "mock"  // 或 "switch"
}
```

**响应**:
```json
{
    "code": 200,
    "message": "模式切换成功",
    "data": {
        "previous_mode": "switch",
        "current_mode": "mock",
        "requires_restart": false
    }
}
```

### 4.2 获取当前模式

```
GET /api/v1/system/mode
```

**响应**:
```json
{
    "code": 200,
    "data": {
        "mode": "mock",
        "mode_description": "离线测试模式",
        "adapter_summary": {
            "port": "cli",
            "system": "cli",
            "vlan": "netconf"
        }
    }
}
```

### 4.3 配置适配器

```
PUT /api/v1/adapters/:function/config
Content-Type: application/json

{
    "adapter_type": "netconf",
    "enabled": true,
    "priority": 1,
    "config": {
        "host": "192.168.1.1",
        "port": 830,
        "username": "admin",
        "password": "***"
    }
}
```

---

## 5. 实现任务清单

### 5.1 核心层实现

- [ ] 创建 `internal/service/mode/mode_resolver.go` - 模式解析器
- [ ] 创建 `internal/service/provider/data_provider.go` - 数据提供者接口
- [ ] 创建 `internal/service/provider/mock_provider.go` - 模拟数据实现
- [ ] 创建 `internal/service/adapter/switch_adapter.go` - 交换机适配器接口
- [ ] 创建 `internal/service/adapter/cli_adapter.go` - CLI 适配器实现
- [ ] 创建 `internal/service/adapter/netconf_adapter.go` - Netconf 适配器实现
- [ ] 创建 `internal/service/adapter/rest_adapter.go` - REST API 适配器实现
- [ ] 创建 `internal/service/adapter/adapter_factory.go` - 适配器工厂

### 5.2 数据库层实现

- [ ] 创建数据库迁移脚本 `scripts/migrate_v2.sql`
- [ ] 创建模型 `internal/model/sys_config.go`
- [ ] 创建模型 `internal/model/adapter_config.go`
- [ ] 创建模型 `internal/model/mock_port.go`
- [ ] 创建模型 `internal/model/mock_system_info.go`
- [ ] 创建 DAO `internal/dao/config_dao.go`
- [ ] 创建 DAO `internal/dao/mock_data_dao.go`

### 5.3 API 层实现

- [ ] 创建 Handler `internal/handler/mode_handler.go`
- [ ] 创建 Handler `internal/handler/adapter_handler.go`
- [ ] 在 `cmd/main.go` 中注册新路由
- [ ] 在 `cmd/main.go` 中添加启动时模式初始化

### 5.4 现有代码改造

- [ ] 改造 `internal/datamodel/dashboard.go` 使用 DataProvider
- [ ] 改造 `internal/datamodel/user.go` 使用 DataProvider
- [ ] 更新所有端口相关功能使用 DataProvider

### 5.5 测试与验证

- [ ] 编写单元测试 `service/mode/mode_resolver_test.go`
- [ ] 编写适配器测试 `service/adapter/cli_adapter_test.go`
- [ ] 编写集成测试 `tests/dual_mode_test.go`
- [ ] 手动验证离线测试模式功能
- [ ] 手动验证交换机模式功能

---

## 6. 验收标准

### 6.1 功能验收

| 用例 | 预期结果 |
|------|----------|
| 启动服务，默认为 mock 模式 | 服务正常启动，Dashboard 显示模拟数据 |
| API 切换到 switch 模式 | 切换成功，无需重启服务 |
| mock 模式下配置端口 | 配置写入数据库，页面显示更新 |
| switch 模式下配置端口 | 配置下发到交换机，返回实际状态 |
| 用户模块在 switch 模式 | 启动时同步交换机用户到数据库 |
| 配置适配器优先级 | 高优先级适配器优先使用 |

### 6.2 性能验收

| 指标 | 目标值 |
|------|--------|
| 模式切换响应时间 | ≤ 500ms |
| mock 模式 API 响应 | ≤ 100ms |
| switch 模式 CLI 响应 | ≤ 2s |
| switch 模式 Netconf 响应 | ≤ 1s |

### 6.3 质量验收

| 指标 | 目标值 |
|------|--------|
| 单元测试覆盖率 | ≥ 70% |
| 关键路径测试 | 100% 覆盖 |
| 代码审查 | 无严重问题 |

---

## 7. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 现有代码改造工作量大 | 高 | 分阶段实施，先新增后替换 |
| 适配器实现复杂度高 | 中 | 优先实现 CLI，其他逐步添加 |
| 模式切换并发问题 | 中 | 使用读写锁保护模式状态 |
| 模拟数据与实际数据不一致 | 低 | 定期对比验证，更新模拟数据 |

---

## 8. 里程碑

| 阶段 | 内容 | 预计时间 |
|------|------|----------|
| **Phase 1** | 核心接口设计 + 数据库迁移 | 1 周 |
| **Phase 2** | Mock 数据提供者实现 | 1 周 |
| **Phase 3** | CLI 适配器实现 | 1 周 |
| **Phase 4** | 模式切换 API + 管理界面 | 1 周 |
| **Phase 5** | 现有代码改造 + 测试 | 2 周 |
| **Phase 6** | Netconf/REST 适配器 | 后续迭代 |

---

## 9. 参考文档

- [PRD/switch-admin 项目需求规格说明书.md](../PRD/switch-admin 项目需求规格说明书.md)
- [PRD/01-需求总纲与架构规范.md](../PRD/01-需求总纲与架构规范.md)
- [openspec/project.md](./project.md)

---

**审批状态**: `pending`
**审批人**: `-`
**最后更新**: 2026-03-15
