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

# Maintenance 模块分析方法论

## 核心架构模式

Maintenance 模块采用 **分层架构 + 策略模式**，核心设计原则：

### 1. 四层架构

```
Handler 层 (internal/handler/maintenance/)
    ↓ HTTP 请求处理
Service 层 (internal/service/maintenance_service.go)
    ↓ 业务代理 + 模式切换
Provider 接口层 (internal/service/provider/interface.go)
    ↓ 接口实现
CLI/Mock Provider (internal/service/provider/cli|mock/)
    ↓ 数据源
DataModel 层 (internal/datamodel/maintenance/)
    ↓ UI 渲染
```

### 2. 文件组织规范

**按功能子模块拆分** - 每个子功能独立成文件：
- `reboot_save.go` - 重启/保存配置
- `user_management.go` - 用户管理
- `session_management.go` - 会话管理
- `log_management.go` - 日志管理
- `file_management.go` - 文件管理
- `snmp_config.go` - SNMP 配置
- `worm_protection.go` - 蠕虫防护
- `ddos_protection.go` - DDoS 防护
- `arp_protection.go` - ARP 防护
- `load_config.go` - 加载配置

**跨层文件对应** - 每一层文件命名保持一致，便于导航。

### 3. 各层实现规范

#### Handler 层
```go
type Handler struct {
    service *service.MaintenanceService
}

func New() *Handler {
    return &Handler{service: service.GetMaintenanceService()}
}

// 方法命名与 HTTP 路由对应
// GET /api/v1/users -> GetUsers()
// POST /api/v1/users -> CreateUser()
```

#### Service 层
```go
// 单例模式 + 线程安全
var maintenanceService *MaintenanceService
var maintenanceOnce sync.Once

func (s *MaintenanceService) GetUsers(ctx context.Context) (...) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.getProvider().GetUsers(ctx)
}
```

#### Provider 层
- **接口定义**: `interface.go` 中统一声明 `MaintenanceProvider`
- **CLI Provider**: 真实交换机调用
- **Mock Provider**: 离线测试模拟数据

#### DataModel 层
```go
func getRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    // HTML/CSS/JS 内嵌 + GoAdmin 组件构建
    // fetch API 调用后端 Handler
}
```

### 4. 策略模式 - 双模式切换

通过 `ModeResolver` 实现运行时模式切换：
- **Mock 模式**: `mock.MaintenanceProvider` - 开发测试
- **CLI 模式**: `cli.MaintenanceProvider` - 生产环境

### 5. 新增功能步骤

1. 在 `interface.go` 的 `MaintenanceProvider` 接口添加方法
2. 在 `service/maintenance_service.go` 添加代理方法
3. 在 `handler/maintenance/` 添加 HTTP Handler
4. 在 `provider/cli/` 和 `provider/mock/` 分别实现
5. 在 `datamodel/maintenance/` 添加 UI 页面（可选）

---

# Network 模块实现方法论

## 架构概览

Network 模块与 Maintenance 模块采用相同的**分层架构 + 策略模式**，实现网络功能管理。

### 1. 四层架构

```
Handler 层 (internal/handler/network/)
    ↓ HTTP 请求处理
Service 层 (internal/service/network_service.go)
    ↓ 业务代理 + 模式切换
Provider 接口层 (internal/service/provider/interface.go)
    ↓ 接口实现
CLI/Mock Provider (internal/service/provider/cli|mock/)
    ↓ 数据源
DataModel 层 (internal/datamodel/network/)
    ↓ UI 渲染
```

### 2. 文件组织规范

**按功能子模块拆分** - 每个子功能独立成文件：

| 功能模块 | Handler | Provider (Mock/CLI) | DataModel |
|---------|---------|---------------------|-----------|
| VLAN 管理 | `vlan.go` | `vlan.go` | `vlan.go` |
| 端口管理 | `port.go` | `port.go` | `port.go` |
| 链路聚合 | `lag.go` | `lag.go` | `lag.go` |
| STP 管理 | `stp.go` | `stp.go` | `stp.go` |
| ACL 管理 | `acl.go` | `acl.go` | `acl.go` |

**跨层文件对应** - 每一层文件命名保持一致，便于导航。

### 3. 各层实现规范

#### Handler 层 (`internal/handler/network/`)
```go
type Handler struct {
    service *service.NetworkService
}

func New() *Handler {
    return &Handler{service: service.GetNetworkService()}
}

// 路由对应示例：
// GET /api/v1/network/vlans -> GetVLANList()
// POST /api/v1/network/vlans -> CreateVLAN()
// PUT /api/v1/network/vlans/:id -> UpdateVLAN()
// DELETE /api/v1/network/vlans/:id -> DeleteVLAN()
```

#### Service 层 (`internal/service/network_service.go`)
```go
// 单例模式 + 线程安全
var networkService *NetworkService
var networkOnce sync.Once

func GetNetworkService() *NetworkService {
    networkOnce.Do(func() {
        networkService = &NetworkService{
            modeResolver: mode.NewModeResolver(...),
        }
    })
    return networkService
}

func (s *NetworkService) GetVLANList(ctx context.Context) (...) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.getProvider().GetVLANList(ctx)
}
```

#### Provider 接口层
```go
// interface.go 中定义
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

#### DataModel 层 (`internal/datamodel/network/`)
```go
// network.go - 页面导出函数
func GetVLANContent(ctx *context.Context) (types.Panel, error) {
    return getVLANContent(ctx)
}

// vlan.go - 实际页面实现
func getVLANContent(ctx *context.Context) (types.Panel, error) {
    // HTML/CSS/JS 内嵌 + GoAdmin 组件构建
    // fetch API 调用后端 Handler (/api/v1/network/vlans)
}
```

### 4. 路由注册 (`cmd/main.go`)

```go
// 导入
networkDatamodel "switch-admin/internal/datamodel/network"
networkHandler "switch-admin/internal/handler/network"

// API 路由注册
r.GET("/api/v1/network/vlans", networkHandler.GetVLANList)
r.POST("/api/v1/network/vlans", networkHandler.CreateVLAN)
// ... 其他路由

// GoAdmin 页面注册
e.HTML("GET", "/admin/network/vlan", networkDatamodel.GetVLANContent, false)
e.HTML("GET", "/admin/network/port", networkDatamodel.GetPortContent, false)
e.HTML("GET", "/admin/network/lag", networkDatamodel.GetLAGContent, false)
e.HTML("GET", "/admin/network/stp", networkDatamodel.GetSTPContent, false)
e.HTML("GET", "/admin/network/acl", networkDatamodel.GetACLContent, false)
```

### 5. 数据模型 (`internal/model/network.go`)

```go
// VLAN 配置
type VLAN struct {
    ID     int      `json:"id"`
    Name   string   `json:"name"`
    Ports  []string `json:"ports"`
    Status string   `json:"status"`
}

// 端口配置
type Port struct {
    Name   string `json:"name"`
    Status string `json:"status"`
    Speed  string `json:"speed"`
    Duplex string `json:"duplex"`
    VLAN   int    `json:"vlan"`
    Type   string `json:"type"`
}

// 链路聚合组
type LAG struct {
    ID     int      `json:"id"`
    Name   string   `json:"name"`
    Ports  []string `json:"ports"`
    Status string   `json:"status"`
    Mode   string   `json:"mode"`
}

// STP 配置
type STPConfig struct {
    Enabled    bool   `json:"enabled"`
    Mode       string `json:"mode"`
    Priority   int    `json:"priority"`
    RootBridge string `json:"root_bridge"`
}

// ACL 配置
type ACL struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Type  string `json:"type"`
    Rules int    `json:"rules"`
    Status string `json:"status"`
}
```

### 6. 新增功能步骤

1. 在 `interface.go` 的 `NetworkProvider` 接口添加方法
2. 在 `service/network_service.go` 添加代理方法
3. 在 `handler/network/` 添加 HTTP Handler
4. 在 `provider/cli/` 和 `provider/mock/` 分别实现
5. 在 `datamodel/network/` 添加 UI 页面
6. 在 `cmd/main.go` 注册 API 路由和页面路由
7. 在 `model/network.go` 添加数据模型（如需要）