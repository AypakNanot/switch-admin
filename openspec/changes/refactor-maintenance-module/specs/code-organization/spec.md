# Code Organization Spec

## 概述

维护模块代码重构，将集中在两个大文件中的代码按功能子模块拆分为多个独立文件，提升可维护性和可读性。

## 能力范围

- Handler 层代码拆分（12 个子模块）
- DataModel 层代码拆分（12 个子模块）
- 保持 API 路由和功能不变

---

## MODIFIED Requirements

### Requirement: Handler 层文件组织

维护模块的 Handler 层代码**MUST**按功能子模块组织为独立文件。

**变更说明：**

- **之前：** 所有 38 个 Handler 方法集中在 `internal/handler/maintenance_handler.go` (1024 行)
- **现在：** 按 12 个子模块拆分为独立文件，每个文件 50-200 行

| 子模块 | 原位置 | 新位置 | 方法数 |
|--------|-------|-------|-------|
| 重启/保存 | maintenance_handler.go | maintenance_reboot_save.go | 3 |
| 系统配置 | maintenance_handler.go | maintenance_system_config.go | 5 |
| 文件管理 | maintenance_handler.go | maintenance_file_management.go | 5 |
| 日志管理 | maintenance_handler.go | maintenance_log_management.go | 2 |
| SNMP 配置 | maintenance_handler.go | maintenance_snmp_config.go | 4 |
| SNMP Trap | maintenance_handler.go | maintenance_snmp_trap.go | 5 |
| 蠕虫防护 | maintenance_handler.go | maintenance_worm_protection.go | 5 |
| DDoS 防护 | maintenance_handler.go | maintenance_ddos_protection.go | 2 |
| ARP 防护 | maintenance_handler.go | maintenance_arp_protection.go | 2 |
| 用户管理 | maintenance_handler.go | maintenance_user_management.go | 4 |
| 会话管理 | maintenance_handler.go | maintenance_session_management.go | 2 |
| 加载配置 | maintenance_handler.go | maintenance_load_config.go | 2 |

#### Scenario: 开发者查找代码

当开发者需要修改 SNMP 配置相关功能时，**MUST**能直接定位到 `maintenance_snmp_config.go` 文件，无需在 1000+ 行文件中搜索。

#### Scenario: 代码审查

审查文件上传功能变更时，PR 只修改 `maintenance_file_management.go`，变更范围清晰。

#### Scenario: 并行开发

开发者 A 修改 `maintenance_snmp_config.go`，开发者 B 修改 `maintenance_user_management.go`，两者不会产生 Git 合并冲突。

---

### Requirement: Handler 子模块结构

每个 Handler 子模块文件**MUST**遵循统一的结构模板。

**代码模板：**

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// RebootSaveHandler 重启/保存子模块处理器
type RebootSaveHandler struct{}

// NewRebootSaveHandler 创建重启/保存处理器
func NewRebootSaveHandler() *RebootSaveHandler {
    return &RebootSaveHandler{}
}

// SaveConfig 保存配置
func (h *RebootSaveHandler) SaveConfig(c *gin.Context) {
    // 实现...
}
```

#### Scenario: 新增子模块

添加新的维护功能（如"批量操作"）时，开发者**MUST**创建 `maintenance_batch_operations.go` 并遵循相同模板。

#### Scenario: 主处理器整合

`MaintenanceHandler` 结构体**MUST**持有所有子模块处理器的引用，通过构造函数统一初始化。

---

### Requirement: DataModel 层文件组织

维护模块的 DataModel 层代码**MUST**按功能子模块组织为独立文件。

**变更说明：**

- **之前：** 所有 12 个页面函数集中在 `internal/datamodel/maintenance.go` (1601 行)
- **现在：** 按 12 个子模块拆分为独立文件，每个文件 100-200 行

| 子模块 | 原位置 | 新位置 | 页面函数 |
|--------|-------|-------|---------|
| 重启/保存 | maintenance.go | maintenance_reboot_save.go | GetRebootSaveContent |
| 系统配置 | maintenance.go | maintenance_system_config.go | GetMaintenanceSystemConfigContent |
| 文件管理 | maintenance.go | maintenance_file_management.go | GetFilesContent |
| 日志管理 | maintenance.go | maintenance_log_management.go | GetLogsContent |
| SNMP 配置 | maintenance.go | maintenance_snmp_config.go | GetSNMPContent |
| SNMP Trap | maintenance.go | maintenance_snmp_trap.go | GetSNMPTrapContent |
| 蠕虫防护 | maintenance.go | maintenance_worm_protection.go | GetWormProtectionContent |
| DDoS 防护 | maintenance.go | maintenance_ddos_protection.go | GetDDoSProtectionContent |
| ARP 防护 | maintenance.go | maintenance_arp_protection.go | GetARPProtectionContent |
| 用户管理 | maintenance.go | maintenance_user_management.go | GetUsersContent |
| 会话管理 | maintenance.go | maintenance_session_management.go | GetSessionsContent |
| 加载配置 | maintenance.go | maintenance_load_config.go | GetLoadConfigContent |

#### Scenario: UI 页面修改

修改 SNMP 配置页面时，开发者**MUST**只修改 `maintenance_snmp_config.go`，不影响其他页面。

---

### Requirement: DataModel 子模块结构

每个 DataModel 子模块文件**MUST**遵循统一的结构模板。

**代码模板：**

```go
package datamodel

import (
    "html/template"
    "github.com/GoAdminGroup/go-admin/context"
    tmpl "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/types"
)

// GetRebootSaveContent 重启/保存页面内容
func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    components := tmpl.Default(ctx)

    content := `
    <div class="maintenance-section">
        <!-- 页面内容 -->
    </div>
    <script>
    function saveConfig() {
        // JavaScript 逻辑
    }
    </script>
    `

    // ... 渲染逻辑
}
```

#### Scenario: 页面样式隔离

每个子模块页面的 CSS 样式**MUST**定义在各自文件内，避免全局样式冲突。

---

### Requirement: 路由映射保持不变

所有 API 路由和 UI 路由**MUST**保持原有映射关系不变。

**变更说明：**

- **之前：** 路由指向 `MaintenanceHandler` 的方法
- **现在：** 路由仍然指向 `MaintenanceHandler`，内部委托给子模块处理器

#### Scenario: API 调用

前端调用 `POST /api/v1/system/save-config` **MUST**正常工作，无需修改任何前端代码。

#### Scenario: UI 访问

用户访问 `/admin/maintenance/snmp` **MUST**正常显示 SNMP 配置页面。

---

### Requirement: 共享依赖处理

子模块间共享的依赖（如 ConfigDAO）**MUST**通过主处理器统一注入。

**代码示例：**

```go
// MaintenanceHandler 持有共享依赖
type MaintenanceHandler struct {
    configDAO *dao.ConfigDAO
    rebootSaveHandler *RebootSaveHandler
    // ...
}

// 构造函数中注入依赖
func NewMaintenanceHandler() *MaintenanceHandler {
    return &MaintenanceHandler{
        configDAO: dao.NewConfigDAO(),
        rebootSaveHandler: NewRebootSaveHandler(),
        // ...
    }
}
```

#### Scenario: ConfigDAO 复用

所有子模块**MUST**使用同一个 ConfigDAO 实例，避免重复创建数据库连接。

---

## RENAMED Requirements

### Requirement: 主处理器职责

`maintenance_handler.go` 的职责**MUST**重新定义为：
- 持有所有子模块处理器引用
- 提供统一的构造函数
- 可选：提供共享依赖（如 ConfigDAO）

**变更说明：**

- **之前：** `MaintenanceHandler` 包含所有 38 个方法
- **现在：** `MaintenanceHandler` 作为容器，子模块处理器作为成员

#### Scenario: 方法委托

`h.SaveConfig()` **MUST**内部调用 `h.rebootSaveHandler.SaveConfig()`。

---

## 验收标准

1. [x] 创建 24 个子模块文件（12 个 Handler + 12 个 DataModel）
2. [x] `go build` 编译无错误
3. [x] `go vet ./...` 检查通过
4. [x] 所有 API 端点功能正常（38 个方法）
5. [x] 所有 UI 页面显示正常（12 个页面）
6. [x] 每个文件不超过 300 行
7. [x] 代码通过 `go fmt` 格式化

## 非目标

本变更**不涉及**：
- API 路由变更
- 功能逻辑修改
- UI 样式调整
- 数据库结构变更
