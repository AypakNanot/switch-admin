# switch-admin API 文档

> BroadEdge-S3652 智能交换机 Web 管理系统 - RESTful API 参考
>
> 版本：1.0
> 更新日期：2026-03-15

---

## 快速开始

### 基础信息

- **基础 URL**: `http://localhost:9033`
- **认证方式**: 部分接口需要登录（GoAdmin Session）
- **数据格式**: JSON (`Content-Type: application/json`)

### 健康检查

```bash
curl http://localhost:9033/health
```

**响应示例**:
```json
{"status": "ok"}
```

---

## 系统模式 API

系统支持两种运行模式：
- **mock** - 离线测试模式：使用数据库模拟数据，不依赖真实交换机硬件
- **switch** - 交换机模式：连接真实交换机，通过 CLI/Netconf 等方式操作硬件

### 获取当前模式

获取系统当前运行模式。

**请求**:
```http
GET /api/mode
```

**cURL 示例**:
```bash
curl http://localhost:9033/api/mode
```

**成功响应** (HTTP 200):
```json
{
  "code": 200,
  "data": {
    "mode": "switch",
    "description": "交换机模式"
  }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| mode | string | 运行模式代码：`mock` 或 `switch` |
| description | string | 运行模式中文描述 |

---

### 切换运行模式

切换系统运行模式，无需重启服务。

**请求**:
```http
POST /api/mode
Content-Type: application/json

{
  "mode": "mock"
}
```

**cURL 示例**:
```bash
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"mock"}'
```

**成功响应** (HTTP 200):
```json
{
  "code": 200,
  "data": {
    "previous_mode": "switch",
    "current_mode": "mock"
  },
  "message": "模式切换成功"
}
```

**错误响应** - 无效模式 (HTTP 400):
```json
{
  "code": 400,
  "message": "Invalid mode. Must be 'mock' or 'switch'"
}
```

**错误响应** - 无效 JSON (HTTP 400):
```json
{
  "code": 400,
  "message": "Invalid JSON format"
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| mode | string | 是 | 目标模式，只能是 `mock` 或 `switch` |

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| previous_mode | string | 切换前的模式 |
| current_mode | string | 切换后的模式 |

---

### 获取系统配置

获取系统完整配置信息。

**请求**:
```http
GET /api/system/config
```

**cURL 示例**:
```bash
curl http://localhost:9033/api/system/config
```

**成功响应** (HTTP 200):
```json
{
  "code": 200,
  "data": {
    "mode": "switch",
    "mode_description": "交换机模式",
    "database": "SQLite3 (data/admin.db)",
    "goadmin_version": "v1.2.26"
  }
}
```

---

## 错误码说明

| HTTP 状态码 | 说明 |
|------------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 |

---

## 使用示例

### 示例 1：检查当前模式并切换

```bash
# 1. 检查当前模式
curl http://localhost:9033/api/mode

# 2. 切换到 mock 模式
curl -X POST http://localhost:9033/api/mode \
  -H "Content-Type: application/json" \
  -d '{"mode":"mock"}'

# 3. 验证切换结果
curl http://localhost:9033/api/mode
```

### 示例 2：在脚本中使用

```bash
#!/bin/bash

API_BASE="http://localhost:9033"

# 获取当前模式
current_mode=$(curl -s $API_BASE/api/mode | jq -r '.data.mode')
echo "当前模式：$current_mode"

# 切换模式
if [ "$current_mode" = "mock" ]; then
  target_mode="switch"
else
  target_mode="mock"
fi

echo "切换到：$target_mode"
curl -X POST $API_BASE/api/mode \
  -H "Content-Type: application/json" \
  -d "{\"mode\":\"$target_mode\"}"
```

---

## GoAdmin 管理界面

| 路径 | 说明 |
|------|------|
| `/admin` | 管理后台首页 |
| `/admin/dashboard` | Dashboard（系统概览） |
| `/admin/system/config` | 系统配置页面 |
| `/admin/signin` | 登录页面 |

---

## 附录：模式说明

### mock 模式（离线测试模式）

- **用途**: 开发测试、演示环境、无交换机硬件时使用
- **数据源**: SQLite 数据库模拟数据
- **特点**:
  - 不依赖真实交换机硬件
  - 数据可手动配置
  - 适合开发和功能演示

### switch 模式（交换机模式）

- **用途**: 生产环境、连接真实交换机
- **数据源**: 真实交换机硬件
- **特点**:
  - 通过 CLI/Netconf 等方式操作交换机
  - 实时读取交换机状态
  - 支持配置下发

---

*文档版本 1.0 | 最后更新：2026-03-15*
