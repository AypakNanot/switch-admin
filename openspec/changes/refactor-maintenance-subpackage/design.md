# Design Document - refactor-maintenance-subpackage

## 1. 概述

本文档描述维护模块按功能分组重构的详细设计，包括目录结构、包组织、导入路径等。

## 2. 当前问题

### 2.1 文件集中度

虽然已经拆分为 24 个子模块文件，但仍在同一包目录下：

```
internal/handler/
├── maintenance_handler.go           (263 行)
├── maintenance_reboot_save.go       (83 行)
├── maintenance_system_config.go     (105 行)
├── maintenance_file_management.go   (131 行)
├── ... (共 13 个文件)
```

### 2.2 维护问题

- **导航困难**: 13 个 maintenance_* 文件在同一目录，难以快速定位
- **命名冗长**: 每个文件都需要 `maintenance_` 前缀来避免冲突
- **扩展成本高**: 新增子模块时文件列表更长

## 3. 设计原则

### 3.1 单一包职责

维护模块是一个包，所有功能文件共享同一个包名 (`maintenance`)。

### 3.2 按功能分文件

每个功能子模块是独立的 `.go` 文件，文件名表达功能归属。

### 3.3 导入清晰

导入路径简单明了：
```go
import (
    maintHandler "switch-admin/internal/handler/maintenance"
    maintDatamodel "switch-admin/internal/datamodel/maintenance"
)
```

## 4. 详细设计

### 4.1 Handler 层包结构

```
internal/handler/maintenance/
├── handler.go              # 主处理器（Handler 结构体 + New 函数）
├── reboot_save.go          # 重启/保存 Handler 方法
├── system_config.go        # 系统配置 Handler 方法
├── file_management.go      # 文件管理 Handler 方法
├── log_management.go       # 日志管理 Handler 方法
├── snmp_config.go          # SNMP 配置 Handler 方法
├── snmp_trap.go            # SNMP Trap Handler 方法
├── worm_protection.go      # 蠕虫防护 Handler 方法
├── ddos_protection.go      # DDoS 防护 Handler 方法
├── arp_protection.go       # ARP 防护 Handler 方法
├── user_management.go      # 用户管理 Handler 方法
├── session_management.go   # 会话管理 Handler 方法
└── load_config.go          # 加载配置 Handler 方法
```

#### 4.1.1 Handler 结构体定义

```go
// File: internal/handler/maintenance/handler.go
package maintenance

type Handler struct{}

func New() *Handler {
    return &Handler{}
}

// 重启/保存相关方法也在 handler.go 中
func (h *Handler) SaveConfig(c *gin.Context) {
    // 实现...
}
```

#### 4.1.2 功能文件组织

```go
// File: internal/handler/maintenance/system_config.go
package maintenance

func (h *Handler) GetSystemConfig(c *gin.Context) {
    // 实现...
}

func (h *Handler) UpdateNetworkConfig(c *gin.Context) {
    // 实现...
}
```

所有方法都添加到同一个 `Handler` 结构体上，只是分布在不同的文件中。

### 4.2 DataModel 层包结构

```
internal/datamodel/maintenance/
├── maintenance.go          # 导出页面函数（保持兼容）
├── reboot_save.go          # 重启/保存页面内容
├── system_config.go        # 系统配置页面内容
├── file_management.go      # 文件管理页面内容
├── log_management.go       # 日志管理页面内容
├── snmp_config.go          # SNMP 配置页面内容
├── snmp_trap.go            # SNMP Trap 页面内容
├── worm_protection.go      # 蠕虫防护页面内容
├── ddos_protection.go      # DDoS 防护页面内容
├── arp_protection.go       # ARP 防护页面内容
├── user_management.go      # 用户管理页面内容
├── session_management.go   # 会话管理页面内容
└── load_config.go          # 加载配置页面内容
```

#### 4.2.1 导出函数模板

```go
// File: internal/datamodel/maintenance/maintenance.go
package maintenance

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/template/types"
)

// GetRebootSaveContent 重启/保存配置页面
func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    return getRebootSaveContent(ctx)
}

// GetUsersContent 用户管理页面
func GetUsersContent(ctx *context.Context) (types.Panel, error) {
    return getUsersContent(ctx)
}

// ... 其他导出函数
```

#### 4.2.2 内部实现函数

```go
// File: internal/datamodel/maintenance/reboot_save.go
package maintenance

import (
    "html/template"
    "github.com/GoAdminGroup/go-admin/context"
    tmpl "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/types"
)

// getRebootSaveContent 重启/保存配置页面（内部函数）
func getRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    components := tmpl.Default(ctx)
    // ... 页面内容
}
```

### 4.3 导入路径映射

| 原路径 | 新路径 |
|-------|-------|
| `internal/handler/maintenance_reboot_save.go` | `internal/handler/maintenance/reboot_save.go` |
| `internal/handler/maintenance_system_config.go` | `internal/handler/maintenance/system_config.go` |
| `internal/datamodel/maintenance_reboot_save.go` | `internal/datamodel/maintenance/reboot_save.go` |
| `internal/datamodel/maintenance_system_config.go` | `internal/datamodel/maintenance/system_config.go` |

### 4.4 路由更新

```go
// File: cmd/main.go

// 之前
import "switch-admin/internal/handler"
maintenanceHandler := handler.NewMaintenanceHandler()

// 之后
import maintHandler "switch-admin/internal/handler/maintenance"
maintenanceHandler := maintHandler.New()
```

## 5. 包依赖关系

```
cmd/main.go
    ↓
internal/handler/maintenance (包)
    ↓
internal/service (如有)
    ↓
internal/dao
```

所有功能文件在同一个包内，**文件间无导入成本**，方法直接调用。

## 6. 迁移策略

### 6.1 批量迁移

由于是同一包内文件重组，可以批量进行：

```bash
# 1. 创建新目录
mkdir -p internal/handler/maintenance
mkdir -p internal/datamodel/maintenance

# 2. 移动并重命名文件
mv internal/handler/maintenance_*.go internal/handler/maintenance/
mv internal/datamodel/maintenance_*.go internal/datamodel/maintenance/

# 3. 批量修改包名和文件名
# 4. 编译验证
go build
```

### 6.2 回滚方案

如遇问题，可快速回滚：

```bash
git reset --hard HEAD
```

## 7. 验收标准

### 7.1 代码组织

- [ ] 所有功能文件在同一包内
- [ ] 文件名表达功能归属
- [ ] 包名统一为 `maintenance`

### 7.2 编译质量

- [ ] `go build` 无错误
- [ ] `go vet ./...` 通过
- [ ] `go fmt ./...` 格式化

### 7.3 功能验证

- [ ] 所有 API 正常响应
- [ ] 所有 UI 页面正常显示
- [ ] 无功能回归

## 8. 优势总结

### 8.1 与子包方案对比

| 维度 | 同包分文件方案 | 子包方案 |
|------|--------------|---------|
| **导入复杂度** | 低（1 个导入） | 高（12 个导入） |
| **文件导航** | 中（平铺直叙） | 低（需进入子目录） |
| **方法调用** | 零成本（同包） | 跨包调用 |
| **扩展成本** | 低（添加文件） | 中（添加子包） |
| **代码量** | 少 | 多（需要委托代码） |

### 8.2 适用场景

本方案适合：
- 功能模块数量适中（10-20 个）
- 模块间共享数据和状态
- 需要频繁协作的功能

子包方案适合：
- 功能模块相对独立
- 需要独立测试
- 可能独立复用

## 9. 后续优化建议

### 9.1 单元测试

可在同包内编写测试：

```
internal/handler/maintenance/
├── handler.go
├── system_config.go
├── system_config_test.go  # 测试文件
```

### 9.2 代码分组

如文件数量继续增长，可考虑添加注释分组：

```go
//
// 系统配置
//

func (h *Handler) GetSystemConfig(...) { ... }
func (h *Handler) UpdateNetworkConfig(...) { ... }

//
// 文件管理
//

func (h *Handler) GetFiles(...) { ... }
```
