# Ping 分层架构设计

## 1. 架构模式

采用经典的三层架构 + Provider 模式：

```
Handler (HTTP 层) → Service (业务逻辑) → Provider (数据源)
                         ↓
                    ModeResolver (模式选择)
```

## 2. 核心接口设计

### 2.1 Provider 接口

```go
// PingProvider 定义 Ping 操作的统一接口
type PingProvider interface {
    // ExecutePing 执行 Ping 诊断
    // req: Ping 请求参数（目标 IP、VRF、Count、Timeout 等）
    // 返回：Ping 结果（包含统计信息）和错误
    ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error)
}
```

### 2.2 ModeResolver 接口

```go
// ModeResolver 模式解析器
type ModeResolver interface {
    // GetPingProvider 返回当前模式对应的 Ping Provider
    GetPingProvider() provider.PingProvider
    // GetCurrentMode 返回当前运行模式
    GetCurrentMode() string
}
```

## 3. 实现细节

### 3.1 Mock Provider

```go
type MockPingProvider struct{}

func (p *MockPingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error) {
    // 模拟网络延迟
    time.Sleep(time.Duration(req.Count) * time.Second)

    // 生成随机结果（80% 成功率）
    results := make([]model.PingSeqResult, 0, req.Count)
    for i := 0; i < req.Count; i++ {
        success := rand.Float64() < 0.8
        // ... 生成模拟数据
    }

    return &model.PingResult{
        TaskID:     taskID,
        Target:     req.Target,
        Status:     "completed",
        SeqResults: results,
        Statistics: statistics,
    }, nil
}
```

### 3.2 CLI Provider

```go
type CLIPingProvider struct {
    // 可注入的 Command 执行器（便于测试）
    execFunc func(command string, args ...string) ([]byte, error)
}

func (p *CLIPingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingResult, error) {
    // 1. 构建 Ping 命令
    cmd := "ping"
    args := p.buildPingArgs(req)

    // 2. 处理 VRF 路由表
    if req.VrfID != "" {
        // Windows: ping -S <source>
        // Linux: ip vrf exec <vrf> ping ...
        args = p.wrapVRFCommand(req.VrfID, args)
    }

    // 3. 执行命令
    output, err := p.execFunc(cmd, args...)
    if err != nil {
        return nil, fmt.Errorf("ping 命令执行失败：%w", err)
    }

    // 4. 解析输出
    return p.parsePingOutput(output, req)
}

// buildPingArgs 构建 Ping 命令参数
func (p *CLIPingProvider) buildPingArgs(req model.PingRequest) []string {
    args := []string{}

    // Count
    if runtime.GOOS == "windows" {
        args = append(args, "-n", strconv.Itoa(req.Count))
    } else {
        args = append(args, "-c", strconv.Itoa(req.Count))
    }

    // Timeout
    if runtime.GOOS == "windows" {
        args = append(args, "-w", strconv.Itoa(req.Timeout*1000)) // Windows 单位毫秒
    } else {
        args = append(args, "-W", strconv.Itoa(req.Timeout))
    }

    // Interval
    if runtime.GOOS != "windows" {
        args = append(args, "-i", strconv.Itoa(req.Interval))
    }

    // Target
    args = append(args, req.Target)

    return args
}
```

### 3.3 Ping 输出解析器

```go
// parsePingOutput 解析 Ping 命令输出
func (p *CLIPingProvider) parsePingOutput(output []byte, req model.PingRequest) (*model.PingResult, error) {
    result := &model.PingResult{
        Target:     req.Target,
        Status:     "completed",
        SeqResults: make([]model.PingSeqResult, 0),
        Statistics: model.PingStatistics{},
    }

    // 解析每一行
    lines := strings.Split(string(output), "\n")
    var rttTimes []float64

    for _, line := range lines {
        // Windows: 来自 192.168.1.1 的回复：字节=32 时间=1ms TTL=64
        // Linux: 64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=1.23 ms

        if runtime.GOOS == "windows" {
            // 解析 Windows 格式
            if strings.Contains(line, "来自") && strings.Contains(line, "回复") {
                // 提取时间
                rtt := p.extractWindowsRTT(line)
                result.SeqResults = append(result.SeqResults, model.PingSeqResult{
                    Seq:    len(result.SeqResults) + 1,
                    Success: true,
                    Time:   rtt,
                })
                rttTimes = append(rttTimes, rtt)
            } else if strings.Contains(line, "请求超时") {
                result.SeqResults = append(result.SeqResults, model.PingSeqResult{
                    Seq:     len(result.SeqResults) + 1,
                    Success: false,
                    Time:    0,
                })
            }
        } else {
            // 解析 Linux 格式
            if strings.Contains(line, "bytes from") {
                rtt := p.extractLinuxRTT(line)
                result.SeqResults = append(result.SeqResults, model.PingSeqResult{
                    Seq:    len(result.SeqResults) + 1,
                    Success: true,
                    Time:   rtt,
                })
                rttTimes = append(rttTimes, rtt)
            } else if strings.Contains(line, "timeout") {
                result.SeqResults = append(result.SeqResults, model.PingSeqResult{
                    Seq:     len(result.SeqResults) + 1,
                    Success: false,
                    Time:    0,
                })
            }
        }
    }

    // 计算统计信息
    var successCount int
    var totalRtt float64
    for _, r := range result.SeqResults {
        if r.Success {
            successCount++
            totalRtt += r.Time
        }
    }

    result.Statistics = model.PingStatistics{
        Sent:        req.Count,
        Received:    successCount,
        Lost:        req.Count - successCount,
        LossPercent: float64(req.Count-successCount) / float64(req.Count) * 100,
        MinRTT:      min(rttTimes),
        MaxRTT:      max(rttTimes),
        AvgRTT:      totalRtt / float64(successCount),
    }

    return result, nil
}
```

### 3.4 ModeResolver 实现

```go
type ModeResolver struct {
    currentMode   string
    mockProvider  *provider.MockPingProvider
    cliProvider   *provider.CLIPingProvider
}

func NewModeResolver() *ModeResolver {
    return &ModeResolver{
        currentMode:  getConfig().Mode,
        mockProvider: &provider.MockPingProvider{},
        cliProvider:  &provider.CLIPingProvider{},
    }
}

func (mr *ModeResolver) GetPingProvider() provider.PingProvider {
    switch mr.currentMode {
    case "switch", "cli":
        return mr.cliProvider
    default:
        return mr.mockProvider
    }
}

func (mr *ModeResolver) GetCurrentMode() string {
    return mr.currentMode
}
```

### 3.5 Service 层重构

```go
type DiagnosticService struct {
    pingTaskResults sync.Map // map[string]*model.PingResult
    modeResolver    *mode.ModeResolver
}

func NewDiagnosticService() *DiagnosticService {
    return &DiagnosticService{
        modeResolver: mode.NewModeResolver(),
    }
}

func (s *DiagnosticService) ExecutePing(taskID string, req model.PingRequest) error {
    // 通过 ModeResolver 获取 Provider
    provider := s.modeResolver.GetPingProvider()

    // 执行 Ping
    result, err := provider.ExecutePing(context.Background(), req)
    if err != nil {
        return err
    }

    // 存储结果
    result.TaskID = taskID
    s.pingTaskResults.Store(taskID, result)

    return nil
}
```

## 4. 流程图

```
用户请求 Ping
    ↓
Handler: CreatePingTask()
    ↓
Service: ExecutePing(taskID, req)
    ↓
ModeResolver: GetPingProvider()
    ↓
┌─────────────────────────────────────┐
│           模式判断                   │
├─────────────────────────────────────┤
│ Mock 模式 → MockProvider.ExecutePing │
│ CLI 模式  → CLIPingProvider.ExecutePing│
└─────────────────────────────────────┘
    ↓
返回 PingResult
    ↓
Service: 存储结果到 sync.Map
    ↓
Handler: 返回 task_id
    ↓
前端轮询：GetPingTaskResult(task_id)
    ↓
Service: 从 sync.Map 读取结果
    ↓
返回完整 Ping 结果
```

## 5. 关键设计决策

### 5.1 为什么使用 Provider 模式？

- **可测试性**: Mock Provider 便于单元测试
- **可扩展性**: 轻松添加 NETCONF、SNMP 等新实现
- **隔离性**: 各 Provider 实现相互独立
- **可切换性**: 通过 ModeResolver 动态切换

### 5.2 异步任务设计

- 使用 `sync.Map` 存储任务结果
- 前端轮询方式获取结果
- 支持多个并发 Ping 任务

### 5.3 跨平台处理

- Windows/Linux Ping 命令参数不同
- 输出格式解析不同
- 使用 `runtime.GOOS` 判断平台

## 6. 错误处理

```go
// CLI Provider 错误处理
if err != nil {
    if exitErr, ok := err.(*exec.ExitError); ok {
        // 命令执行失败（如网络不可达）
        return &model.PingResult{
            Status: "failed",
            Error:  fmt.Sprintf("网络不可达：%v", exitErr),
        }, nil
    }
    return nil, fmt.Errorf("ping 命令执行失败：%w", err)
}
```

## 7. 测试策略

### 7.1 单元测试

```go
func TestCLIPingProvider_ExecutePing(t *testing.T) {
    // Mock exec 函数
    provider := &CLIPingProvider{
        execFunc: func(cmd string, args ...string) ([]byte, error) {
            return []byte(mockPingOutput), nil
        },
    }

    result, err := provider.ExecutePing(ctx, req)
    // 断言结果
}
```

### 7.2 集成测试

- 真实执行 Ping localhost
- 验证结果解析正确性

## 8. 最终目录结构

```
internal/service/provider/
├── interface.go              # Provider 接口定义
│   ├── PingProvider
│   ├── TracerouteProvider
│   └── CableTestProvider
├── cli/                      # CLI Provider 实现
│   └── ping.go               # 执行系统 ping 命令
├── mock/                     # Mock Provider 实现
│   └── ping.go               # 生成模拟数据
└── netconf/                  # NETCONF Provider (未来扩展)
```

### 8.1 命名简化

重构后移除了文件名中的 Provider 类型前缀：
- `cli_ping_provider.go` → `cli/ping.go` (package `cli`)
- `mock_ping_provider.go` → `mock/ping.go` (package `mock`)
- `CLIPingProvider` → `PingProvider` (在 cli 包中)
- `MockPingProvider` → `PingProvider` (在 mock 包中)

通过目录隔离，不同 Provider 可以使用相同的结构体名称，通过包名区分：
- `cli.PingProvider` - CLI 模式的 Ping Provider
- `mock.PingProvider` - Mock 模式的 Ping Provider
