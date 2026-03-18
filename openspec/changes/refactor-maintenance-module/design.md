# Design Document - refactor-maintenance-module

## 1. 概述

本文档描述维护模块代码重构的详细设计，包括文件结构、接口定义、代码组织方式等。

## 2. 当前问题

### 2.1 代码集中度问题

| 文件 | 代码行数 | 功能数量 | 平均每个功能 |
|------|---------|---------|------------|
| maintenance_handler.go | 1024 行 | 38 个方法 | ~27 行/方法 |
| maintenance.go | 1601 行 | 12 个页面 | ~133 页/页面 |

### 2.2 维护困难

- **导航困难**: 在 1000+ 行文件中定位特定功能需要大量滚动
- **合并冲突**: 多人修改同一文件容易产生冲突
- **代码审查**: 变更范围不清晰，审查成本高
- **理解成本**: 新成员难以快速理解模块结构

## 3. 设计原则

### 3.1 单一职责原则 (SRP)

每个文件只负责一个功能子模块，例如：
- `maintenance_reboot_save.go` 只处理重启/保存相关功能
- `maintenance_snmp_config.go` 只处理 SNMP 配置相关功能

### 3.2 高内聚低耦合

- **内聚**: 相关功能组织在同一文件中
- **低耦合**: 子模块间无直接依赖，通过主处理器协调

### 3.3 向后兼容

- 保持现有 API 路由不变
- 保持现有功能行为不变
- 保持现有测试用例有效性

## 4. 详细设计

### 4.1 Handler 层文件结构

#### 4.1.1 主处理器 (maintenance_handler.go)

```go
package handler

import (
    "github.com/gin-gonic/gin"
)

// MaintenanceHandler 维护模块主处理器
type MaintenanceHandler struct {
    rebootSaveHandler       *RebootSaveHandler
    systemConfigHandler     *SystemConfigHandler
    fileManagementHandler   *FileManagementHandler
    logManagementHandler    *LogManagementHandler
    snmpConfigHandler       *SNMPConfigHandler
    snmpTrapHandler         *SNMPTrapHandler
    wormProtectionHandler   *WormProtectionHandler
    ddosProtectionHandler   *DDoSProtectionHandler
    arpProtectionHandler    *ARPProtectionHandler
    userManagementHandler   *UserManagementHandler
    sessionManagementHandler *SessionManagementHandler
    loadConfigHandler       *LoadConfigHandler
}

// NewMaintenanceHandler 创建维护模块处理器
func NewMaintenanceHandler() *MaintenanceHandler {
    return &MaintenanceHandler{
        rebootSaveHandler:       NewRebootSaveHandler(),
        systemConfigHandler:     NewSystemConfigHandler(),
        fileManagementHandler:   NewFileManagementHandler(),
        logManagementHandler:    NewLogManagementHandler(),
        snmpConfigHandler:       NewSNMPConfigHandler(),
        snmpTrapHandler:         NewSNMPTrapHandler(),
        wormProtectionHandler:   NewWormProtectionHandler(),
        ddosProtectionHandler:   NewDDoSProtectionHandler(),
        arpProtectionHandler:    NewARPProtectionHandler(),
        userManagementHandler:   NewUserManagementHandler(),
        sessionManagementHandler: NewSessionManagementHandler(),
        loadConfigHandler:       NewLoadConfigHandler(),
    }
}
```

#### 4.1.2 子模块处理器模板

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
// POST /api/v1/system/save-config
func (h *RebootSaveHandler) SaveConfig(c *gin.Context) {
    // 实现...
}

// RebootSwitch 重启交换机
// POST /api/v1/system/reboot
func (h *RebootSaveHandler) RebootSwitch(c *gin.Context) {
    // 实现...
}

// FactoryReset 恢复出厂配置
// POST /api/v1/system/factory-reset
func (h *RebootSaveHandler) FactoryReset(c *gin.Context) {
    // 实现...
}
```

### 4.2 DataModel 层文件结构

#### 4.2.1 主页面 (maintenance.go)

保留必要的框架结构和通用组件：

```go
package datamodel

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/template/types"
)

// GetMaintenanceLayout 获取维护模块通用布局
func GetMaintenanceLayout(ctx *context.Context) (types.Panel, error) {
    // 通用布局组件
}
```

#### 4.2.2 子模块页面模板

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
    <style>
        .maintenance-section { margin-bottom: 30px; }
    </style>

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

### 4.3 子模块详细映射

| 子模块 | Handler 文件 | DataModel 文件 | API 路由数量 |
|--------|-------------|---------------|-------------|
| 重启/保存 | maintenance_reboot_save.go | maintenance_reboot_save.go | 3 |
| 系统配置 | maintenance_system_config.go | maintenance_system_config.go | 5 |
| 文件管理 | maintenance_file_management.go | maintenance_file_management.go | 5 |
| 日志管理 | maintenance_log_management.go | maintenance_log_management.go | 2 |
| SNMP 配置 | maintenance_snmp_config.go | maintenance_snmp_config.go | 4 |
| SNMP Trap | maintenance_snmp_trap.go | maintenance_snmp_trap.go | 5 |
| 蠕虫防护 | maintenance_worm_protection.go | maintenance_worm_protection.go | 5 |
| DDoS 防护 | maintenance_ddos_protection.go | maintenance_ddos_protection.go | 2 |
| ARP 防护 | maintenance_arp_protection.go | maintenance_arp_protection.go | 2 |
| 用户管理 | maintenance_user_management.go | maintenance_user_management.go | 4 |
| 会话管理 | maintenance_session_management.go | maintenance_session_management.go | 2 |
| 加载配置 | maintenance_load_config.go | maintenance_load_config.go | 2 |

## 5. 共享依赖处理

### 5.1 ConfigDAO

所有子模块共享 `ConfigDAO` 依赖：

```go
// 方案 A: 主处理器持有共享依赖
type MaintenanceHandler struct {
    configDAO *dao.ConfigDAO
    rebootSaveHandler *RebootSaveHandler
    // ...
}

type RebootSaveHandler struct {
    configDAO *dao.ConfigDAO  // 从主处理器传入
}
```

### 5.2 通用工具函数

如有通用工具函数，提取到 `internal/handler/maintenance_utils.go`：

```go
package handler

// validateConfirmation 验证高危操作确认
func validateConfirmation(confirmation string) bool {
    return confirmation == "CONFIRM"
}
```

## 6. 路由映射保持不变

`cmd/main.go` 中的路由注册无需修改：

```go
// 路由保持不变
maintenanceHandler := handler.NewMaintenanceHandler()
r.POST("/api/v1/system/save-config", maintenanceHandler.SaveConfig)
r.POST("/api/v1/system/reboot", maintenanceHandler.RebootSwitch)
// ... 其他路由
```

## 7. 实施策略

### 7.1 分阶段实施

**阶段一**: Handler 层拆分
1. 创建子模块文件 (1-6)
2. 创建子模块文件 (7-12)
3. 重构主处理器
4. 编译验证

**阶段二**: DataModel 层拆分
5. 创建子模块文件 (1-6)
6. 创建子模块文件 (7-12)
7. 重构主页面
8. 编译验证

**阶段三**: 测试验证
9. 启动服务测试
10. 代码审查与清理

### 7.2 回滚策略

如遇到问题，可通过 git 快速回滚：

```bash
git reset --hard HEAD
```

## 8. 验收标准

### 8.1 代码质量

- [ ] 每个文件不超过 300 行
- [ ] 每个方法不超过 50 行
- [ ] 通过 `go fmt` 格式化
- [ ] 通过 `go vet` 检查

### 8.2 功能验证

- [ ] 所有 API 端点响应正常
- [ ] 所有 UI 页面显示正常
- [ ] 无功能回归

### 8.3 文档

- [ ] 更新 tasks.md 记录进度
- [ ] 提交信息清晰描述变更

## 9. 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|-------|------|---------|
| 路由映射错误 | 低 | 高 | 编译测试 + 功能测试 |
| 共享依赖遗漏 | 中 | 中 | 代码审查检查 |
| 循环导入 | 低 | 高 | 使用 go build 验证 |
| 代码重复 | 中 | 低 | 提取公共函数 |

## 10. 未来扩展

重构后的架构支持：

1. **独立测试**: 可为每个子模块编写单元测试
2. **文档生成**: 可按子模块生成 API 文档
3. **权限控制**: 可按子模块细化权限
4. **性能优化**: 可针对特定子模块优化
