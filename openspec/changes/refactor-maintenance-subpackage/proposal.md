# Proposal: 维护模块按功能分组重构

## 概述

当前维护模块的 Handler 和 DataModel 层文件分散在包目录下，缺乏清晰的组织。本提案建议按功能分组到不同的 `.go` 文件中，所有文件共享同一个包名 (`maintenance`)，提升代码组织性。

## 目标

1. **提升可维护性**：每个功能模块独立文件，职责清晰
2. **简化包结构**：不使用子包，所有功能在同一包内
3. **便于扩展**：新增功能模块只需添加新文件
4. **保持兼容性**：API 路由和功能行为保持不变

## 非目标

- 不修改现有功能逻辑
- 不改变 API 路由映射
- 不调整 UI 页面结构

## 设计方案

### 核心原则

1. **一个模块一个包** - 维护模块整体是一个包 (`maintenance`)
2. **按功能分文件** - 每个功能子模块是一个独立的 `.go` 文件
3. **统一命名** - 文件名与功能对应

### Handler 层结构

```
internal/handler/maintenance/
├── handler.go              # 主处理器（Handler 结构体 + New 函数）
├── reboot_save.go          # 重启/保存 Handler
├── system_config.go        # 系统配置 Handler
├── file_management.go      # 文件管理 Handler
├── log_management.go       # 日志管理 Handler
├── snmp_config.go          # SNMP 配置 Handler
├── snmp_trap.go            # SNMP Trap Handler
├── worm_protection.go      # 蠕虫防护 Handler
├── ddos_protection.go      # DDoS 防护 Handler
├── arp_protection.go       # ARP 防护 Handler
├── user_management.go      # 用户管理 Handler
├── session_management.go   # 会话管理 Handler
└── load_config.go          # 加载配置 Handler
```

### DataModel 层结构

```
internal/datamodel/maintenance/
├── maintenance.go          # 导出页面函数（保持兼容）
├── reboot_save.go          # 重启/保存页面
├── system_config.go        # 系统配置页面
├── file_management.go      # 文件管理页面
├── log_management.go       # 日志管理页面
├── snmp_config.go          # SNMP 配置页面
├── snmp_trap.go            # SNMP Trap 页面
├── worm_protection.go      # 蠕虫防护页面
├── ddos_protection.go      # DDoS 防护页面
├── arp_protection.go       # ARP 防护页面
├── user_management.go      # 用户管理页面
├── session_management.go   # 会话管理页面
└── load_config.go          # 加载配置页面
```

### 代码组织模式

#### Handler 层
- 所有文件共享同一个包名 `package maintenance`
- `Handler` 结构体定义在 `handler.go` 中
- 其他文件中的方法直接添加到 `Handler` 结构体上

```go
// handler.go
package maintenance

type Handler struct{}

func New() *Handler {
    return &Handler{}
}

// reboot_save.go
package maintenance

func (h *Handler) SaveConfig(c *gin.Context) {
    // 实现
}
```

#### DataModel 层
- 所有文件共享同一个包名 `package maintenance`
- 导出函数（大写开头）在 `maintenance.go` 中定义
- 内部实现函数（小写开头）在各功能文件中定义
- 导出函数委托调用内部实现函数

```go
// maintenance.go
package maintenance

func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    return getRebootSaveContent(ctx)
}

// reboot_save.go
package maintenance

func getRebootSaveContent(ctx *context.Context) (types.Panel, error) {
    // 实现
}
```

## 实施步骤

### 阶段一：Handler 层按功能分组
1. 创建 `internal/handler/maintenance/` 目录
2. 创建 handler.go（主处理器）
3. 按功能创建 11 个文件
4. 编译验证

### 阶段二：DataModel 层按功能分组
1. 创建 `internal/datamodel/maintenance/` 目录
2. 创建 maintenance.go（导出函数）
3. 按功能创建 12 个文件
4. 编译验证

### 阶段三：测试验证
1. 全量编译
2. 代码格式化
3. 功能测试

## 验收标准

1. [ ] 所有功能文件创建完成
2. [ ] `go build` 编译无错误
3. [ ] `go vet ./...` 检查通过
4. [ ] 所有 API 端点功能正常
5. [ ] 所有 UI 页面显示正常
6. [ ] 代码通过 `go fmt` 格式化

## 实施状态

- [x] Handler 层文件创建完成（11 个文件）
- [x] DataModel 层文件创建完成（13 个文件）
- [x] 编译验证通过
- [x] 代码格式化完成

## 优势

1. **清晰的职责分离** - 每个功能模块独立文件，易于定位和维护
2. **统一的包管理** - 所有功能在同一个包内，不需要子包导入
3. **简洁的调用链** - main.go 中的调用方式保持不变
4. **易于扩展** - 新增功能只需添加新文件
