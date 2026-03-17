# add-ping-cli-provider Tasks

## 实现任务列表

### 1. 定义 Provider 接口
- [ ] 创建 `service/provider/diagnostic_provider.go`
- [ ] 定义 `DiagnosticProvider` 接口
- [ ] 定义 `PingProvider` 接口（专注 Ping 功能）

### 2. 实现 Mock Provider
- [ ] 创建 `service/provider/mock_ping_provider.go`
- [ ] 保留现有随机数据生成逻辑
- [ ] 实现 `PingProvider` 接口

### 3. 实现 CLI Provider
- [ ] 创建 `service/provider/cli_ping_provider.go`
- [ ] 实现 `ExecutePing` 方法执行系统 Ping 命令
- [ ] 实现 Ping 输出解析器（提取 RTT、丢包率等）
- [ ] 支持 VRF 路由表选择（Windows: `ping -S`，Linux: `ip vrf exec`）
- [ ] 支持 Count/Timeout/Interval 参数传递

### 4. 扩展 ModeResolver
- [ ] 扩展 `service/mode/mode_resolver.go`
- [ ] 添加 `GetPingProvider()` 方法
- [ ] 根据配置返回 MockProvider 或 CLIProvider

### 5. 重构 Service 层
- [ ] 重构 `service/diagnostic_service.go`
- [ ] 移除硬编码的 Mock 逻辑
- [ ] 使用 ModeResolver 获取 Provider
- [ ] 调用 Provider 接口执行 Ping

### 6. 更新异步任务管理
- [ ] 确保任务结果存储正确
- [ ] 更新 `GetPingTaskResult` 方法
- [ ] 添加错误处理与日志记录

### 7. 测试验证
- [ ] 编译构建无错误
- [ ] Mock 模式测试
- [ ] CLI 模式测试
- [ ] 模式切换测试

## 依赖关系

```
Task 1 (接口定义)
    ↓
Task 2 (Mock Provider) ←────────┐
    ↓                           │
Task 3 (CLI Provider)           │
    ↓                           │
Task 4 (ModeResolver) ←─────────┘
    ↓
Task 5 (Service 重构)
    ↓
Task 6 (任务管理)
    ↓
Task 7 (测试验证)
```

## 并行工作

- Task 2 和 Task 3 可并行（都依赖 Task 1 的接口定义）
- Task 4 可与 Task 2/3 并行（只需接口定义，不依赖具体实现）
