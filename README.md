# switch-admin

BroadEdge-S3652 智能交换机 Web 管理系统

## 架构设计

### 分层架构

本项目采用经典的分层架构设计，确保各层职责清晰、易于测试和扩展：

```
┌─────────────────────────────────────────────────────────────────┐
│                        表现层 (Presentation)                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  GoAdmin Pages  │  │  HTTP Handlers  │  │   API Routes    │  │
│  │  (datamodel/)   │  │  (handler/)     │  │  (cmd/main.go)  │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                        业务层 (Service)                          │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              DiagnosticService / Other Services          │    │
│  │  • 任务管理 (创建/查询/删除)                               │    │
│  │  • 业务逻辑编排                                           │    │
│  │  • 通过 ModeResolver 选择 Provider                        │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                      模式解析层 (ModeResolver)                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   ModeResolver                           │    │
│  │  • 读取当前运行模式 (mock / switch)                       │    │
│  │  • 根据模式返回对应的 Provider 实例                        │    │
│  │  • 支持运行时动态切换模式                                  │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                     提供者层 (Provider)                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  MockProvider   │  │  CLI Provider   │  │ NETCONF Provider│  │
│  │  (模拟数据)      │  │  (系统命令)      │  │  (未来扩展)     │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                        数据源层 (Data Source)                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  SQLite Database│  │  System Commands│  │  Switch Hardware│  │
│  │  (mock 模式)      │  │  (ping 命令)     │  │  (CLI/NETCONF)  │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 项目结构

```
switch-admin/
├── cmd/                      # 应用程序入口
│   └── main.go              # 主程序入口
├── internal/                 # 内部业务逻辑
│   ├── datamodel/           # 数据模型和页面定义
│   │   ├── bootstrap.go     # Bootstrap 配置
│   │   ├── dashboard.go     # Dashboard 页面
│   │   ├── generators.go    # 表格生成器配置
│   │   └── user.go          # 用户管理页面
│   ├── handler/             # HTTP 处理器
│   └── service/             # 业务服务
├── pkg/                      # 公共包
├── data/                     # 数据库文件
│   └── admin.db             # SQLite 数据库
├── uploads/                  # 上传文件目录
├── logs/                     # 日志目录
├── dist/                     # 编译输出目录
└── go.mod                    # Go 模块依赖
```

## 快速开始

### 环境要求

- Go 1.20+
- GCC (用于 CGO)
- SQLite3

### 安装依赖

```bash
go mod tidy
```

### 运行（开发模式）

```bash
CGO_ENABLED=1 CC=gcc go run ./cmd/main.go
```

### 编译

```bash
CGO_ENABLED=1 CC=gcc go build -o gin-admin.exe ./cmd/main.go
```

### 访问

启动后访问：http://localhost:9033/admin

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/mode` | GET | 获取当前运行模式 |
| `/api/mode` | POST | 切换运行模式 |
| `/api/system/config` | GET | 获取系统配置 |

### 模式切换示例

```bash
# 获取当前模式
curl http://localhost:9033/api/mode

# 切换到 mock 模式（离线测试模式）
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"mock"}'

# 切换到 switch 模式（交换机模式）
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"switch"}'
```

详细 API 文档请参考：[docs/API.md](docs/API.md)

默认账号密码请参考数据库中的初始数据。

### 模式判断流程

系统通过 `ModeResolver` 实现运行模式的判断和 Provider 路由：

```
                    ┌─────────────────────┐
                    │   请求到达 Service    │
                    │  (CreatePingTask)   │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  ModeResolver       │
                    │  .GetPingProvider() │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   判断当前模式       │
                    │  GetCurrentMode()   │
                    └──────────┬──────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
     ┌────────────────┐ ┌────────────────┐ ┌──────────────┐
     │ ModeMock       │ │ ModeSwitch     │ │ (未来扩展)   │
     │ (离线测试模式)  │ │ (交换机模式)    │ │              │
     └───────┬────────┘ └───────┬────────┘ └──────────────┘
             │                  │
             ▼                  ▼
     ┌────────────────┐ ┌────────────────┐
     │ MockPingProvider│ │ CLIPingProvider│
     │ (生成模拟数据)   │ │ (执行系统命令)  │
     └────────────────┘ └────────────────┘
```

### 调用链示例（以 Ping 功能为例）

#### 1. HTTP 请求 → Handler → Service

```
用户请求
  ↓
POST /api/v1/diagnostic/ping
  ↓
handler/diagnostic_handler.go:CreatePingTask()
  - 解析请求参数 (target, count, timeout, interval)
  - 调用 service.CreatePingTask(req)
  ↓
service/diagnostic_service.go:CreatePingTask()
  - 生成 task_id
  - 创建任务记录 (status: running)
  - 异步调用 executePing(taskID, req)
  - 返回 task_id 给前端
```

#### 2. Service → ModeResolver → Provider

```
executePing(taskID, req)
  ↓
modeResolver.GetPingProvider()
  - 读取 currentMode (mock / switch)
  - 返回对应的 Provider 实例
  ↓
Provider.ExecutePing(ctx, req)
  - MockPingProvider: 生成模拟数据
  - CLIPingProvider: 执行系统 ping 命令
```

#### 3. Provider 执行逻辑

**Mock 模式：**
```
MockPingProvider.ExecutePing()
  - 模拟网络延迟
  - 生成随机结果 (80% 成功率)
  - 计算统计信息 (丢包率/RTT)
  - 返回 PingTaskResponse
```

**CLI 模式：**
```
CLIPingProvider.ExecutePing()
  - 构建 ping 命令参数
    • Windows: ping -n 3 -w 2000 127.0.0.1
    • Linux: ping -c 3 -W 2 127.0.0.1
  - 执行系统命令 (exec.Command)
  - 解析输出 (正则匹配)
  - 计算统计信息
  - 返回 PingTaskResponse
```

#### 4. 结果存储与查询

```
Service (executePing 完成后)
  - 存储结果到 pingResults map
  - 更新任务状态 (status: completed)

前端轮询
  ↓
GET /api/v1/diagnostic/ping/{task_id}
  ↓
handler.GetPingTaskResult(taskID)
  ↓
service.GetPingTaskResult(taskID)
  - 从 pingResults map 读取结果
  - 返回 JSON 响应
```

### 已实现

- [x] 系统 Dashboard（首页）
  - 系统信息展示（产品型号、软件版本、运行时间、序列号）
  - 端口状态概览（端口总数、活跃端口、Down 端口、端口利用率）
  - 端口状态列表
  - 系统详细信息
- [x] 双模式架构
  - 离线测试模式 (mock)：使用数据库模拟数据
  - 交换机模式 (switch)：连接真实交换机硬件
  - 支持运行时动态切换，无需重启
- [x] 诊断功能（分层架构）
  - Ping 诊断：支持真实 CLI 模式和 Mock 模式
  - Traceroute 诊断（Mock 模式）
  - 虚拟电缆检测（Mock 模式）
- [x] 系统配置 API
  - 模式切换 API：`POST /api/mode`
  - 模式查询 API：`GET /api/mode`
  - 系统配置查询：`GET /api/system/config`

### 计划实现

- [ ] Traceroute CLI Provider
- [ ] CableTest CLI Provider
- [ ] NETCONF Provider 实现
- [ ] 端口状态管理页面
- [ ] 端口统计清零功能
- [ ] 系统信息详情页

## 技术栈

- **后端框架**: Go + Gin
- **管理框架**: GoAdmin
- **数据库**: SQLite3 (modernc.org/sqlite - 纯 Go 实现)
- **主题**: AdminLTE / Sword
- **UI 组件**: GoAdmin Template Components

## 核心接口

### PingProvider 接口

```go
type PingProvider interface {
    ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error)
}
```

### 实现

| Provider | 说明 | 使用场景 |
|----------|------|----------|
| MockPingProvider | 生成模拟数据 | 开发测试、离线演示 |
| CLIPingProvider | 执行系统 ping 命令 | 生产环境、真实诊断 |

### ModeResolver

```go
type ModeResolver interface {
    GetPingProvider() provider.PingProvider
    GetCurrentMode() RunMode
    SwitchMode(newMode RunMode) error
}
```

## 项目结构

```
switch-admin/
├── cmd/                      # 应用程序入口
│   └── main.go              # 主程序入口
├── internal/                 # 内部业务逻辑
│   ├── datamodel/           # 数据模型和页面定义
│   │   ├── bootstrap.go     # Bootstrap 配置
│   │   ├── dashboard.go     # Dashboard 页面
│   │   ├── diagnostic.go    # 诊断功能页面
│   │   └── ...
│   ├── handler/             # HTTP 处理器
│   │   ├── diagnostic_handler.go
│   │   └── ...
│   ├── model/               # 数据模型
│   │   ├── ping.go
│   │   ├── traceroute.go
│   │   └── ...
│   ├── service/             # 业务服务层
│   │   ├── diagnostic_service.go
│   │   ├── mode/            # 模式解析
│   │   │   └── mode_resolver.go
│   │   └── provider/        # Provider 实现
│   │       ├── diagnostic_provider.go  # 接口定义
│   │       ├── mock_ping_provider.go   # Mock 实现
│   │       └── cli_ping_provider.go    # CLI 实现
│   └── dao/                 # 数据访问层
├── pkg/                      # 公共包
├── scripts/                  # 工具脚本
├── data/                     # 数据库文件
│   └── admin.db             # SQLite 数据库
├── uploads/                  # 上传文件目录
├── logs/                     # 日志目录
├── dist/                     # 编译输出目录
├── openspec/                 # OpenSpec 规范文档
└── go.mod                    # Go 模块依赖
```

## 跨平台命令参数

CLIPingProvider 自动适配不同操作系统的 Ping 命令参数：

| 参数 | Windows | Linux |
|------|---------|-------|
| Count | `-n <count>` | `-c <count>` |
| Timeout | `-w <ms>` | `-W <seconds>` |
| Interval | 不支持 | `-i <seconds>` |

## 开发说明

## 开发说明

### 添加新页面

在 `internal/datamodel/` 目录下创建新的页面处理函数，然后在 `cmd/main.go` 中注册：

```go
// 1. 创建页面处理函数
func GetNewPageContent(ctx *context.Context) (types.Panel, error) {
    // ...
}

// 2. 在 main.go 中注册路由
e.HTML("GET", "/admin/newpage", datamodel.GetNewPageContent)
```

### 实现新的 Provider

实现新的诊断功能 Provider 需要：

1. **定义接口** (`internal/service/provider/diagnostic_provider.go`):
```go
type NewFeatureProvider interface {
    Execute(ctx context.Context, req model.NewFeatureRequest) (*model.NewFeatureResponse, error)
}
```

2. **实现 Mock Provider** (`internal/service/provider/mock_new_feature_provider.go`):
```go
type MockNewFeatureProvider struct{}

func (p *MockNewFeatureProvider) Execute(ctx context.Context, req model.NewFeatureRequest) (*model.NewFeatureResponse, error) {
    // 生成模拟数据
}
```

3. **实现 CLI Provider** (`internal/service/provider/cli_new_feature_provider.go`):
```go
type CLINewFeatureProvider struct {
    execFunc func(command string, args ...string) ([]byte, error)
}

func (p *CLINewFeatureProvider) Execute(ctx context.Context, req model.NewFeatureRequest) (*model.NewFeatureResponse, error) {
    // 1. 构建命令
    // 2. 执行命令
    // 3. 解析输出
    // 4. 返回结果
}
```

4. **扩展 ModeResolver** (`internal/service/mode/mode_resolver.go`):
```go
func (r *ModeResolver) GetNewFeatureProvider() provider.NewFeatureProvider {
    r.mu.RLock()
    defer r.mu.RUnlock()

    switch r.currentMode {
    case ModeSwitch:
        return r.cliNewFeatureProvider
    default:
        return r.mockNewFeatureProvider
    }
}
```

5. **在 Service 中使用**:
```go
provider := s.modeResolver.GetNewFeatureProvider()
result, err := provider.Execute(ctx, req)
```

### 数据库迁移

数据库文件位于 `data/admin.db`，使用 SQLite3 格式。

### 模式切换

通过 API 动态切换运行模式：

```bash
# 切换到 mock 模式（离线测试模式）
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"mock"}'

# 切换到 switch 模式（交换机模式）
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"switch"}'
```

### 调试技巧

1. **查看当前模式**: `GET /api/mode`
2. **查看模式描述**: `GET /api/system/config`
3. **日志文件**: `logs/` 目录
4. **数据库检查**: 使用 `scripts/check_menu.go` 查看菜单数据

## 注意事项

1. **Go Modules**: 项目使用 Go Modules 模式，bat 脚本已设置 `GOPATH=` 强制使用 Modules
2. **CGO**: 使用 `modernc.org/sqlite`（纯 Go 实现），无需 CGO 即可编译
3. **跨平台**:
   - Windows: 直接运行 `build_and_run.bat`
   - Linux/macOS: 使用 `go run ./cmd/main.go`
4. **端口占用**: 默认使用 9033 端口，确保该端口未被占用
5. **数据库初始化**: 首次运行时会自动创建数据库表和初始数据
6. **模式持久化**: 运行模式会持久化到数据库，重启后仍然有效

## 故障排查

### 常见问题

1. **端口被占用**
   ```
   listen tcp :9033: bind: Only one usage of each socket address is already permitted
   ```
   解决：`taskkill /F /IM switch-admin.exe`

2. **数据库文件不存在**
   解决：确保 `data/` 目录存在，程序会自动创建 `admin.db`

3. **Ping 功能无响应**
   - Mock 模式：检查是否正确生成模拟数据
   - CLI 模式：检查系统 ping 命令是否可用

### 日志位置

- 控制台输出：启动日志和错误信息
- 日志文件：`logs/` 目录（如配置）
