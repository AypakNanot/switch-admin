# switch-admin API 文档

> BroadEdge-S3652 智能交换机 Web 管理系统 - RESTful API 参考
>
> 版本：2.0
> 更新日期：2026-03-17

---

## 目录

1. [基础信息](#基础信息)
2. [系统模块 API](#系统模块-api)
3. [配置模块 API](#配置模块-api)
4. [网络模块 API](#网络模块-api)
5. [维护模块 API](#维护模块-api)
6. [安全模块 API](#安全模块-api)

---

## 基础信息

- **基础 URL**: `http://localhost:9033`
- **认证方式**: 部分接口需要登录（GoAdmin Session）
- **数据格式**: JSON (`Content-Type: application/json`)

### 通用响应格式

```json
{
  "code": 200,
  "data": { ... },
  "message": "操作成功"
}
```

### 错误响应格式

```json
{
  "code": 400,
  "message": "错误描述",
  "error": "错误代码"
}
```

---

## 系统模块 API

### 健康检查

**请求**:
```http
GET /health
```

**响应示例**:
```json
{"status": "ok"}
```

---

### 获取运行模式

获取系统当前运行模式。

**请求**:
```http
GET /api/mode
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "mode": "switch",
    "description": "交换机模式"
  }
}
```

---

### 切换运行模式

切换系统运行模式，无需重启服务。

**请求**:
```http
POST /api/mode
Content-Type: application/json

{"mode": "mock"}
```

**响应示例**:
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

**错误响应**:
```json
{
  "code": 400,
  "message": "Invalid mode. Must be 'mock' or 'switch'"
}
```

---

### 获取系统配置

**请求**:
```http
GET /api/system/config
```

**响应示例**:
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

## 配置模块 API

### 端口管理

#### 获取端口列表

**请求**:
```http
GET /api/v1/ports
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "ports": [
      {
        "port_id": "GE1/0/1",
        "admin_status": "enable",
        "link_status": "up",
        "speed_duplex": "1000F",
        "flow_control": "off",
        "description": "Server-A",
        "aggregation": "-"
      }
    ]
  }
}
```

---

#### 更新端口配置

**请求**:
```http
PUT /api/v1/ports/:port_id
Content-Type: application/json

{
  "admin_status": "enable",
  "description": "Server-A",
  "speed_duplex": "1000F",
  "flow_control": "off"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "端口配置保存成功"
}
```

---

### 链路聚合

#### 获取链路聚合列表

**请求**:
```http
GET /api/v1/link-aggregation
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "aggregations": [
      {
        "group_id": 1,
        "name": "Ag1",
        "mode": "LACP",
        "member_ports": ["GE1/0/1", "GE1/0/2", "GE1/0/3", "GE1/0/4"],
        "load_balance": "src-dst-ip",
        "min_active": 2,
        "status": "normal"
      }
    ]
  }
}
```

---

#### 创建链路聚合组

**请求**:
```http
POST /api/v1/link-aggregation
Content-Type: application/json

{
  "group_id": 1,
  "mode": "LACP",
  "description": "Uplink-Agg",
  "load_balance": "src-dst-ip",
  "member_ports": ["GE1/0/1", "GE1/0/2"],
  "min_active": 2,
  "lacp_timeout": "slow",
  "priority": 100
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "链路聚合组创建成功"
}
```

---

#### 更新链路聚合组

**请求**:
```http
PUT /api/v1/link-aggregation/:id
Content-Type: application/json
```

**响应示例**:
```json
{
  "code": 200,
  "message": "链路聚合组更新成功"
}
```

---

#### 删除链路聚合组

**请求**:
```http
DELETE /api/v1/link-aggregation/:id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "链路聚合组删除成功"
}
```

---

## 网络模块 API

### 路由管理

#### 获取 IPv4 路由表

**请求**:
```http
GET /api/v1/routes/table
```

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| dest_ip | string | - | 目的 IP（模糊匹配） |
| protocol | string | - | 协议类型（Static/Connected/OSPF/RIP/BGP） |
| page | int | 1 | 页码 |
| page_size | int | 50 | 每页数量 |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "total": 10,
    "page": 1,
    "page_size": 50,
    "total_pages": 1,
    "items": [
      {
        "dest_ip": "0.0.0.0",
        "dest_mask": "0.0.0.0",
        "protocol": "Static",
        "out_port": "vlan1",
        "next_hop": "192.168.1.1",
        "metric": 1,
        "preference": 60
      }
    ]
  }
}
```

---

#### 获取静态路由列表

**请求**:
```http
GET /api/v1/routes/static
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": "1",
        "dest_ip": "0.0.0.0",
        "dest_mask": "0.0.0.0",
        "next_hop": "192.168.1.1",
        "distance": 1,
        "status": "active",
        "status_desc": ""
      }
    ]
  }
}
```

---

#### 获取单条静态路由

**请求**:
```http
GET /api/v1/routes/static/:id
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": "1",
    "dest_ip": "0.0.0.0",
    "dest_mask": "0.0.0.0",
    "next_hop": "192.168.1.1",
    "distance": 1,
    "status": "active"
  }
}
```

---

#### 创建静态路由

**请求**:
```http
POST /api/v1/routes/static
Content-Type: application/json

{
  "dest_ip": "10.0.0.0",
  "dest_mask": "255.255.255.0",
  "next_hop": "192.168.1.1",
  "distance": 1
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "静态路由配置成功",
  "data": {
    "id": "4",
    "dest_ip": "10.0.0.0",
    "dest_mask": "255.255.255.0",
    "next_hop": "192.168.1.1",
    "distance": 1,
    "status": "active"
  }
}
```

---

#### 更新静态路由

**请求**:
```http
PUT /api/v1/routes/static/:id
Content-Type: application/json
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功"
}
```

---

#### 删除静态路由

**请求**:
```http
DELETE /api/v1/routes/static/:id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "deleted_id": "1"
  }
}
```

---

### 诊断工具

#### Ping 测试

##### 创建 Ping 任务

**请求**:
```http
POST /api/v1/diagnostic/ping
Content-Type: application/json

{
  "target": "8.8.8.8",
  "count": 4,
  "timeout": 2,
  "interval": 1,
  "vrf_id": "mgmt vrf"
}
```

**参数说明**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| target | string | - | 目标 IP 或域名 |
| count | int | 4 | Ping 次数 |
| timeout | int | 2 | 超时时间（秒） |
| interval | int | 1 | 间隔时间（秒） |
| vrf_id | string | "mgmt vrf" | VRF ID |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "task_id": "ping_12345",
    "status": "running"
  }
}
```

---

##### 获取 Ping 任务结果

**请求**:
```http
GET /api/v1/diagnostic/ping/:task_id
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "target": "8.8.8.8",
    "packets_sent": 4,
    "packets_received": 4,
    "packet_loss": "0%",
    "rtt_min": 10,
    "rtt_avg": 12,
    "rtt_max": 15,
    "results": [...]
  }
}
```

---

##### 删除 Ping 任务

**请求**:
```http
DELETE /api/v1/diagnostic/ping/:task_id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Task deleted"
}
```

---

#### Traceroute 跟踪

##### 创建 Traceroute 任务

**请求**:
```http
POST /api/v1/diagnostic/traceroute
Content-Type: application/json

{
  "target": "8.8.8.8",
  "max_hops": 30,
  "timeout": 2,
  "probes": 3,
  "vrf_id": "mgmt vrf"
}
```

**参数说明**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| target | string | - | 目标 IP 或域名 |
| max_hops | int | 30 | 最大跳数 |
| timeout | int | 2 | 超时时间（秒） |
| probes | int | 3 | 每跳探测次数 |
| vrf_id | string | "mgmt vrf" | VRF ID |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "task_id": "traceroute_12345",
    "status": "running"
  }
}
```

---

##### 获取 Traceroute 任务结果

**请求**:
```http
GET /api/v1/diagnostic/traceroute/:task_id
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "target": "8.8.8.8",
    "hops": [
      {
        "hop": 1,
        "ip": "192.168.1.1",
        "hostname": "gateway.local",
        "rtt": [1.2, 1.1, 1.3]
      }
    ]
  }
}
```

---

##### 删除 Traceroute 任务

**请求**:
```http
DELETE /api/v1/diagnostic/traceroute/:task_id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Task deleted"
}
```

---

#### 电缆检测

##### 获取可检测端口列表

**请求**:
```http
GET /api/v1/diagnostic/cable/ports
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "ports": [
      {"port_id": "GE1/0/1", "type": "electrical"},
      {"port_id": "GE1/0/2", "type": "electrical"}
    ]
  }
}
```

---

##### 执行电缆检测

**请求**:
```http
POST /api/v1/diagnostic/cable
Content-Type: application/json

{
  "port_id": "GE1/0/1"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "port_id": "GE1/0/1",
    "pair_a": {"length": 50, "status": "normal"},
    "pair_b": {"length": 50, "status": "normal"},
    "pair_c": {"length": 50, "status": "normal"},
    "pair_d": {"length": 50, "status": "normal"}
  }
}
```

**错误响应 - 端口不支持**:
```json
{
  "code": 400,
  "error": "PORT_NOT_ELECTRICAL",
  "message": "光口不支持虚拟电缆检测",
  "data": {
    "port_id": "GE1/0/49",
    "port_type": "optical"
  }
}
```

---

## 维护模块 API

### 系统控制

#### 保存配置

**请求**:
```http
POST /api/v1/system/save-config
```

**响应示例**:
```json
{
  "code": 200,
  "message": "配置保存成功"
}
```

---

#### 重启设备

**请求**:
```http
POST /api/v1/system/reboot
Content-Type: application/json

{
  "save_before_reboot": true
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "设备正在重启，请稍后刷新页面"
}
```

---

#### 恢复出厂配置

**请求**:
```http
POST /api/v1/system/factory-reset
Content-Type: application/json

{
  "confirmation": "CONFIRM"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "设备正在恢复出厂配置并重启",
  "data": {
    "default_ip": "192.168.1.1",
    "default_subnet": "255.255.255.0",
    "default_user": "admin",
    "default_password": "admin"
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "error": "CONFIRMATION_REQUIRED",
  "message": "请输入 'CONFIRM' 以确认恢复出厂配置"
}
```

---

### 系统配置

#### 获取系统配置

**请求**:
```http
GET /api/v1/system/config
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "network": {
      "ip_address": "192.168.1.100",
      "subnet": "255.255.255.0",
      "gateway": "192.168.1.1"
    },
    "temperature": {
      "low_alarm": 5,
      "high_warn": 65,
      "high_alarm": 80
    },
    "device_info": {
      "device_name": "Switch-001",
      "contact": "admin@example.com",
      "location": "机房 A-01"
    },
    "datetime": "2026-03-17 10:30:00",
    "timezone": "UTC+8"
  }
}
```

---

#### 更新网络配置

**请求**:
```http
PUT /api/v1/system/network
Content-Type: application/json

{
  "ip_address": "192.168.1.100",
  "subnet": "255.255.255.0",
  "gateway": "192.168.1.1"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "网络配置更新成功"
}
```

---

#### 更新温度阈值配置

**请求**:
```http
PUT /api/v1/system/temperature
Content-Type: application/json

{
  "low_alarm": 5,
  "high_warn": 65,
  "high_alarm": 80
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "温度阈值配置更新成功"
}
```

---

#### 更新设备信息

**请求**:
```http
PUT /api/v1/system/info
Content-Type: application/json

{
  "device_name": "Switch-001",
  "contact": "admin@example.com",
  "location": "机房 A-01"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "设备信息更新成功"
}
```

---

#### 更新时间日期

**请求**:
```http
PUT /api/v1/system/datetime
Content-Type: application/json

{
  "datetime": "2026-03-17 10:30:00",
  "timezone": "UTC+8"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "时间日期更新成功"
}
```

---

### 配置管理

#### 获取配置文件列表

**请求**:
```http
GET /api/v1/config/files
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "files": [
      {
        "file_path": "flash:/startup-config.conf",
        "modified": "2026-03-11 10:17:31",
        "size": "1.8K"
      },
      {
        "file_path": "flash:/backup-config.conf",
        "modified": "2026-03-10 08:30:00",
        "size": "1.7K"
      }
    ],
    "total": 2
  }
}
```

---

#### 加载配置文件

**请求**:
```http
POST /api/v1/config/load
Content-Type: application/json

{
  "file_path": "flash:/backup-config.conf"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "配置加载成功，部分配置可能需要重启生效"
}
```

---

### 文件管理

#### 获取文件列表

**请求**:
```http
GET /api/v1/files
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "files": [
      {"filename": "startup-config.conf", "directory": "flash:/", "size": "1.8K", "modified": "2026-03-11 10:17:31"}
    ],
    "storage": {
      "flash_total": "3.9G",
      "flash_free": "3.6G",
      "boot_total": "2.9G",
      "boot_free": "2.3G"
    }
  }
}
```

---

#### 上传文件

**请求**:
```http
POST /api/v1/files/upload
Content-Type: multipart/form-data

file: <file>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文件上传成功"
}
```

---

#### 上传固件

**请求**:
```http
POST /api/v1/files/firmware
Content-Type: multipart/form-data

firmware: <file>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "固件上传成功，正在校验...",
  "data": {
    "filename": "firmware.bin",
    "size": 123456789
  }
}
```

---

#### 下载文件

**请求**:
```http
GET /api/v1/files/download?file_path=flash:/file.conf
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文件下载成功",
  "data": {
    "file_path": "flash:/file.conf",
    "content": "文件内容"
  }
}
```

---

#### 删除文件

**请求**:
```http
DELETE /api/v1/files
Content-Type: application/json

{
  "file_paths": ["flash:/file1.conf", "flash:/file2.conf"]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文件删除成功"
}
```

---

### 日志管理

#### 获取日志列表

**请求**:
```http
GET /api/v1/logs
```

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| start_time | string | - | 开始时间 |
| end_time | string | - | 结束时间 |
| level | string | All | 日志级别 |
| module | string | All | 模块 |
| page | int | 1 | 页码 |
| page_size | int | 50 | 每页数量 |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "logs": [
      {
        "time": "2026-03-12 13:32:15",
        "module": "DHCLIENT",
        "level": "Info",
        "content": "Interface vlan1 renew success"
      }
    ],
    "total": 300,
    "page": 1,
    "page_size": 50
  }
}
```

---

#### 清除所有日志

**请求**:
```http
DELETE /api/v1/logs
```

**响应示例**:
```json
{
  "code": 200,
  "message": "日志清除成功"
}
```

---

### SNMP 配置

#### 获取 SNMP 配置

**请求**:
```http
GET /api/v1/snmp/config
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "enabled": true,
    "version": "v2c",
    "communities": [
      {"name": "public", "access": "Read-Only"},
      {"name": "private", "access": "Read-Write"}
    ]
  }
}
```

---

#### 更新 SNMP 配置

**请求**:
```http
PUT /api/v1/snmp/config
Content-Type: application/json

{
  "enabled": true,
  "version": "v2c"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP 配置更新成功"
}
```

---

#### 获取 SNMP 团体列表

**请求**:
```http
GET /api/v1/snmp/communities
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "communities": [
      {"name": "public", "access": "Read-Only"}
    ]
  }
}
```

---

#### 添加 SNMP 团体

**请求**:
```http
POST /api/v1/snmp/communities
Content-Type: application/json

{
  "name": "newcommunity",
  "access": "Read-Only"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP 团体添加成功"
}
```

---

#### 删除 SNMP 团体

**请求**:
```http
DELETE /api/v1/snmp/communities/:name
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP 团体删除成功"
}
```

---

### SNMP Trap 配置

#### 获取 SNMP Trap 配置

**请求**:
```http
GET /api/v1/snmp/trap/config
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "trap_enabled": {
      "coldstart": false,
      "warmstart": false,
      "linkup": false,
      "linkdown": false,
      "system": false,
      "loopback": false
    },
    "trap_hosts": [
      {"address": "192.168.1.100", "port": 162, "vrf": "mgmt-if", "community": "public"}
    ]
  }
}
```

---

#### 更新 SNMP Trap 配置

**请求**:
```http
PUT /api/v1/snmp/trap/config
Content-Type: application/json

{
  "trap_enabled": {
    "coldstart": true,
    "linkup": true,
    "linkdown": true
  }
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP Trap 配置更新成功"
}
```

---

#### 获取 SNMP Trap 目标主机列表

**请求**:
```http
GET /api/v1/snmp/trap/hosts
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "trap_hosts": [
      {"id": 1, "address": "192.168.1.100", "port": 162, "vrf": "mgmt-if", "community": "public"}
    ]
  }
}
```

---

#### 添加 SNMP Trap 目标主机

**请求**:
```http
POST /api/v1/snmp/trap/hosts
Content-Type: application/json

{
  "address": "192.168.1.100",
  "port": 162,
  "vrf": "mgmt-if",
  "community": "public"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP Trap 目标主机添加成功"
}
```

---

#### 删除 SNMP Trap 目标主机

**请求**:
```http
DELETE /api/v1/snmp/trap/hosts/:id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "SNMP Trap 目标主机删除成功"
}
```

---

#### 发送测试 Trap

**请求**:
```http
POST /api/v1/snmp/trap/hosts/:id/test
```

**响应示例**:
```json
{
  "code": 200,
  "message": "测试 Trap 已发送，请检查网管平台是否收到告警"
}
```

---

### 用户管理

#### 获取用户列表

**请求**:
```http
GET /api/v1/users
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "users": [
      {
        "username": "admin",
        "role": 0,
        "role_name": "super-admin",
        "created_at": "2026-01-01 00:00:00"
      }
    ],
    "total": 2
  }
}
```

---

#### 创建用户

**请求**:
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "role": 2
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户创建成功"
}
```

---

#### 更新用户

**请求**:
```http
PUT /api/v1/users/:username
Content-Type: application/json

{
  "password": "newpassword",
  "role": 1
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户更新成功"
}
```

---

#### 批量删除用户

**请求**:
```http
DELETE /api/v1/users
Content-Type: application/json

{
  "usernames": ["user1", "user2"]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户删除成功"
}
```

---

### 会话管理

#### 获取会话列表

**请求**:
```http
GET /api/v1/sessions
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "sessions": [
      {
        "session_id": "1773322174",
        "username": "admin",
        "timeout": "2026-03-15 23:01:14",
        "client_ip": "192.168.1.50",
        "is_current": true
      }
    ],
    "total": 1
  }
}
```

---

#### 删除会话（强制踢出）

**请求**:
```http
DELETE /api/v1/sessions/:session_id
```

**响应示例**:
```json
{
  "code": 200,
  "message": "会话已终止"
}
```

---

## 安全模块 API

### 蠕虫攻击防护

#### 获取蠕虫规则列表

**请求**:
```http
GET /api/v1/security/worm/rules
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "rules": [
      {
        "id": 1,
        "name": "NachiBlasterD",
        "protocol": "tcp",
        "port": 707,
        "attack_count": 0,
        "enabled": false
      }
    ],
    "total": 7
  }
}
```

---

#### 添加蠕虫规则

**请求**:
```http
POST /api/v1/security/worm/rules
Content-Type: application/json

{
  "name": "CustomWorm",
  "protocol": "tcp",
  "port": 1234,
  "enabled": true
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "蠕虫规则添加成功"
}
```

---

#### 更新蠕虫规则

**请求**:
```http
PUT /api/v1/security/worm/rules/:id
Content-Type: application/json

{
  "name": "CustomWorm",
  "protocol": "tcp",
  "port": 1234,
  "enabled": true
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "蠕虫规则更新成功"
}
```

---

#### 批量删除蠕虫规则

**请求**:
```http
DELETE /api/v1/security/worm/rules
Content-Type: application/json

{
  "ids": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "蠕虫规则删除成功"
}
```

---

#### 清除攻击统计

**请求**:
```http
POST /api/v1/security/worm/clear-stats
```

**响应示例**:
```json
{
  "code": 200,
  "message": "攻击统计已清零"
}
```

---

### DDoS 攻击防护

#### 获取 DDoS 防护配置

**请求**:
```http
GET /api/v1/security/ddos/config
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "icmp_flooding": 0,
    "udp_flooding": 0,
    "syn_flooding": 0,
    "small_packet_size": 28,
    "smurf_protection": false,
    "fraggle_protection": false,
    "mac_equal_protection": false,
    "ip_equal_protection": false
  }
}
```

---

#### 更新 DDoS 防护配置

**请求**:
```http
PUT /api/v1/security/ddos/config
Content-Type: application/json

{
  "icmp_flooding": 100,
  "udp_flooding": 100,
  "syn_flooding": 100,
  "small_packet_size": 28,
  "smurf_protection": true,
  "fraggle_protection": true,
  "mac_equal_protection": true,
  "ip_equal_protection": true
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "DDoS 防护配置更新成功"
}
```

---

### ARP 攻击防护

#### 获取 ARP 防护配置

**请求**:
```http
GET /api/v1/security/arp/config
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "arp_rate_limit": 100
  }
}
```

---

#### 更新 ARP 防护配置

**请求**:
```http
PUT /api/v1/security/arp/config
Content-Type: application/json

{
  "arp_rate_limit": 100
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "ARP 防护配置更新成功"
}
```

**错误响应 - 危险设置**:
```json
{
  "code": 400,
  "error": "DANGEROUS_SETTING",
  "message": "设置为 0 将导致交换机停止学习动态 ARP，可能失去管理权限"
}
```

---

## 错误码说明

| HTTP 状态码 | 说明 |
|------------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## GoAdmin 管理界面

| 路径 | 说明 |
|------|------|
| `/admin` | 管理后台首页 |
| `/admin/dashboard` | Dashboard（系统概览） |
| `/admin/system/config` | 系统配置页面 |
| `/admin/signin` | 登录页面 |

---

*文档版本 2.0 | 最后更新：2026-03-17*
