# 实现任务清单

**变更 ID**: `dual-mode-adapter-architecture`
**状态**: `pending`

---

## Phase 1: 核心接口与数据库迁移

### 1.1 创建核心接口

- [ ] `internal/service/provider/data_provider.go` - 数据提供者接口
  - [ ] 定义 `DataProvider` 接口
  - [ ] 定义 `PortStatus` 结构体
  - [ ] 定义 `SystemInfo` 结构体

- [ ] `internal/service/adapter/switch_adapter.go` - 交换机适配器接口
  - [ ] 定义 `SwitchAdapter` 接口
  - [ ] 定义连接管理接口
  - [ ] 定义错误类型

- [ ] `internal/service/mode/mode_resolver.go` - 模式解析器
  - [ ] 定义 `ModeResolver` 结构
  - [ ] 实现 `GetCurrentMode()` 方法
  - [ ] 实现 `SwitchMode()` 方法
  - [ ] 实现模式切换锁机制

### 1.2 数据库迁移

- [ ] `scripts/migrate_v2.sql` - 数据库迁移脚本
  - [ ] 创建 `sys_config` 表
  - [ ] 创建 `adapter_config` 表
  - [ ] 创建 `mock_port` 表
  - [ ] 创建 `mock_system_info` 表
  - [ ] 插入默认配置数据
  - [ ] 插入模拟测试数据

- [ ] `internal/model/sys_config.go` - 系统配置模型
- [ ] `internal/model/adapter_config.go` - 适配器配置模型
- [ ] `internal/model/mock_port.go` - 端口模拟数据模型
- [ ] `internal/model/mock_system_info.go` - 系统信息模拟模型

### 1.3 DAO 层

- [ ] `internal/dao/config_dao.go` - 配置数据访问
  - [ ] `GetRunMode()` - 获取运行模式
  - [ ] `SetRunMode(mode string)` - 设置运行模式
  - [ ] `GetAdapterConfig(functionName string)` - 获取适配器配置

- [ ] `internal/dao/mock_data_dao.go` - 模拟数据访问
  - [ ] `GetAllMockPorts()` - 获取所有模拟端口
  - [ ] `UpdateMockPort()` - 更新模拟端口
  - [ ] `GetMockSystemInfo()` - 获取模拟系统信息

---

## Phase 2: Mock 数据提供者实现

### 2.1 核心实现

- [ ] `internal/service/provider/mock_provider.go`
  - [ ] 实现 `GetPortStatus()` - 从数据库读取端口状态
  - [ ] 实现 `GetAllPorts()` - 从数据库读取所有端口
  - [ ] 实现 `SetPortAdminStatus()` - 更新数据库端口状态
  - [ ] 实现 `GetSystemInfo()` - 从数据库读取系统信息
  - [ ] 实现 `ClearPortStats()` - 清零单个端口统计
  - [ ] 实现 `ClearAllPortStats()` - 清零所有端口统计

### 2.2 模拟数据增强

- [ ] `internal/service/provider/mock_provider_stats.go`
  - [ ] 实现流量统计自动增长（模拟真实流量）
  - [ ] 实现可配置的更新频率
  - [ ] 实现统计重置功能

---

## Phase 3: CLI 适配器实现

### 3.1 CLI 适配器核心

- [ ] `internal/service/adapter/cli_adapter.go`
  - [ ] 实现 SSH 连接管理
  - [ ] 实现命令发送
  - [ ] 实现响应解析

- [ ] `internal/service/adapter/cli_parser.go`
  - [ ] 实现端口状态响应解析
  - [ ] 实现系统信息响应解析
  - [ ] 实现错误码解析

### 3.2 CLI 命令集

- [ ] `internal/service/adapter/cli_commands.go`
  - [ ] `GetPortStatusCommand(portID string)` - 获取端口状态命令
  - [ ] `SetPortAdminCommand(portID string, enabled bool)` - 设置端口状态命令
  - [ ] `GetSystemInfoCommand()` - 获取系统信息命令
  - [ ] `ClearPortStatsCommand(portID string)` - 清零端口统计命令

---

## Phase 4: 适配器框架

### 4.1 适配器注册

- [ ] `internal/service/adapter/adapter_registry.go`
  - [ ] 实现适配器注册表
  - [ ] 实现按功能获取适配器
  - [ ] 实现适配器优先级管理

### 4.2 适配器工厂

- [ ] `internal/service/adapter/adapter_factory.go`
  - [ ] 实现 `CreateAdapter(type string, config map[string]interface{})`
  - [ ] 支持 `cli` 类型
  - [ ] 支持 `netconf` 类型（预留）
  - [ ] 支持 `rest` 类型（预留）

### 4.3 Netconf/REST 预留

- [ ] `internal/service/adapter/netconf_adapter.go` - 框架（待实现）
- [ ] `internal/service/adapter/rest_adapter.go` - 框架（待实现）

---

## Phase 5: API 层实现

### 5.1 Handler 实现

- [ ] `internal/handler/mode_handler.go`
  - [ ] `GetMode()` - 获取当前模式 API
  - [ ] `SwitchMode()` - 切换模式 API

- [ ] `internal/handler/adapter_handler.go`
  - [ ] `ListAdapters()` - 列出所有适配器配置
  - [ ] `GetAdapterConfig(functionName string)` - 获取适配器配置
  - [ ] `UpdateAdapterConfig()` - 更新适配器配置

### 5.2 路由注册

- [ ] `cmd/main.go`
  - [ ] 注册 `/api/v1/system/mode` 路由
  - [ ] 注册 `/api/v1/adapters` 路由
  - [ ] 启动时初始化模式解析器
  - [ ] 启动时初始化适配器工厂

### 5.3 启动初始化

- [ ] `internal/bootstrap.go`
  - [ ] 读取数据库运行模式
  - [ ] 初始化模式解析器
  - [ ] 加载适配器配置
  - [ ] 交换机模式下同步用户到数据库

---

## Phase 6: 现有代码改造

### 6.1 Dashboard 改造

- [ ] `internal/datamodel/dashboard.go`
  - [ ] 使用 `DataProvider` 替代直接数据
  - [ ] 根据模式自动选择数据源
  - [ ] 移除硬编码数据

### 6.2 用户管理改造

- [ ] `internal/datamodel/user.go`
  - [ ] 启动时同步交换机用户
  - [ ] 支持从数据库读取用户
  - [ ] 支持写入到交换机/数据库

### 6.3 其他页面改造

- [ ] `internal/datamodel/generators.go`
- [ ] `internal/datamodel/bootstrap.go`

---

## Phase 7: 测试与验证

### 7.1 单元测试

- [ ] `internal/service/mode/mode_resolver_test.go`
- [ ] `internal/service/provider/mock_provider_test.go`
- [ ] `internal/service/adapter/cli_adapter_test.go`
- [ ] `internal/dao/config_dao_test.go`
- [ ] `internal/dao/mock_data_dao_test.go`

### 7.2 集成测试

- [ ] `tests/dual_mode_integration_test.go`
  - [ ] 测试模式切换流程
  - [ ] 测试 mock 模式数据读写
  - [ ] 测试适配器配置更新

### 7.3 手动测试

- [ ] 验证 mock 模式 Dashboard 正常显示
- [ ] 验证 mock 模式端口配置生效
- [ ] 验证模式切换无需重启
- [ ] 验证适配器配置持久化

---

## Phase 8: 文档与部署

### 8.1 文档更新

- [ ] 更新 `README.md` - 添加双模式架构说明
- [ ] 创建 `docs/DUAL_MODE_GUIDE.md` - 双模式使用指南
- [ ] 创建 `docs/ADAPTER_DEVELOPMENT.md` - 适配器开发指南

### 8.2 部署脚本

- [ ] `scripts/deploy_mock_data.sh` - 部署模拟数据
- [ ] `scripts/switch_mode.sh` - 快速切换模式脚本

---

## 任务依赖关系

```
Phase 1 (核心接口)
    ↓
Phase 2 (Mock Provider)
    ↓
Phase 3 (CLI Adapter)
    ↓
Phase 4 (Adapter Framework)
    ↓
Phase 5 (API Layer)
    ↓
Phase 6 (代码改造)
    ↓
Phase 7 (测试验证)
    ↓
Phase 8 (文档部署)
```

---

## 进度追踪

| 阶段 | 状态 | 完成日期 |
|------|------|----------|
| Phase 1 | pending | - |
| Phase 2 | pending | - |
| Phase 3 | pending | - |
| Phase 4 | pending | - |
| Phase 5 | pending | - |
| Phase 6 | pending | - |
| Phase 7 | pending | - |
| Phase 8 | pending | - |
