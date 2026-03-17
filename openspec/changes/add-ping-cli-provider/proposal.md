# add-ping-cli-provider

实现 Ping 功能的分层架构与 CLI 模式实现

## 概述

当前 Ping 功能完全基于 Mock 数据（`rand.Float64() < 0.8`），无法执行真实的 Ping 诊断。本变更实现分层架构模式，支持多种数据源（Mock、CLI、NETCONF 等），并具体实现 CLI 模式的真实 Ping 功能。

## 目标

1. 建立分层架构模式：Page → Handler → Service → ModeResolver → Provider Interface → Implementations
2. 定义统一的 DiagnosticProvider 接口
3. 实现 CLI Provider 执行真实 Ping 命令
4. 保留 Mock Provider 用于测试
5. 支持通过模式切换选择不同 Provider

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│ Page Layer (datamodel/ping.go)                              │
│   - HTML/JS 表单 UI                                           │
│   - 前端展示逻辑                                              │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Handler Layer (handler/diagnostic_handler.go)               │
│   - HTTP 请求处理                                             │
│   - 参数验证与绑定                                            │
│   - JSON 响应格式化                                           │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Service Layer (service/diagnostic_service.go)               │
│   - 业务逻辑编排                                              │
│   - 任务管理（创建/查询/删除）                                │
│   - 通过 ModeResolver 选择 Provider                           │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ ModeResolver (service/mode/mode_resolver.go)                │
│   - 读取配置（mock vs switch 模式）                           │
│   - 返回对应的 Provider 实例                                  │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Provider Interface (service/provider/diagnostic_provider.go)│
│   - DiagnosticProvider 接口定义                             │
│   - 统一方法签名                                              │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Provider Implementations                                    │
│   - MockProvider: 生成模拟数据（测试用）                      │
│   - CLIProvider: 执行系统 Ping 命令（真实功能）                │
│   - NETCONFProvider: NETCONF 协议（未来扩展）                 │
│   - DBProvider: 数据库缓存（未来扩展）                        │
└─────────────────────────────────────────────────────────────┘
```

## 关键变更

### 1. Provider 接口定义

```go
type DiagnosticProvider interface {
    // ExecutePing 执行 Ping 诊断
    ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error)
}
```

### 2. CLI Provider 实现

- 使用 `exec.Command` 执行系统 Ping 命令
- 解析 Ping 输出提取统计信息
- 支持 VRF 路由表选择
- 支持 Count/Timeout/Interval 参数

### 3. ModeResolver 扩展

- 新增 `GetDiagnosticProvider()` 方法
- 根据配置返回 MockProvider 或 CLIProvider

### 4. Service 层重构

- 移除硬编码的 Mock 逻辑
- 通过 ModeResolver 获取 Provider
- 调用 Provider 接口执行实际操作

## 任务列表

1. 创建 `service/provider/diagnostic_provider.go` 定义接口
2. 创建 `service/provider/cli_ping_provider.go` 实现 CLI Ping
3. 扩展 `service/mode/mode_resolver.go` 添加 Provider 解析
4. 重构 `service/diagnostic_service.go` 使用 Provider 模式
5. 更新 Handler 层支持异步任务结果查询

## 验收标准

1. CLI 模式下 Ping 功能执行真实系统 Ping 命令
2. Mock 模式下 Ping 功能生成模拟数据
3. 可通过配置切换模式
4. 代码符合项目现有风格
5. 编译运行无错误

## 实现状态

**✅ 已完成**

1. ✅ Provider 接口定义 - `internal/service/provider/interface.go`
2. ✅ CLI Provider 实现 - `internal/service/provider/cli/ping.go`
3. ✅ Mock Provider 实现 - `internal/service/provider/mock/ping.go`
4. ✅ ModeResolver 扩展 - 添加 `GetPingProvider()` 方法
5. ✅ Service 层重构 - 使用 Provider 模式
6. ✅ 目录结构优化 - 按 Provider 类型分目录 (cli/, mock/, netconf/)
7. ✅ 测试验证 - CLI 模式 Ping 127.0.0.1 返回真实 RTT 数据 (10-41ms)

### 测试结果

```bash
# CLI 模式测试
$ curl -s -X POST http://localhost:9033/api/v1/diagnostic/ping \
  -H "Content-Type: application/json" \
  -d '{"target":"127.0.0.1","count":4,"timeout":2}'

{
  "code": 200,
  "data": {
    "task_id": "ping_1773762227143379500",
    "status": "completed",
    "target": "127.0.0.1",
    "results": [
      {"seq": 1, "time": "26.03ms", "ttl": 64, "status": "success"},
      {"seq": 2, "time": "10.98ms", "ttl": 64, "status": "success"},
      {"seq": 3, "time": "41.44ms", "ttl": 64, "status": "success"},
      {"seq": 4, "time": "41.60ms", "ttl": 64, "status": "success"}
    ],
    "statistics": {
      "sent": 4,
      "received": 4,
      "loss_rate": "0%",
      "min_time": "10.98ms",
      "avg_time": "30.01ms",
      "max_time": "41.60ms"
    }
  }
}
```

## 未来扩展

- NETCONF 模式：通过 NETCONF 协议远程执行 Ping
- DB 模式：从数据库读取缓存的 Ping 结果
- SNMP 模式：通过 SNMP 执行 Ping 诊断
