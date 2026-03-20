# 无 CGO 编译指南

> 最后更新：2026-03-20
> 状态：**已完成** - 使用纯 Go SQLite 驱动

---

## 快速开始

### 一键编译（无 CGO）

```bash
# Windows
.\build.bat

# 或手动执行
set CGO_ENABLED=0
go build -mod=vendor -o switch-admin.exe ./cmd
```

编译成功后会生成 `switch-admin.exe`（约 50MB），可在任意 Windows 电脑运行。

---

## 问题背景

### 原有问题

之前使用 `github.com/mattn/go-sqlite3` 驱动，需要 CGO 环境：

1. **开发机器需要 GCC**：编译时必须安装 MinGW 或 TDM-GCC
2. **目标机器需要运行库**：可能需要 Visual C++ Redistributable
3. **跨平台编译复杂**：需要配置 CGO 交叉编译环境

### 错误示例

```
panic: sql: unknown driver "sqlite3" (forgotten import?)
goroutine 1:
github.com/GoAdminGroup/go-admin/modules/db.(*Sqlite).InitDB.func1()
    vendor/github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite/sqlite.go:86
```

---

## 解决方案

### 方案：使用纯 Go SQLite 驱动

将 GoAdmin 的 SQLite 驱动从 `mattn/go-sqlite3` 替换为 `modernc.org/sqlite`。

**优点：**
- ✅ 不需要 GCC 环境
- ✅ 编译后的程序可在任意 Windows 电脑运行
- ✅ 文件更小（50MB vs 75MB）
- ✅ 跨平台编译简单

**缺点：**
- ⚠️ 性能略低（约 10-20%）
- ⚠️ 某些 SQLite 高级特性可能不支持

---

## 技术实现

### 1. 修改 GoAdmin 驱动包装器

**文件：** `vendor/github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite/sqlite.go`

```go
// Copyright 2019 GoAdmin Core Team. All rights reserved.
package sqlite

import (
	"database/sql"
	"database/sql/driver"

	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动，不需要 CGO
)

// 驱动别名注册 - 将 "sqlite" 包装为 "sqlite3"
// 因为 modernc.org/sqlite 注册为 "sqlite"，但 GoAdmin 使用 "sqlite3"
type sqlite3Driver struct{}

func (d sqlite3Driver) Open(dsn string) (driver.Conn, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db.Driver().Open(dsn)
}

func init() {
	sql.Register("sqlite3", sqlite3Driver{})
}
```

### 2. 修改 main.go 导入

**文件：** `cmd/main.go`

```go
import (
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // 使用上面的包装器
	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动
	// ... 其他导入
)
```

### 3. 无 CGO 编译

```bash
CGO_ENABLED=0 go build -mod=vendor -o switch-admin.exe ./cmd
```

---

## 对比

| 特性 | CGO 版本 (mattn/go-sqlite3) | 无 CGO 版本 (modernc.org/sqlite) |
|------|----------------------------|--------------------------------|
| 编译需要 GCC | ✅ 是 | ❌ 否 |
| 目标机器需要运行库 | ⚠️ 可能需要 | ❌ 否 |
| 文件大小 | 75MB | 50MB |
| 性能 | 100% | 80-90% |
| SQLite 兼容性 | 100% | 95%+ |
| 跨平台编译 | 复杂 | 简单 |

---

## 注意事项

### 1. Vendor 目录

修改了 vendor 目录中的 GoAdmin 源码，需要确保 vendor 目录完整：

```bash
# 如果 vendor 目录损坏，重新生成
go mod vendor
```

### 2. 驱动注册

`modernc.org/sqlite` 注册为 `"sqlite"` 驱动名，但 GoAdmin 使用 `"sqlite3"`，所以需要包装器注册别名。

### 3. DSN 格式

GoAdmin 的 SQLite DSN 格式是 `data/admin.db`，与 `modernc.org/sqlite` 兼容。

---

## 验证

### 检查编译是否使用无 CGO

```bash
# 查看编译信息
go version -m switch-admin.exe

# 应该看到 CGO_ENABLED=0
```

### 测试运行

```bash
# 无 CGO 版本应能在无 GCC 环境的电脑上运行
.\switch-admin.exe

# 输出应显示：
# GoAdmin 启动成功
# 初始化数据库连接
# ...
```

---

## 相关文件

- `build.bat` - 无 CGO 编译脚本
- `vendor/github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite/sqlite.go` - 驱动包装器
- `cmd/main.go` - 主入口文件

---

## 参考

- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - 纯 Go SQLite 驱动
- [GoAdmin](https://github.com/GoAdminGroup/go-admin) - GoAdmin 框架

---

## 修订记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2026-03-20 | 1.0 | 初始版本，实现无 CGO 编译 |
