# Diagnostic Provider Spec

## 概述

定义诊断功能（Ping）的统一 Provider 接口，支持多种实现方式（Mock、CLI 等）。

## 能力范围

- 定义 PingProvider 接口
- 规范 Provider 实现标准
- 支持模式切换

---

## ADDED Requirements

### Requirement: PingProvider 接口定义

系统**MUST**定义统一的 PingProvider 接口，规范所有 Ping 诊断功能的行为。

```go
// PingProvider Ping 诊断接口
type PingProvider interface {
    // ExecutePing 执行 Ping 诊断
    // - ctx: 上下文，支持取消操作
    // - req: Ping 请求参数（目标 IP、VRF、Count、Timeout 等）
    // - 返回：Ping 结果（包含统计信息）和错误
    ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error)
}
```

#### Scenario: 接口一致性
当 Service 层调用 DiagnosticProvider 时，无论底层是 Mock 实现还是 CLI 实现，返回的数据结构**MUST**完全一致，上层无需关心具体实现。

#### Scenario: 实现 Provider
开发者实现新的 Ping Provider 时，只需实现 `ExecutePing` 方法，无需关心任务管理、HTTP 处理等上层逻辑。

---

### Requirement: MockPingProvider 实现

系统**MUST**保留 Mock 实现用于测试和开发环境。

```go
// MockPingProvider Mock 模式的 Ping Provider
type MockPingProvider struct{}

// ExecutePing 生成模拟 Ping 结果
func (p *MockPingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error)
```

#### Scenario: 开发测试
开发人员在本地测试时，使用 Mock 模式无需连接真实设备，快速验证 UI 和业务流程。

---

### Requirement: CLIPingProvider 实现

系统**MUST**实现 CLI Provider 执行真实系统 Ping 命令。

```go
// CLIPingProvider CLI 模式的 Ping Provider
type CLIPingProvider struct {
    execFunc func(command string, args ...string) ([]byte, error)
}

// ExecutePing 执行系统 Ping 命令并解析结果
func (p *CLIPingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error)
```

#### Scenario: 真实环境
生产环境使用 CLI 模式，通过执行系统 `ping` 命令获取真实网络诊断结果。

#### Scenario: 跨平台支持
CLIPingProvider **MUST**支持 Windows 和 Linux 系统，自动适配不同的 Ping 命令参数和输出格式。

---

### Requirement: 命令参数构建

CLIPingProvider **MUST**正确构建 Ping 命令参数。

| 参数 | Windows | Linux |
|------|---------|-------|
| Count | `-n <count>` | `-c <count>` |
| Timeout | `-w <ms>` | `-W <seconds>` |
| Interval | 不支持 | `-i <seconds>` |
| Target | 直接追加 | 直接追加 |

#### Scenario: Count 参数
用户请求 Ping 5 次，Windows **MUST**构建 `ping -n 5 <target>`，Linux **MUST**构建 `ping -c 5 <target>`。

#### Scenario: Timeout 参数
用户设置超时 2 秒，Windows **MUST**构建 `-w 2000`（毫秒），Linux **MUST**构建 `-W 2`（秒）。

---

### Requirement: Ping 输出解析

CLIPingProvider **MUST**解析系统 Ping 命令输出，提取关键信息。

**Windows 输出示例：**
```
来自 192.168.1.1 的回复：字节=32 时间=1ms TTL=64
来自 192.168.1.1 的回复：字节=32 时间=2ms TTL=64
请求超时。
```

**Linux 输出示例：**
```
64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=1.23 ms
64 bytes from 192.168.1.1: icmp_seq=2 ttl=64 time=2.45 ms
Request timeout for icmp_seq 3
```

**提取信息：**
- 序列号（Seq）
- 是否成功（Success）
- RTT 时间（Time，单位毫秒）
- 统计信息（发送/接收/丢失/丢包率/最小/最大/平均 RTT）

#### Scenario: 解析成功响应
解析 `来自 192.168.1.1 的回复：字节=32 时间=1ms TTL=64` **MUST**得到 `PingSeqResult{Seq: 1, Success: true, Time: 1}`。

#### Scenario: 解析超时响应
解析 `请求超时` **MUST**得到 `PingSeqResult{Seq: 3, Success: false, Time: 0}`。

---

### Requirement: ModeResolver 扩展

ModeResolver **MUST**支持返回 Ping Provider。

```go
// ModeResolver 模式解析器
type ModeResolver interface {
    // GetPingProvider 返回当前模式对应的 Ping Provider
    GetPingProvider() provider.PingProvider

    // GetCurrentMode 返回当前运行模式
    GetCurrentMode() string
}
```

#### Scenario: 模式切换
用户配置为 `switch` 模式时，`GetPingProvider()` **MUST**返回 `CLIPingProvider`；配置为 `mock` 模式时**MUST**返回 `MockPingProvider`。

---

### Requirement: Service 层重构

DiagnosticService **MUST**使用 Provider 模式，移除硬编码的 Mock 逻辑。

```go
type DiagnosticService struct {
    pingTaskResults sync.Map
    modeResolver    *mode.ModeResolver
}

func (s *DiagnosticService) ExecutePing(taskID string, req model.PingRequest) error {
    provider := s.modeResolver.GetPingProvider()
    result, err := provider.ExecutePing(context.Background(), req)
    if err != nil {
        return err
    }
    result.TaskID = taskID
    s.pingTaskResults.Store(taskID, result)
    return nil
}
```

#### Scenario: 执行 Ping
用户通过 UI 发起 Ping 请求，Handler 调用 `ExecutePing`，Service **MUST**通过 ModeResolver 获取 Provider，Provider 执行实际操作。

---

## MODIFIED Requirements

### Requirement: 模式配置

系统配置中的 Mode 字段现在 MUST 决定 ModeResolver 返回的 Provider 类型。

**变更说明：**

- **之前：** 模式配置仅用于条件判断
- **现在：** 模式配置 MUST 决定 ModeResolver 返回的 Provider 类型

```go
// config.Config 新增字段
type Config struct {
    Mode string `json:"mode"` // "mock" 或 "switch"
}
```

#### Scenario: 配置生效
系统启动时读取配置文件，ModeResolver 根据 Mode 字段返回对应的 Provider。

---

## 验收标准

1. [ ] 编译构建无错误
2. [ ] Mock 模式：Ping 生成模拟数据
3. [ ] CLI 模式：Ping 执行真实系统命令
4. [ ] 模式切换：修改配置后立即生效
5. [ ] 跨平台：Windows 和 Linux 均能正常工作
6. [ ] 错误处理：网络不可达时正确返回错误信息
