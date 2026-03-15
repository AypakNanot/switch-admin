# 架构设计文档

**变更 ID**: `dual-mode-adapter-architecture`
**版本**: 1.0
**创建日期**: 2026-03-15

---

## 1. 架构决策

### 1.1 核心设计原则

#### 1.1.1 接口隔离原则

```
┌─────────────────────────────────────────────────────────────────┐
│                      业务逻辑层 (Service)                        │
│                              │                                   │
│                              │ 依赖接口                           │
│                              ▼                                   │
│                    ┌───────────────────┐                         │
│                    │   DataProvider    │  ← 业务层不关心实现     │
│                    └───────────────────┘                         │
│                              │                                   │
│              ┌───────────────┴───────────────┐                   │
│              ▼                               ▼                   │
│    ┌───────────────────┐           ┌───────────────────┐         │
│    │  MockProvider     │           │  SwitchAdapter    │         │
│    │  (数据库模拟)      │           │  (交换机硬件)      │         │
│    └───────────────────┘           └───────────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

**决策**: 业务逻辑层只依赖 `DataProvider` 接口，不直接依赖具体实现。

**理由**:
- 便于单元测试（可注入 Mock）
- 便于扩展新的数据源
- 降低耦合度

#### 1.1.2 策略模式

```go
// 模式判断使用策略模式
type ModeStrategy interface {
    GetProvider() DataProvider
    GetAdapter(functionName string) SwitchAdapter
}

type MockModeStrategy struct {
    provider *MockProvider
}

type SwitchModeStrategy struct {
    adapters map[string]SwitchAdapter
}
```

**决策**: 使用策略模式处理不同模式的实现差异。

**理由**:
- 避免大量的 `if mode == "mock"` 判断
- 新增模式时无需修改现有代码
- 各模式实现相互隔离

#### 1.1.3 适配器模式

```
┌─────────────────────────────────────────────────────────────────┐
│                        统一接口                                  │
│                     SwitchAdapter                               │
└─────────────────────────────┬───────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            │                 │                 │
            ▼                 ▼                 ▼
    ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
    │  CLI Adapter  │ │ Netconf Adapter│ │ REST Adapter  │
    │  (SSH 命令行)  │ │ (NETCONF 协议) │ │ (HTTP API)    │
    └───────────────┘ └───────────────┘ └───────────────┘
```

**决策**: 不同交换机接口使用适配器模式统一为 `SwitchAdapter` 接口。

**理由**:
- 不同功能可使用不同适配器
- 新增适配方式无需修改业务代码
- 支持运行时切换适配器

---

## 2. 模式切换流程

### 2.1 模式切换时序图

```
用户请求切换模式
      │
      │ POST /api/v1/system/mode
      ▼
┌─────────────────┐
│ ModeHandler     │
└────┬────────────┘
     │
     │ SwitchMode(newMode)
     ▼
┌─────────────────┐
│ ModeResolver    │
└────┬────────────┘
     │
     ├─► 获取写锁 ──────────────────┐
     │                             │
     ├─► 验证新模式是否有效        │ 并
     │                             │ 发
     ├─► 持久化到数据库            │ 控
     │                             │ 制
     ├─► 清理旧模式资源            │
     │                             │
     ├─► 初始化新模式组件          │
     │                             │
     └─► 释放写锁 ─────────────────┘
     │
     ▼
返回切换成功
```

### 2.2 并发控制

```go
type ModeResolver struct {
    mu sync.RWMutex  // 读写锁
    currentMode string

    mockProvider    *MockProvider
    switchAdapters  map[string]SwitchAdapter
}

func (r *ModeResolver) SwitchMode(newMode string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 模式切换逻辑...
}

func (r *ModeResolver) GetProvider() DataProvider {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // 返回当前模式提供者
}
```

---

## 3. 数据流设计

### 3.1 离线测试模式数据流

```
用户请求
   │
   ▼
┌─────────────────┐
│   Handler       │
└────┬────────────┘
     │
     │ GetPortStatus()
     ▼
┌─────────────────┐
│ ModeResolver    │ ──► 判断 mode == "mock"
└────┬────────────┘
     │
     │ GetProvider()
     ▼
┌─────────────────┐
│ MockProvider    │
└────┬────────────┘
     │
     │ SELECT * FROM mock_port WHERE port_name = ?
     ▼
┌─────────────────┐
│ SQLite Database │
│   (mock_port)   │
└─────────────────┘
```

### 3.2 交换机模式数据流

```
用户请求
   │
   ▼
┌─────────────────┐
│   Handler       │
└────┬────────────┘
     │
     │ GetPortStatus()
     ▼
┌─────────────────┐
│ ModeResolver    │ ──► 判断 mode == "switch"
└────┬────────────┘
     │
     │ GetAdapter("port")
     ▼
┌─────────────────┐
│ CLI Adapter     │
└────┬────────────┘
     │
     │ SSH: show interface status
     ▼
┌─────────────────┐
│  交换机硬件     │
└─────────────────┘
```

### 3.3 用户模块特殊处理

```
程序启动
   │
   ▼
┌─────────────────┐
│ Bootstrap       │
└────┬────────────┘
     │
     │ GetRunMode()
     ▼
┌─────────────────┐
│ ModeResolver    │
└────┬────────────┘
     │
     ├─► mode == "mock"
     │   │
     │   └─► 从数据库读取用户
     │
     └─► mode == "switch"
         │
         ├─► 从交换机读取用户
         │
         ├─► 持久化到数据库
         │
         └─► 后续操作从数据库读取
```

---

## 4. 配置管理设计

### 4.1 系统配置表设计

```sql
-- 系统配置表
CREATE TABLE sys_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key VARCHAR(64) NOT NULL UNIQUE,
    config_value TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 配置项示例
INSERT INTO sys_config (config_key, config_value, description) VALUES
    ('run_mode', 'mock', '运行模式：mock=离线测试，switch=交换机'),
    ('cli_ssh_host', '192.168.1.1', 'CLI SSH 主机地址'),
    ('cli_ssh_port', '22', 'CLI SSH 端口'),
    ('cli_username', 'admin', 'CLI 登录用户名'),
    ('cli_password', '***', 'CLI 登录密码');
```

### 4.2 适配器配置表设计

```sql
-- 适配器配置表（支持一功能多适配器）
CREATE TABLE adapter_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    function_name VARCHAR(64) NOT NULL,      -- 功能：port, vlan, system...
    adapter_type VARCHAR(32) NOT NULL,       -- 类型：cli, netconf, rest
    priority INTEGER DEFAULT 0,              -- 优先级：数字越大优先级越高
    enabled INTEGER DEFAULT 1,               -- 是否启用
    config TEXT,                             -- JSON 格式配置
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(function_name, adapter_type)
);

-- 配置示例：端口功能使用 CLI，系统功能使用 Netconf
INSERT INTO adapter_config (function_name, adapter_type, priority, enabled, config) VALUES
    ('port', 'cli', 1, 1, '{"protocol":"ssh","host":"192.168.1.1","port":22}'),
    ('port', 'netconf', 0, 0, '{"host":"192.168.1.1","port":830}'),
    ('system', 'netconf', 1, 1, '{"host":"192.168.1.1","port":830}'),
    ('system', 'cli', 0, 0, '{"protocol":"ssh","host":"192.168.1.1","port":22}');
```

### 4.3 适配器选择逻辑

```go
func (r *AdapterRegistry) GetAdapter(functionName string) SwitchAdapter {
    // 查询数据库获取配置
    configs := r.dao.GetAdapterConfigs(functionName)

    // 过滤启用的配置
    enabledConfigs := filter(configs, func(c *AdapterConfig) bool {
        return c.Enabled
    })

    if len(enabledConfigs) == 0 {
        return nil
    }

    // 按优先级排序，选择最高优先级的适配器
    sort.Slice(enabledConfigs, func(i, j int) bool {
        return enabledConfigs[i].Priority > enabledConfigs[j].Priority
    })

    // 创建并返回适配器
    return r.factory.CreateAdapter(
        enabledConfigs[0].AdapterType,
        enabledConfigs[0].Config,
    )
}
```

---

## 5. 错误处理设计

### 5.1 错误类型定义

```go
// ErrModeSwitchFailed 模式切换失败
type ErrModeSwitchFailed struct {
    FromMode string
    ToMode   string
    Reason   error
}

func (e *ErrModeSwitchFailed) Error() string {
    return fmt.Sprintf("模式切换失败：%s -> %s, 原因：%v",
        e.FromMode, e.ToMode, e.Reason)
}

// ErrAdapterNotConfigured 适配器未配置
type ErrAdapterNotConfigured struct {
    FunctionName string
}

func (e *ErrAdapterNotConfigured) Error() string {
    return fmt.Sprintf("功能 [%s] 未配置适配器", e.FunctionName)
}

// ErrAdapterConnectFailed 适配器连接失败
type ErrAdapterConnectFailed struct {
    AdapterType string
    Reason      error
}

func (e *ErrAdapterConnectFailed) Error() string {
    return fmt.Sprintf("适配器 [%s] 连接失败：%v", e.AdapterType, e.Reason)
}
```

### 5.2 错误处理策略

| 错误类型 | 处理策略 | 用户提示 |
|----------|----------|----------|
| 模式切换失败 | 回滚到原模式 | "模式切换失败，已恢复为 {原模式}" |
| 适配器未配置 | 返回错误，要求配置 | "该功能未配置适配器，请先配置" |
| 适配器连接失败 | 重试 3 次，失败后降级 | "连接交换机失败，请检查网络" |
| 模拟数据不存在 | 返回默认值 | - |

---

## 6. 性能优化设计

### 6.1 连接池

```go
type CLIAdapterPool struct {
    pool *sync.Pool
    config *CLIConfig
}

func NewCLIAdapterPool(config *CLIConfig) *CLIAdapterPool {
    return &CLIAdapterPool{
        pool: &sync.Pool{
            New: func() interface{} {
                adapter := NewCLIAdapter(config)
                adapter.Connect()
                return adapter
            },
        },
        config: config,
    }
}

func (p *CLIAdapterPool) Get() *CLIAdapter {
    return p.pool.Get().(*CLIAdapter)
}

func (p *CLIAdapterPool) Put(adapter *CLIAdapter) {
    p.pool.Put(adapter)
}
```

### 6.2 缓存策略

```go
type CachedProvider struct {
    cache *ristretto.Cache
    provider DataProvider
    ttl time.Duration
}

func (p *CachedProvider) GetPortStatus(portID string) (*PortStatus, error) {
    // 先查缓存
    if cached, found := p.cache.Get("port:" + portID); found {
        return cached.(*PortStatus), nil
    }

    // 缓存未命中，查数据源
    status, err := p.provider.GetPortStatus(portID)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    p.cache.Set("port:" + portID, status, 1)
    return status, nil
}
```

---

## 7. 扩展点设计

### 7.1 新增适配器

```go
// 实现 SwitchAdapter 接口即可
type CustomAdapter struct {
    // 自定义配置
}

func (a *CustomAdapter) Name() string {
    return "custom"
}

func (a *CustomAdapter) Connect() error {
    // 实现连接逻辑
}

// ... 实现其他接口方法

// 注册到工厂
func init() {
    RegisterAdapter("custom", NewCustomAdapter)
}
```

### 7.2 新增模式

```go
// 实现 ModeStrategy 接口
type HybridModeStrategy struct {
    mockProvider *MockProvider
    switchAdapters map[string]SwitchAdapter
}

func (s *HybridModeStrategy) GetProvider() DataProvider {
    // 混合模式：部分数据来自 mock，部分来自交换机
    return &HybridProvider{
        mock: s.mockProvider,
        switch: s.switchAdapters,
    }
}
```

---

## 8. 测试策略

### 8.1 单元测试

```go
func TestMockProvider_GetPortStatus(t *testing.T) {
    provider := NewMockProvider(db)

    status, err := provider.GetPortStatus("GE1/0/1")

    assert.NoError(t, err)
    assert.Equal(t, "GE1/0/1", status.Name)
    assert.Equal(t, "up", status.LinkStatus)
}

func TestModeResolver_SwitchMode(t *testing.T) {
    resolver := NewModeResolver(db)

    err := resolver.SwitchMode("switch")

    assert.NoError(t, err)
    assert.Equal(t, "switch", resolver.GetCurrentMode())
}
```

### 8.2 集成测试

```go
func TestDualMode_Integration(t *testing.T) {
    // 1. 启动服务
    server := StartTestServer()
    defer server.Stop()

    // 2. 测试 mock 模式
    resp := server.Get("/api/v1/system/mode")
    assert.Equal(t, "mock", resp.Mode)

    // 3. 切换到 switch 模式
    server.Post("/api/v1/system/mode", map[string]string{"mode": "switch"})

    // 4. 验证切换成功
    resp = server.Get("/api/v1/system/mode")
    assert.Equal(t, "switch", resp.Mode)
}
```

---

## 9. 监控与日志

### 9.1 关键日志点

```go
// 模式切换日志
log.Infof("模式切换：%s -> %s", oldMode, newMode)

// 适配器选择日志
log.Debugf("功能 [%s] 选择适配器 [%s], 优先级 %d",
    functionName, adapterType, priority)

// 连接状态日志
log.Infof("适配器 [%s] 连接成功", adapter.Name())
log.Errorf("适配器 [%s] 连接失败：%v", adapter.Name(), err)
```

### 9.2 监控指标

```go
// Prometheus 指标
var (
    modeSwitchCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "switch_admin_mode_switches_total",
            Help: "模式切换次数",
        },
        []string{"from_mode", "to_mode"},
    )

    adapterRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "switch_admin_adapter_request_duration_seconds",
            Help: "适配器请求耗时",
        },
        []string{"adapter_type", "function"},
    )
)
```

---

**审批状态**: `pending`
**最后更新**: 2026-03-15
