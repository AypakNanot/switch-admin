# switch-admin

BroadEdge-S3652 智能交换机 Web 管理系统

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

默认账号密码请参考数据库中的初始数据。

## 功能模块

### 已实现

- [x] 系统 Dashboard（首页）
  - 系统信息展示（产品型号、软件版本、运行时间、序列号）
  - 端口状态概览（端口总数、活跃端口、Down 端口、端口利用率）
  - 端口状态列表
  - 系统详细信息

### 计划实现

- [ ] 端口状态管理页面
- [ ] 端口统计清零功能
- [ ] 系统信息详情页
- [ ] 用户管理（已有）

## 技术栈

- **后端框架**: Go + Gin
- **管理框架**: GoAdmin
- **数据库**: SQLite3
- **主题**: AdminLTE
- **UI 组件**: GoAdmin Template Components

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

### 数据库迁移

数据库文件位于 `data/admin.db`，使用 SQLite3 格式。

## 注意事项

1. 运行时需要设置 `CGO_ENABLED=1`，因为使用了 sqlite3 驱动
2. Windows 系统下需要使用 GCC 编译器（推荐安装 MSYS2 或 MinGW）
3. 不要在 Bash (MINGW64) 中直接运行 .exe 文件，应使用 `go run` 命令
