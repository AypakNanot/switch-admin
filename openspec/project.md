# Project Context

## Purpose

**switch-admin** 是一个基于 Web 的智能交换机管理系统，为网络运维人员提供可视化的设备配置、监控、维护和诊断能力。

**核心目标**：
- 降低交换机运维门槛，可视化界面替代命令行操作
- 缩短故障 MTTR（平均修复时间），通过实时监控和诊断工具
- 提升配置效率，批量操作和配置校验
- 保障网络安全，权限分级和安全防护

**目标用户**：
- 网络运维工程师 - 日常配置、监控、故障排查
- 系统管理员 - 用户管理、系统维护、安全策略
- 网络规划工程师 - 拓扑规划、VLAN 划分、路由设计

## Tech Stack

**后端**：
- Go 1.25.5
- Gin v1.12.0 - Web 框架
- GoAdmin v1.2.26 - 管理后台框架
- Xorm v1.0.2 - ORM 框架
- SQLite3 - 嵌入式数据库
- Zap v1.19.1 - 日志库

**前端**：
- AdminLTE 主题
- GoAdmin Template Components
- jQuery
- Bootstrap

**开发工具**：
- GCC (用于 CGO/SQLite)
- Go Modules

## Project Conventions

### Code Style

**Go 代码规范**：
- 遵循 `gofmt` 格式化
- 使用 `go vet` 进行静态检查
- 函数命名：大驼峰（导出），小驼峰（内部）
- 错误处理：显式返回 error，不使用 `_` 忽略关键错误
- 包导入顺序：标准库 → 第三方库 → 内部包

**目录结构**：
```
switch-admin/
├── cmd/           # 应用入口
├── internal/      # 内部包（外部不可引用）
│   ├── datamodel/ # 数据模型和页面
│   ├── handler/   # HTTP 处理器
│   └── service/   # 业务逻辑
├── pkg/           # 公共包
├── data/          # 数据库文件
├── uploads/       # 上传文件
├── logs/          # 日志文件
└── PRD/           # 需求文档
```

### Architecture Patterns

**分层架构**：
```
路由层 (Gin) → 处理器层 (Handler) → 服务层 (Service) → 数据层 (Xorm/SQL)
                                    ↓
                              数据模型层 (datamodel/)
```

**GoAdmin 集成模式**：
- 使用 `context.Context` 处理 HTTP 请求
- 使用 `types.Panel` 返回页面内容
- 使用 `template/components` 构建 UI

**数据库模式**：
- SQLite3 嵌入式数据库
- 数据文件：`data/admin.db`
- 使用 Xorm ORM 进行数据访问

### Testing Strategy

**当前状态**：待完善

**目标策略**：
- 单元测试：核心业务逻辑（`*_test.go`）
- 集成测试：API 端点测试
- 手动测试：UI 界面功能验证

### Git Workflow

**分支策略**：
- `master` - 主分支，可部署状态
- 功能开发使用独立分支或 worktree

**提交规范**：
- 使用清晰的提交信息
- 功能提交使用英文或中文描述
- 关联 PRD 或 Issue 编号（如有）

## Domain Context

### 交换机管理域知识

**核心概念**：
| 术语 | 说明 |
|------|------|
| **端口 (Port)** | 交换机物理接口，如 GE1/0/1 |
| **VLAN** | 虚拟局域网，ID 范围 1-4094 |
| **PVID** | 端口默认 VLAN ID |
| **Trunk** | 干道模式，支持多 VLAN 传输 |
| **Access** | 接入模式，单 VLAN |
| **LACP** | 链路聚合控制协议 |
| **SNMP** | 简单网络管理协议 |

**端口状态**：
- 管理状态 (Admin Status): enable/disable
- 链路状态 (Link Status): up/down
- 速率：10M/100M/1000M/10G
- 双工：Full/Half

**流量统计**：
- 使用 64 位无符号整数（uint64）
- 10G 端口 3.4 秒可溢出 32 位计数器
- API 传输大数使用字符串格式

### 功能模块

**四大模块**：
1. **配置模块** - 端口、VLAN、链路聚合、风暴控制等
2. **监控模块** - Dashboard、端口状态、流量监控
3. **维护模块** - 用户管理、日志、SNMP、安全防护
4. **网络模块** - 静态路由、Ping、Traceroute、VCD

## Important Constraints

### 硬件约束（关键！）

| 约束 | 规格 | 工程影响 |
|------|------|----------|
| **内存** | 256MB-512MB | 禁止大数据量内存驻留，CSV 导出需流式处理 |
| **CPU** | 单核~双核 ARM/MIPS | 限制并发，防 DDoS 式刷新 |
| **存储** | 128MB-512MB Flash | 限制配置写入频率，防 Flash 磨损 |
| **带宽** | 10M/100M/1000M/10G | 流量统计需 64 位计数器 |

### 工程约束

**64 位计数器约束**：
```
⚠️ 所有流量统计字段必须使用 uint64_t
⚠️ API 传输超过 JS 安全整数的数值使用 String 类型
⚠️ 前端使用 BigInt 解析大数字符串
```

**统计清零偏移量机制**：
```
⚠️ 采用 Offset 机制，禁止直接清零硬件计数器
⚠️ Offset 存储在 /var/run/ 目录（重启清空）
⚠️ 展示值 = 当前硬件值 - Offset
```

**防 DDoS 刷新机制**：
```
⚠️ 前端使用 Visibility API，页面隐藏时暂停轮询
⚠️ 后端数据聚合缓存，API 响应≤100ms
⚠️ 禁止每个 HTTP 请求触发多次底层调用
```

### 安全约束

**权限控制**：
- super-admin: 超级管理员（唯一，不可删除）
- admin: 普通管理员（可管理下级）
- operator: 操作员（大部分配置权限）
- viewer: 只读用户

**高危操作防护**：
- 修改管理 VLAN → 红色警告 + 双重确认 + 30 秒倒计时
- 重启设备 → 二次确认
- 恢复出厂 → 输入 "CONFIRM" 确认
- 全局统计清零 → 输入 "CLEAR ALL" 确认

## External Dependencies

### Go 模块依赖

| 模块 | 用途 |
|------|------|
| github.com/GoAdminGroup/go-admin | 管理后台框架 |
| github.com/GoAdminGroup/themes | AdminLTE 主题 |
| github.com/gin-gonic/gin | Web 框架 |
| github.com/mattn/go-sqlite3 | SQLite3 驱动 (CGO) |
| xorm.io/xorm | ORM 框架 |
| go.uber.org/zap | 日志库 |

### 系统依赖

| 依赖 | 用途 |
|------|------|
| GCC | CGO 编译（SQLite3 驱动） |
| SQLite3 | 数据库引擎 |

### 外部系统

| 系统 | 集成方式 | 说明 |
|------|----------|------|
| 交换机底层驱动 | 系统调用/CLI | 获取端口状态、配置参数 |
| SNMP | SNMP v2/v3 | 可选的监控协议 |

## OpenSpec Usage

**创建变更提案**：
```bash
# 创建新的变更
openspec change add <change-id>

# 查看变更列表
openspec list

# 查看规格列表
openspec list --specs
```

**变更 ID 命名规范**：
- 使用动词前缀：`add-`, `update-`, `remove-`, `refactor-`
- 使用 kebab-case：`add-vlan-config-page`
- 保持唯一性和描述性

**变更提案结构**：
```
openspec/changes/<change-id>/
├── proposal.md    # 变更提案
├── tasks.md       # 实现任务清单
└── design.md      # 设计文档（可选）
```
