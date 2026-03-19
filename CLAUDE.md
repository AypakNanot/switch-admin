<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

---

# 项目架构总览

## 核心架构模式

本项目采用 **五层分层架构 + 策略模式**，核心设计原则：

### 1. 五层架构

```
Handler 层 (internal/handler/*)
    ↓ HTTP 请求处理 (Gin)
Service 层 (internal/service/*_service.go)
    ↓ 业务逻辑 + 模式切换
Provider 接口层 (internal/service/provider/*)
    ↓ 接口实现 (CLI/Mock/REST)
CLI/Mock Provider (internal/service/provider/{config|network|maintenance|diagnostic}/*)
    ↓ 数据源 (真实交换机 CLI / 数据库模拟)
DataModel 层 (internal/datamodel/*)
    ↓ UI 渲染 (GoAdmin Panel)
```

### 2. 项目统计

| 统计项 | 数量 |
|--------|------|
| 总代码行数 | ~20,683 行 Go |
| Handler 文件 | 37+ |
| DataModel 文件 | 47+ |
| Service 文件 | 26+ |
| Provider 文件 | 23+ |
| Model 文件 | 11 |

---

# 核心模块实现方法论

## 模块总览

| 模块 | 核心功能 | Handler | Service | Provider | DataModel |
|------|----------|---------|---------|----------|-----------|
| **System** | 系统配置、Dashboard、用户管理 | `system/` | - | - | `system/` |
| **Diagnostic** | Ping、Traceroute、电缆检测 | `diagnostic/` | `DiagnosticService` | `diagnostic/` | `diagnostic/` |
| **Maintenance** | 系统维护、安全防护、SNMP | `maintenance/` | `MaintenanceService` | `maintenance/` | `maintenance/` |
| **Network** | VLAN、端口、链路聚合、STP、ACL | `network/` | `NetworkService` | `network/` | `network/` |
| **Config** | 端口配置、风暴控制、ERPS、PoE 等 14 项 | `config/` | `ConfigService` | `config/` | `config/` |

---

# Diagnostic（诊断）模块

## 架构实现

### Handler 层 (`internal/handler/diagnostic/`)
```go
type DiagnosticHandler struct {
    service *service.DiagnosticService
}

func NewDiagnosticHandler() *DiagnosticHandler {
    return &DiagnosticHandler{service: service.GetDiagnosticService()}
}

// 路由对应：
// POST /api/v1/diagnostic/ping -> CreatePingTask()
// GET  /api/v1/diagnostic/ping/:task_id -> GetPingTaskResult()
// DELETE /api/v1/diagnostic/ping/:task_id -> DeletePingTask()
// POST /api/v1/diagnostic/traceroute -> CreateTracerouteTask()
// GET  /api/v1/diagnostic/cable/ports -> GetDetectablePorts()
// POST /api/v1/diagnostic/cable -> ExecuteCableTest()
```

### Service 层 (`internal/service/diagnostic_service.go`)
```go
// 单例模式 + 异步任务管理
type DiagnosticService struct {
    pingTasks    map[string]*model.PingTask
    pingResults  map[string]*model.PingTaskResponse
    traceTasks   map[string]*model.TracerouteTask
    traceResults map[string]*model.TracerouteResponse
    cableTasks   map[string]string
    cableResults map[string]*model.CableTestResult
    mu           sync.RWMutex
    modeResolver *mode.ModeResolver
}

func GetDiagnosticService() *DiagnosticService {
    diagnosticOnce.Do(func() {
        diagnosticService = &DiagnosticService{...}
        go diagnosticService.cleanupRoutine() // 定期清理超时任务
    })
    return diagnosticService
}
```

### Provider 接口层
```go
// interface.go
type DiagnosticProvider interface {
    ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error)
    ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error)
    ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error)
}
```

### DataModel 层
| 页面 | 文件 | 路由 |
|------|------|------|
| Ping 诊断 | `ping.go` | `/admin/network/ping` |
| Traceroute | `traceroute.go` | `/admin/network/traceroute` |
| 电缆检测 | `cable.go` | `/admin/network/cable-test` |

---

# Maintenance（维护）模块

## 文件组织规范

**按功能子模块拆分** - 每个子功能独立成文件：

| 功能模块 | Handler | Provider | DataModel |
|---------|---------|----------|-----------|
| 重启/保存 | - | `reboot_save.go` | `reboot_save.go` |
| 用户管理 | `user_management.go` | `user_management.go` | `user_management.go` |
| 会话管理 | `session_management.go` | `session_management.go` | `session_management.go` |
| 日志管理 | `log_management.go` | `log_management.go` | `log_management.go` |
| 文件管理 | `file_management.go` | `file_management.go` | `file_management.go` |
| SNMP 配置 | `snmp_config.go` | `snmp_config.go` | `snmp_config.go` |
| SNMP Trap | `snmp_trap.go` | `snmp_trap.go` | `snmp_trap.go` |
| 蠕虫防护 | `worm_protection.go` | `worm_protection.go` | `worm_protection.go` |
| DDoS 防护 | `ddos_protection.go` | `ddos_protection.go` | `ddos_protection.go` |
| ARP 防护 | `arp_protection.go` | `arp_protection.go` | `arp_protection.go` |
| 加载配置 | `load_config.go` | `load_config.go` | `load_config.go` |
| 系统配置 | - | `system_config.go` | `system_config.go` |

## 核心 Provider 接口
```go
type MaintenanceProvider interface {
    // 系统配置
    GetSystemConfig(ctx context.Context) (*model.SystemConfig, error)
    UpdateNetworkConfig(ctx context.Context, req model.NetworkConfigRequest) error
    UpdateTemperatureConfig(ctx context.Context, req model.TemperatureConfigRequest) error
    UpdateDeviceInfo(ctx context.Context, req model.DeviceInfoRequest) error
    UpdateDateTime(ctx context.Context, req model.DateTimeRequest) error

    // 重启/保存
    SaveConfig(ctx context.Context) error
    RebootSwitch(ctx context.Context, delay int) error
    FactoryReset(ctx context.Context) error

    // 用户管理
    GetUsers(ctx context.Context) (*model.UserListResponse, error)
    CreateUser(ctx context.Context, req model.UserRequest) error
    DeleteUser(ctx context.Context, username string) error
    DeleteUsers(ctx context.Context, usernames []string) error

    // 会话管理
    GetSessions(ctx context.Context) (*model.SessionListResponse, error)
    DeleteSession(ctx context.Context, sessionID string) error
    DeleteSessions(ctx context.Context, sessionIDs []string) error

    // 日志管理
    GetLogs(ctx context.Context) (*model.LogListResponse, error)
    ClearLogs(ctx context.Context, levels []string) error

    // 文件管理
    GetFiles(ctx context.Context, path string) (*model.FileListResponse, error)
    UploadFile(ctx context.Context, req model.FileUploadRequest) error
    DeleteFile(ctx context.Context, path string) error
    DeleteFiles(ctx context.Context, paths []string) error
    DownloadFile(ctx context.Context, path string) ([]byte, string, error)

    // SNMP 配置
    GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error)
    UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error
    GetTrapHosts(ctx context.Context) ([]model.TrapHost, error)
    AddTrapHost(ctx context.Context, req model.TrapHostRequest) error
    DeleteTrapHost(ctx context.Context, host string) error
    TestTrap(ctx context.Context, host string) error
    GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error)
    AddCommunity(ctx context.Context, name, access, description string) error
    DeleteCommunity(ctx context.Context, name string) error

    // 安全防护
    GetWormRules(ctx context.Context) (*model.WormRuleList, error)
    AddWormRule(ctx context.Context, req model.WormRuleRequest) error
    UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error
    DeleteWormRule(ctx context.Context, id string) error
    DeleteWormRules(ctx context.Context, ids []string) error
    ClearWormStats(ctx context.Context) error
    GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error)
    UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error
    GetARPConfig(ctx context.Context) (*model.ARPConfig, error)
    UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error

    // 加载配置
    GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error)
    LoadConfig(ctx context.Context, configFile string) error
}
```

---

# Network（网络）模块

## 文件组织规范

| 功能模块 | Handler | Provider | DataModel | API 路由前缀 |
|---------|---------|----------|-----------|-------------|
| VLAN 管理 | `vlan.go` | `vlan.go` | `vlan.go` | `/api/v1/network/vlans` |
| 端口管理 | `port.go` | `port.go` | `port.go` | `/api/v1/network/ports` |
| 链路聚合 | `lag.go` | `lag.go` | `lag.go` | `/api/v1/network/lags` |
| STP 管理 | `stp.go` | `stp.go` | `stp.go` | `/api/v1/network/stp` |
| ACL 管理 | `acl.go` | `acl.go` | `acl.go` | `/api/v1/network/acls` |
| IP 路由 | `route_table.go`, `static_route.go` | - | `route_table.go`, `static_route.go` | `/api/v1/routes/*` |

## 核心 Provider 接口
```go
type NetworkProvider interface {
    // VLAN 管理
    GetVLANList(ctx context.Context) (*model.VLANListResponse, error)
    CreateVLAN(ctx context.Context, req model.VLANRequest) error
    UpdateVLAN(ctx context.Context, id int, req model.VLANRequest) error
    DeleteVLAN(ctx context.Context, id int) error
    DeleteVLANs(ctx context.Context, ids []int) error
    AddVLANPort(ctx context.Context, vlanID int, port string, mode string) error
    RemoveVLANPort(ctx context.Context, vlanID int, port string) error

    // 端口管理
    GetPortList(ctx context.Context) (*model.NetworkPortListResponse, error)
    GetPortDetail(ctx context.Context, portName string) (*model.PortDetail, error)
    UpdatePort(ctx context.Context, portName string, req model.PortUpdateRequest) error
    ResetPort(ctx context.Context, portName string) error
    RestartPort(ctx context.Context, portName string) error

    // 链路聚合管理
    GetLAGList(ctx context.Context) (*model.LAGListResponse, error)
    CreateLAG(ctx context.Context, req model.LAGRequest) error
    UpdateLAG(ctx context.Context, id int, req model.LAGRequest) error
    DeleteLAG(ctx context.Context, id int) error
    AddLAGPort(ctx context.Context, lagID int, port string) error
    RemoveLAGPort(ctx context.Context, lagID int, port string) error

    // STP 管理
    GetSTPConfig(ctx context.Context) (*model.STPConfig, error)
    UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error
    GetSTPStatus(ctx context.Context) (*model.STPStatus, error)

    // ACL 管理
    GetACLList(ctx context.Context) (*model.ACLListResponse, error)
    CreateACL(ctx context.Context, req model.ACLRequest) error
    UpdateACL(ctx context.Context, id int, req model.ACLRequest) error
    DeleteACL(ctx context.Context, id int) error
    GetACLRules(ctx context.Context, aclID int) (*model.ACLRuleListResponse, error)
    AddACLRule(ctx context.Context, aclID int, req model.ACLRuleRequest) error
    UpdateACLRule(ctx context.Context, aclID int, ruleID int, req model.ACLRuleRequest) error
    DeleteACLRule(ctx context.Context, aclID int, ruleID int) error
}
```

---

# Config（配置）模块

## 文件组织规范

| 功能模块 | Handler | Provider | DataModel |
|---------|---------|----------|-----------|
| 端口配置 | `port.go` | `config_port.go` | `port.go` |
| 链路聚合 | `lag.go` | `config_lag.go` | `lag.go` |
| VLAN 配置 | `vlan_config.go` | `config_vlan.go` | `vlan_config.go` |
| 风暴控制 | `storm_control.go` | `config_storm.go` | `storm_control.go` |
| 流量控制 | `flow_control.go` | `config_flow.go` | `flow_control.go` |
| 端口隔离 | `port_isolation.go` | `config_isolation.go` | `port_isolation.go` |
| 端口监控 | `port_monitor.go` | `config_monitor.go` | `port_monitor.go` |
| MAC 地址表 | `mac_table.go` | `config_mac.go` | `mac_table.go` |
| STP 配置 | `stp_config.go` | `config_stp.go` | `stp_config.go` |
| ERPS 配置 | `erps.go` | `config_erps.go` | `erps.go` |
| PoE 配置 | `poe.go` | `config_poe.go` | `poe.go` |
| 端口镜像 | `port_mirror.go` | `config_mirror.go` | `port_mirror.go` |
| 组播配置 | `multicast.go` | `config_multicast.go` | `multicast.go` |
| 资源管理 | `resource.go` | `config_resource.go` | `resource.go` |
| 堆叠配置 | `stack.go` | `config_stack.go` | `stack.go` |

## 核心 Provider 接口
```go
type ConfigProvider interface {
    // 端口配置
    GetPortList(ctx context.Context) (*model.PortConfigListResponse, error)
    GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error)
    UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error

    // 链路聚合配置
    GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error)
    CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error
    UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error
    DeleteLinkAggregation(ctx context.Context, id int) error

    // 其他配置（14 项）在各自子文件中定义
}
```

---

# 策略模式 - 双模式切换

## ModeResolver 实现

通过 `internal/service/mode/mode_resolver.go` 实现运行时模式切换：

```go
type RunMode string

const (
    ModeMock   RunMode = "mock"    // 离线测试模式 - 使用数据库模拟数据
    ModeSwitch RunMode = "switch"  // 交换机模式 - 使用真实交换机硬件
)

type ModeResolver struct {
    mu sync.RWMutex
    currentMode RunMode
    configDAO ConfigDAO

    // Provider 缓存
    mockDiagnosticProvider      *diagnostic.MockProvider
    cliDiagnosticProvider       *diagnostic.CLIProvider
    mockMaintenanceProvider     *maintenance.MockProvider
    cliMaintenanceProvider      *maintenance.CLIProvider
    mockNetworkProvider         *network.MockProvider
    cliNetworkProvider          *network.CLIProvider
    mockConfigProvider          *config.MockProvider
    cliConfigProvider           *config.CLIProvider
}

func (r *ModeResolver) GetDiagnosticProvider() provider.DiagnosticProvider
func (r *ModeResolver) GetMaintenanceProvider() provider.MaintenanceProvider
func (r *ModeResolver) GetNetworkProvider() provider.NetworkProvider
func (r *ModeResolver) GetConfigProvider() provider.ConfigProvider
```

## 模式切换 API

```
GET  /api/mode          - 获取当前模式
POST /api/mode          - 切换模式 {"mode": "mock"|"switch"}
GET  /api/system/config - 获取系统配置（含模式信息）
```

---

# 新增功能步骤

## 通用步骤

1. **定义数据模型** (`internal/model/*.go`)
   - 添加新的 struct 定义

2. **扩展 Provider 接口** (`internal/service/provider/interface.go` 或子目录 `interface.go`)
   - 在对应 Provider 接口添加方法签名

3. **实现 Service 代理** (`internal/service/*_service.go`)
   - 添加代理方法，通过 `getProvider()` 调用

4. **实现 Handler** (`internal/handler/*/`)
   - 添加 HTTP 处理方法

5. **实现 Provider** (`internal/service/provider/{module}/cli.go` 和 `mock.go`)
   - CLI Provider: 真实交换机调用
   - Mock Provider: 数据库模拟数据

6. **实现 DataModel** (`internal/datamodel/*/`)
   - 添加 UI 页面（GoAdmin Panel）

7. **注册路由** (`cmd/main.go`)
   - API 路由：`r.GET/POST/PUT/DELETE(...)`
   - 页面路由：`e.HTML("GET", path, datamodelFunc, false)`

---

# 路由注册总览 (`cmd/main.go`)

## API 路由

### System
```go
r.GET("/health", ...)
r.GET("/api/mode", ...)
r.POST("/api/mode", ...)
r.GET("/api/system/config", ...)
```

### Diagnostic
```go
r.GET("/api/v1/diagnostic/cable/ports", ...)
r.POST("/api/v1/diagnostic/cable", ...)
r.POST("/api/v1/diagnostic/ping", ...)
r.GET("/api/v1/diagnostic/ping/:task_id", ...)
r.DELETE("/api/v1/diagnostic/ping/:task_id", ...)
r.POST("/api/v1/diagnostic/traceroute", ...)
r.GET("/api/v1/diagnostic/traceroute/:task_id", ...)
r.DELETE("/api/v1/diagnostic/traceroute/:task_id", ...)
```

### Maintenance
```go
r.POST("/api/v1/system/save-config", ...)
r.POST("/api/v1/system/reboot", ...)
r.POST("/api/v1/system/factory-reset", ...)
r.GET("/api/v1/system/config", ...)
r.PUT("/api/v1/system/network", ...)
r.PUT("/api/v1/system/temperature", ...)
r.PUT("/api/v1/system/info", ...)
r.PUT("/api/v1/system/datetime", ...)
r.GET("/api/v1/config/files", ...)
r.POST("/api/v1/config/load", ...)
r.GET("/api/v1/files", ...)
r.POST("/api/v1/files/upload", ...)
r.POST("/api/v1/files/firmware", ...)
r.GET("/api/v1/files/download", ...)
r.DELETE("/api/v1/files", ...)
r.GET("/api/v1/logs", ...)
r.DELETE("/api/v1/logs", ...)
r.GET("/api/v1/snmp/config", ...)
r.PUT("/api/v1/snmp/config", ...)
r.GET("/api/v1/snmp/communities", ...)
r.POST("/api/v1/snmp/communities", ...)
r.DELETE("/api/v1/snmp/communities/:name", ...)
r.GET("/api/v1/snmp/trap/config", ...)
r.PUT("/api/v1/snmp/trap/config", ...)
r.GET("/api/v1/snmp/trap/hosts", ...)
r.POST("/api/v1/snmp/trap/hosts", ...)
r.DELETE("/api/v1/snmp/trap/hosts/:id", ...)
r.POST("/api/v1/snmp/trap/hosts/:id/test", ...)
r.GET("/api/v1/security/worm/rules", ...)
r.POST("/api/v1/security/worm/rules", ...)
r.PUT("/api/v1/security/worm/rules/:id", ...)
r.DELETE("/api/v1/security/worm/rules", ...)
r.POST("/api/v1/security/worm/clear-stats", ...)
r.GET("/api/v1/security/ddos/config", ...)
r.PUT("/api/v1/security/ddos/config", ...)
r.GET("/api/v1/security/arp/config", ...)
r.PUT("/api/v1/security/arp/config", ...)
r.GET("/api/v1/users", ...)
r.POST("/api/v1/users", ...)
r.PUT("/api/v1/users/:username", ...)
r.DELETE("/api/v1/users", ...)
r.GET("/api/v1/sessions", ...)
r.DELETE("/api/v1/sessions/:session_id", ...)
```

### Network
```go
r.GET("/api/v1/routes/table", ...)
r.GET("/api/v1/routes/static", ...)
r.GET("/api/v1/routes/static/:id", ...)
r.POST("/api/v1/routes/static", ...)
r.PUT("/api/v1/routes/static/:id", ...)
r.DELETE("/api/v1/routes/static/:id", ...)
r.GET("/api/v1/network/vlans", ...)
r.POST("/api/v1/network/vlans", ...)
r.PUT("/api/v1/network/vlans/:id", ...)
r.DELETE("/api/v1/network/vlans/:id", ...)
r.DELETE("/api/v1/network/vlans", ...)
r.POST("/api/v1/network/vlans/:id/ports", ...)
r.DELETE("/api/v1/network/vlans/:id/ports", ...)
r.GET("/api/v1/network/ports", ...)
r.GET("/api/v1/network/ports/:name", ...)
r.PUT("/api/v1/network/ports/:name", ...)
r.POST("/api/v1/network/ports/:name/reset", ...)
r.POST("/api/v1/network/ports/:name/restart", ...)
r.GET("/api/v1/network/lags", ...)
r.POST("/api/v1/network/lags", ...)
r.PUT("/api/v1/network/lags/:id", ...)
r.DELETE("/api/v1/network/lags/:id", ...)
r.POST("/api/v1/network/lags/:id/ports", ...)
r.DELETE("/api/v1/network/lags/:id/ports", ...)
r.GET("/api/v1/network/stp/config", ...)
r.PUT("/api/v1/network/stp/config", ...)
r.GET("/api/v1/network/stp/status", ...)
r.GET("/api/v1/network/acls", ...)
r.POST("/api/v1/network/acls", ...)
r.PUT("/api/v1/network/acls/:id", ...)
r.DELETE("/api/v1/network/acls/:id", ...)
r.GET("/api/v1/network/acls/:id/rules", ...)
r.POST("/api/v1/network/acls/:id/rules", ...)
r.PUT("/api/v1/network/acls/:id/rules/:ruleID", ...)
r.DELETE("/api/v1/network/acls/:id/rules/:ruleID", ...)
```

### Config
```go
r.GET("/api/v1/config/ports", ...)
r.GET("/api/v1/config/ports/:port_id", ...)
r.PUT("/api/v1/config/ports/:port_id", ...)
r.GET("/api/v1/link-aggregation", ...)
r.POST("/api/v1/link-aggregation", ...)
r.PUT("/api/v1/link-aggregation/:id", ...)
r.DELETE("/api/v1/link-aggregation/:id", ...)
```

## GoAdmin 页面路由

```go
// Dashboard 和系统配置
e.HTML("GET", "/admin/dashboard", systemDatamodel.GetDashboardContent, false)
e.HTML("GET", "/admin/system/config", systemDatamodel.GetSystemConfigPage, false)

// 网络模块 - IP 路由
e.HTML("GET", "/admin/network/route-table", networkDatamodel.GetRouteTableContent, false)
e.HTML("GET", "/admin/network/static-route", networkDatamodel.GetStaticRouteContent, false)

// 网络模块 - 诊断工具
e.HTML("GET", "/admin/network/ping", diagnosticDatamodel.GetPingContent, false)
e.HTML("GET", "/admin/network/traceroute", diagnosticDatamodel.GetTracerouteContent, false)
e.HTML("GET", "/admin/network/cable-test", diagnosticDatamodel.GetCableTestContent, false)

// 维护模块 (12 页面)
e.HTML("GET", "/admin/maintenance/reboot-save", maintDatamodel.GetRebootSaveContent, false)
e.HTML("GET", "/admin/maintenance/users", maintDatamodel.GetUsersContent, false)
e.HTML("GET", "/admin/maintenance/system-config", maintDatamodel.GetMaintenanceSystemConfigContent, false)
e.HTML("GET", "/admin/maintenance/load-config", maintDatamodel.GetLoadConfigContent, false)
e.HTML("GET", "/admin/maintenance/files", maintDatamodel.GetFilesContent, false)
e.HTML("GET", "/admin/maintenance/logs", maintDatamodel.GetLogsContent, false)
e.HTML("GET", "/admin/maintenance/snmp", maintDatamodel.GetSNMPContent, false)
e.HTML("GET", "/admin/maintenance/snmp-trap", maintDatamodel.GetSNMPTrapContent, false)
e.HTML("GET", "/admin/maintenance/worm-protection", maintDatamodel.GetWormProtectionContent, false)
e.HTML("GET", "/admin/maintenance/ddos-protection", maintDatamodel.GetDDoSProtectionContent, false)
e.HTML("GET", "/admin/maintenance/arp-protection", maintDatamodel.GetARPProtectionContent, false)
e.HTML("GET", "/admin/maintenance/sessions", maintDatamodel.GetSessionsContent, false)

// 网络模块 (5 页面)
e.HTML("GET", "/admin/network/vlan", networkDatamodel.GetVLANContent, false)
e.HTML("GET", "/admin/network/port", networkDatamodel.GetPortContent, false)
e.HTML("GET", "/admin/network/lag", networkDatamodel.GetLAGContent, false)
e.HTML("GET", "/admin/network/stp", networkDatamodel.GetSTPContent, false)
e.HTML("GET", "/admin/network/acl", networkDatamodel.GetACLContent, false)

// 配置模块 (14 页面)
e.HTML("GET", "/admin/config/ports", configDatamodel.GetPortsContent, false)
e.HTML("GET", "/admin/config/link-aggregation", configDatamodel.GetLinkAggregationContent, false)
e.HTML("GET", "/admin/config/storm-control", configDatamodel.GetStormControlContent, false)
e.HTML("GET", "/admin/config/flow-control", configDatamodel.GetFlowControlContent, false)
e.HTML("GET", "/admin/config/port-isolation", configDatamodel.GetPortIsolationContent, false)
e.HTML("GET", "/admin/config/port-monitor", configDatamodel.GetPortMonitorContent, false)
e.HTML("GET", "/admin/config/vlan", configDatamodel.GetVLANContent, false)
e.HTML("GET", "/admin/config/mac-table", configDatamodel.GetMacTableContent, false)
e.HTML("GET", "/admin/config/stp", configDatamodel.GetSTPContent, false)
e.HTML("GET", "/admin/config/erps", configDatamodel.GetERPSContent, false)
e.HTML("GET", "/admin/config/poe", configDatamodel.GetPoEContent, false)
e.HTML("GET", "/admin/config/port-mirror", configDatamodel.GetPortMirrorContent, false)
e.HTML("GET", "/admin/config/multicast", configDatamodel.GetMulticastContent, false)
e.HTML("GET", "/admin/config/resource", configDatamodel.GetResourceContent, false)
e.HTML("GET", "/admin/config/stack", configDatamodel.GetStackContent, false)
```

---

# 模块对比总结

| 特性 | Diagnostic | Maintenance | Network | Config |
|------|------------|-------------|---------|--------|
| 核心功能 | 网络诊断 | 系统维护、安全防护 | 网络管理 | 配置管理 |
| Handler 文件数 | 1 | 11 | 7 | 17 |
| DataModel 文件数 | 3 | 13 | 7 | 15 |
| Service | `DiagnosticService` | `MaintenanceService` | `NetworkService` | `ConfigService` |
| Provider 接口 | `DiagnosticProvider` | `MaintenanceProvider` | `NetworkProvider` | `ConfigProvider` |
| 异步任务管理 | ✓ (Ping/Traceroute/Cable) | - | - | - |
| 模式切换 | ✓ | ✓ | ✓ | ✓ |
| 单例模式 | ✓ | ✓ | ✓ | ✓ |
| 线程安全 | ✓ | ✓ | ✓ | ✓ |
| GoAdmin 页面 | 3 | 12 | 7 | 15 |
| API 端点 | 7 | 30+ | 25+ | 4 |