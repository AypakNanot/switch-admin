# Tasks - refactor-maintenance-subpackage

## 任务清单

### 阶段一：Handler 层按功能分组

#### Task 1: 创建 maintenance 包
- 创建 `internal/handler/maintenance/` 目录
- 包名为 `maintenance`（不是子包）
- 按功能分组到不同的 `.go` 文件

**验收**: 目录结构完整，所有 Handler 在同一个包内

---

#### Task 2: 创建 handler.go（主处理器）
- 包含 `Handler` 结构体
- 包含 `New()` 构造函数
- 包含重启/保存相关方法

**验收**: 编译通过

---

#### Task 3: 创建 system_config.go
- 系统配置相关 Handler 方法
- 网络配置、温度配置、设备信息、日期时间

**验收**: 编译通过

---

#### Task 4: 创建 file_management.go
- 文件管理相关 Handler 方法
- 文件上传、下载、删除、列表

**验收**: 编译通过

---

#### Task 5: 创建 log_management.go
- 日志管理相关 Handler 方法
- 日志列表、清除

**验收**: 编译通过

---

#### Task 6: 创建 snmp_config.go
- SNMP 配置相关 Handler 方法
- SNMP 配置、团体管理

**验收**: 编译通过

---

#### Task 7: 创建 snmp_trap.go
- SNMP Trap 相关 Handler 方法
- Trap 配置、目标主机管理、测试

**验收**: 编译通过

---

#### Task 8: 创建 worm_protection.go
- 蠕虫防护相关 Handler 方法
- 规则管理、统计清除

**验收**: 编译通过

---

#### Task 9: 创建 ddos_protection.go
- DDoS 防护相关 Handler 方法

**验收**: 编译通过

---

#### Task 10: 创建 arp_protection.go
- ARP 防护相关 Handler 方法

**验收**: 编译通过

---

#### Task 11: 创建 user_management.go
- 用户管理相关 Handler 方法
- 用户 CRUD 操作

**验收**: 编译通过

---

#### Task 12: 创建 session_management.go
- 会话管理相关 Handler 方法
- 会话列表、删除

**验收**: 编译通过

---

#### Task 13: 创建 load_config.go
- 加载配置相关 Handler 方法
- 配置文件列表、加载

**验收**: 编译通过

---

#### Task 14: Handler 层编译验证
- 运行 `go build`
- 运行 `go vet`

**验收**: 无编译错误

---

### 阶段二：DataModel 层按功能分组

#### Task 15: 创建 maintenance 包
- 创建 `internal/datamodel/maintenance/` 目录
- 包名为 `maintenance`

**验收**: 目录结构完整

---

#### Task 16: 创建 maintenance.go（导出文件）
- 导出所有页面内容函数
- 保持与 main.go 中的调用兼容

**验收**: 编译通过

---

#### Task 17-28: 创建各功能页面文件
- reboot_save.go - 重启/保存页面
- system_config.go - 系统配置页面
- file_management.go - 文件管理页面
- log_management.go - 日志管理页面
- snmp_config.go - SNMP 配置页面
- snmp_trap.go - SNMP Trap 页面
- worm_protection.go - 蠕虫防护页面
- ddos_protection.go - DDoS 防护页面
- arp_protection.go - ARP 防护页面
- user_management.go - 用户管理页面
- session_management.go - 会话管理页面
- load_config.go - 加载配置页面

**验收**: 所有文件创建完成，编译通过

---

### 阶段三：测试验证

#### Task 29: 全量编译验证
- 运行完整项目编译
- 确保所有包正常

**验收**: `go build` 成功

---

#### Task 30: 代码格式化
- 运行 `go fmt ./...`

**验收**: 代码格式化完成

---

## 依赖关系

```
Task 1 → Task 2-13 → Task 14
                       ↓
Task 15 → Task 16-28 → Task 29 → Task 30
```

## 进度追踪

- [x] 阶段一完成 (Handler 层)
- [x] 阶段二完成 (DataModel 层)
- [x] 阶段三完成 (测试验证)

---

**实现方案**: 按功能分组到同一包内的不同文件（不是子包）
**验收标准**: 编译通过，代码格式化，无运行时错误
