# Code Organization Spec - refactor-maintenance-subpackage

## 概述

维护模块分包重构，将 Handler 和 DataModel 层的子模块文件组织到独立的子包中，提升代码可维护性和扩展性。

---

## MODIFIED Requirements

### Requirement: Handler 层包组织

Handler 子模块文件**SHALL**移动到对应的子包目录中，每个子包**SHALL**只包含一个处理器文件。

**目录结构 SHALL 遵循以下规范：**

```
internal/handler/maintenance/
├── handler.go                    # 主处理器
├── reboot_save/handler.go        # 重启/保存子包
├── system_config/handler.go      # 系统配置子包
├── file_management/handler.go    # 文件管理子包
├── log_management/handler.go     # 日志管理子包
├── snmp_config/handler.go        # SNMP 配置子包
├── snmp_trap/handler.go          # SNMP Trap 子包
├── worm_protection/handler.go    # 蠕虫防护子包
├── ddos_protection/handler.go    # DDoS 防护子包
├── arp_protection/handler.go     # ARP 防护子包
├── user_management/handler.go    # 用户管理子包
├── session_management/handler.go # 会话管理子包
└── load_config/handler.go        # 加载配置子包
```

#### Scenario: 导航定位

开发者需要修改 SNMP 配置处理器时，**MUST**直接访问路径 `internal/handler/maintenance/snmp_config/handler.go`，无需在文件列表中滚动查找。

#### Scenario: 新增子模块

添加新的维护功能（如"批量配置"）时，**MUST**创建新子包 `internal/handler/maintenance/batch_config/handler.go`，不影响现有结构。

---

### Requirement: Handler 子包命名规范

每个 Handler 子包**SHALL**遵循统一的命名规范：包名**SHALL**与目录名一致，文件名**SHALL**统一使用 `handler.go`，构造函数**SHALL**统一使用 `New()` 返回 `*Handler`。

#### Scenario: 开发者理解代码

新成员查看代码时，**MUST**通过统一的命名快速理解包结构和职责。

#### Scenario: IDE 自动补全

使用 IDE 时，**MUST**能通过路径补全快速找到目标处理器。

---

### Requirement: DataModel 层包组织

DataModel 子模块文件**SHALL**移动到对应的子包目录中，每个子包**SHALL**只包含一个页面文件。

**目录结构 SHALL 遵循以下规范：**

```
internal/datamodel/maintenance/
├── page.go                       # 主页面框架（可选）
├── reboot_save/page.go           # 重启/保存子包
├── system_config/page.go         # 系统配置子包
├── file_management/page.go       # 文件管理子包
├── log_management/page.go        # 日志管理子包
├── snmp_config/page.go           # SNMP 配置子包
├── snmp_trap/page.go             # SNMP Trap 子包
├── worm_protection/page.go       # 蠕虫防护子包
├── ddos_protection/page.go       # DDoS 防护子包
├── arp_protection/page.go        # ARP 防护子包
├── user_management/page.go       # 用户管理子包
├── session_management/page.go    # 会话管理子包
└── load_config/page.go           # 加载配置子包
```

#### Scenario: UI 页面修改

修改 SNMP 配置页面时，**MUST**直接访问 `internal/datamodel/maintenance/snmp_config/page.go`，不影响其他页面。

---

### Requirement: DataModel 子包命名规范

每个 DataModel 子包**SHALL**遵循统一的命名规范：包名**SHALL**与目录名一致，文件名**SHALL**统一使用 `page.go`，页面函数**SHALL**统一使用 `GetContent()` 返回 `(types.Panel, error)`。

#### Scenario: 页面函数调用

路由注册时，**MUST**通过 `reboot_save.GetContent(ctx)` 调用页面函数。

---

### Requirement: 主处理器依赖注入

主处理器**SHALL**持有所有子包处理器的引用，并通过构造函数统一初始化。子处理器实例**SHALL**通过主处理器统一管理。

#### Scenario: 子处理器复用

所有子处理器通过主处理器统一管理，**MUST**避免重复创建实例。

---

### Requirement: 路由导入路径更新

路由注册文件（`cmd/main.go`）的导入路径**SHALL**更新为新的包路径，所有引用**SHALL**指向新的子包位置。

#### Scenario: 路由功能正常

前端调用 API 时，**MUST**正常路由到对应的子包处理器方法。

---

## RENAMED Requirements

### Requirement: 文件命名

原 `maintenance_*.go` 文件**SHALL**重命名为子包中的标准文件名：

| 原文件名 | 新路径 |
|---------|--------|
| `maintenance_reboot_save.go` | `maintenance/reboot_save/handler.go` |
| `maintenance_system_config.go` | `maintenance/system_config/handler.go` |
| `maintenance_file_management.go` | `maintenance/file_management/handler.go` |
| `maintenance_log_management.go` | `maintenance/log_management/handler.go` |
| `maintenance_snmp_config.go` | `maintenance/snmp_config/handler.go` |
| `maintenance_snmp_trap.go` | `maintenance/snmp_trap/handler.go` |
| `maintenance_worm_protection.go` | `maintenance/worm_protection/handler.go` |
| `maintenance_ddos_protection.go` | `maintenance/ddos_protection/handler.go` |
| `maintenance_arp_protection.go` | `maintenance/arp_protection/handler.go` |
| `maintenance_user_management.go` | `maintenance/user_management/handler.go` |
| `maintenance_session_management.go` | `maintenance/session_management/handler.go` |
| `maintenance_load_config.go` | `maintenance/load_config/handler.go` |

---

## 验收标准

1. [ ] Handler 层 12 个子包创建完成
2. [ ] DataModel 层 12 个子包创建完成
3. [ ] `go build` 编译无错误
4. [ ] `go vet ./...` 检查通过
5. [ ] 所有 API 端点功能正常
6. [ ] 所有 UI 页面显示正常
7. [ ] 代码通过 `go fmt` 格式化
8. [ ] 无循环导入

---

## 非目标

本变更**不涉及**：
- API 路由变更
- 功能逻辑修改
- UI 样式调整
- 数据库结构变更
