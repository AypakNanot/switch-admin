# refactor-maintenance-module

维护模块代码重构 - 按功能子模块拆分文件

## 概述

当前维护模块（Maintenance）的所有代码都集中在两个大文件中：
- `internal/handler/maintenance_handler.go` (1024 行)
- `internal/datamodel/maintenance.go` (1601 行)

这两个文件包含了 12 个功能子模块的代码，不利于后期维护和代码导航。本变更将维护模块按功能子模块进行拆分，每个子模块单独使用文件管理。

## 维护模块功能子模块划分

维护模块当前包含以下 12 个功能子模块：

| 序号 | 子模块 | 功能描述 | Handler 方法数 |
|------|--------|----------|---------------|
| 1 | 重启/保存 (reboot_save) | 保存配置、重启设备、恢复出厂 | 3 |
| 2 | 系统配置 (system_config) | 网络配置、温度阈值、设备信息、时间日期 | 5 |
| 3 | 文件管理 (file_management) | 文件列表、上传、下载、删除、固件上传 | 5 |
| 4 | 日志管理 (log_management) | 日志列表、清除日志 | 2 |
| 5 | SNMP 配置 (snmp_config) | SNMP 配置、团体管理 | 4 |
| 6 | SNMP Trap (snmp_trap) | SNMP Trap 配置和目标管理 | 5 |
| 7 | 蠕虫防护 (worm_protection) | 蠕虫规则管理、攻击统计 | 5 |
| 8 | DDoS 防护 (ddos_protection) | DDoS 防护配置 | 2 |
| 9 | ARP 防护 (arp_protection) | ARP 防护配置 | 2 |
| 10 | 用户管理 (user_management) | 用户 CRUD | 4 |
| 11 | 会话管理 (session_management) | 会话列表、强制踢出 | 2 |
| 12 | 加载配置 (load_config) | 配置文件列表、加载配置 | 2 |

## 目标

1. 将 `maintenance_handler.go` 拆分为 12 个独立的子模块文件
2. 将 `maintenance.go` (datamodel) 拆分为 12 个独立的子模块文件
3. 保持原有 API 路由和功能不变
4. 提升代码可维护性和可读性
5. 便于后续功能扩展和问题定位

## 重构方案

### 目录结构

保持现有目录结构不变，在同一目录下按功能子模块拆分文件：

```
internal/handler/
├── maintenance_handler.go          # 主处理器（保留 NewMaintenanceHandler 和路由汇总）
├── maintenance_reboot_save.go      # 重启/保存
├── maintenance_system_config.go    # 系统配置
├── maintenance_file_management.go  # 文件管理
├── maintenance_log_management.go   # 日志管理
├── maintenance_snmp_config.go      # SNMP 配置
├── maintenance_snmp_trap.go        # SNMP Trap
├── maintenance_worm_protection.go  # 蠕虫防护
├── maintenance_ddos_protection.go  # DDoS 防护
├── maintenance_arp_protection.go   # ARP 防护
├── maintenance_user_management.go  # 用户管理
├── maintenance_session_management.go # 会话管理
└── maintenance_load_config.go      # 加载配置

internal/datamodel/
├── maintenance.go                  # 主页面（保留框架和通用组件）
├── maintenance_reboot_save.go      # 重启/保存页面
├── maintenance_system_config.go    # 系统配置页面
├── maintenance_file_management.go  # 文件管理页面
├── maintenance_log_management.go   # 日志管理页面
├── maintenance_snmp_config.go      # SNMP 配置页面
├── maintenance_snmp_trap.go        # SNMP Trap 页面
├── maintenance_worm_protection.go  # 蠕虫防护页面
├── maintenance_ddos_protection.go  # DDoS 防护页面
├── maintenance_arp_protection.go   # ARP 防护页面
├── maintenance_user_management.go  # 用户管理页面
├── maintenance_session_management.go # 会话管理页面
└── maintenance_load_config.go      # 加载配置页面
```

### Handler 层拆分规则

每个子模块文件包含：
1. 独立的结构体（可选，如需独立 DAO 依赖）
2. 对应的 HTTP Handler 方法
3. 辅助函数和验证逻辑

**示例** (`maintenance_reboot_save.go`)：
```go
package handler

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

// RebootSaveHandler 重启/保存子模块处理器
type RebootSaveHandler struct{}

// NewRebootSaveHandler 创建处理器
func NewRebootSaveHandler() *RebootSaveHandler {
    return &RebootSaveHandler{}
}

// SaveConfig 保存配置
func (h *RebootSaveHandler) SaveConfig(c *gin.Context) {
    // 实现...
}

// RebootSwitch 重启交换机
func (h *RebootSaveHandler) RebootSwitch(c *gin.Context) {
    // 实现...
}

// FactoryReset 恢复出厂
func (h *RebootSaveHandler) FactoryReset(c *gin.Context) {
    // 实现...
}
```

### DataModel 层拆分规则

每个子模块文件包含：
1. 独立的页面内容生成函数
2. UI 组件和 JavaScript 逻辑

**示例** (`maintenance_reboot_save.go`)：
```go
package datamodel

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/template/types"
)

// GetRebootSaveContent 重启/保存页面内容
func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    // 实现...
}
```

### 主处理器整合

`maintenance_handler.go` 保留：
1. `MaintenanceHandler` 主结构体（包含各子模块处理器引用）
2. `NewMaintenanceHandler()` 构造函数
3. 可选：统一的初始化方法

**示例**：
```go
package handler

// MaintenanceHandler 维护模块主处理器
type MaintenanceHandler struct {
    rebootSaveHandler      *RebootSaveHandler
    systemConfigHandler    *SystemConfigHandler
    fileManagementHandler  *FileManagementHandler
    // ... 其他子模块
}

// NewMaintenanceHandler 创建维护模块处理器
func NewMaintenanceHandler() *MaintenanceHandler {
    return &MaintenanceHandler{
        rebootSaveHandler:      NewRebootSaveHandler(),
        systemConfigHandler:    NewSystemConfigHandler(),
        fileManagementHandler:  NewFileManagementHandler(),
        // ... 初始化其他子模块
    }
}
```

## 任务列表

### 阶段一：Handler 层拆分（6 个任务）

1. **创建 Handler 子模块文件 (1-6)**
   - `maintenance_reboot_save.go`
   - `maintenance_system_config.go`
   - `maintenance_file_management.go`
   - `maintenance_log_management.go`
   - `maintenance_snmp_config.go`
   - `maintenance_snmp_trap.go`

2. **创建 Handler 子模块文件 (7-12)**
   - `maintenance_worm_protection.go`
   - `maintenance_ddos_protection.go`
   - `maintenance_arp_protection.go`
   - `maintenance_user_management.go`
   - `maintenance_session_management.go`
   - `maintenance_load_config.go`

3. **重构 maintenance_handler.go**
   - 移除已拆分的 Handler 方法
   - 添加子模块处理器结构体引用
   - 更新构造函数

4. **更新 main.go 路由注册**
   - 保持现有路由不变
   - 验证所有路由正确映射

5. **编译验证**
   - 执行 `go build` 验证编译
   - 修复编译错误

6. **功能测试**
   - 测试所有 API 端点
   - 验证功能正常

### 阶段二：DataModel 层拆分（6 个任务）

7. **创建 DataModel 子模块文件 (1-6)**
   - 拆分 UI 页面内容生成函数

8. **创建 DataModel 子模块文件 (7-12)**
   - 拆分剩余 UI 页面

9. **重构 maintenance.go**
   - 保留主页面框架
   - 移除已拆分的页面函数

10. **编译验证**
    - 执行 `go build` 验证
    - 修复编译错误

11. **UI 测试**
    - 访问所有维护模块页面
    - 验证 UI 正常显示

12. **代码审查与清理**
    - 检查是否有重复代码
    - 统一代码风格
    - 添加必要的注释

## 验收标准

1. 所有 24 个子模块文件创建完成
2. 编译无错误，运行正常
3. 所有 API 端点功能正常
4. 所有 UI 页面显示正常
5. 代码行数合理分布（每个文件 50-200 行）
6. 保持原有功能不变
7. 提交代码前通过 `go fmt` 和 `go vet`

## 预期收益

1. **可维护性提升**：每个文件专注于单一功能子模块
2. **代码导航简化**：快速定位特定功能的代码
3. **并行开发支持**：多人可同时修改不同子模块
4. **代码审查友好**：变更范围清晰
5. **问题定位加速**：错误日志可快速对应到具体文件

## 风险与注意事项

1. **路由映射**：确保路由注册不发生变化
2. **共享依赖**：注意 ConfigDAO 等共享依赖的处理
3. **循环导入**：避免子模块间产生循环导入
4. **代码重复**：识别并提取公共逻辑

## 时间估算

- 阶段一（Handler 层）：约 2-3 小时
- 阶段二（DataModel 层）：约 2-3 小时
- 测试与修复：约 1 小时
- 总计：约 5-7 小时
