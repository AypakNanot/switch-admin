# Tasks for refactor-maintenance-module

## 阶段一：Handler 层拆分

### 任务 1: 创建 Handler 子模块文件 (1-6)
- [ ] 创建 `internal/handler/maintenance_reboot_save.go`
  - [ ] 移动 `SaveConfig` 方法
  - [ ] 移动 `RebootSwitch` 方法
  - [ ] 移动 `FactoryReset` 方法
- [ ] 创建 `internal/handler/maintenance_system_config.go`
  - [ ] 移动 `GetSystemConfig` 方法
  - [ ] 移动 `UpdateNetworkConfig` 方法
  - [ ] 移动 `UpdateTemperatureConfig` 方法
  - [ ] 移动 `UpdateDeviceInfo` 方法
  - [ ] 移动 `UpdateDateTime` 方法
- [ ] 创建 `internal/handler/maintenance_file_management.go`
  - [ ] 移动 `GetFiles` 方法
  - [ ] 移动 `UploadFile` 方法
  - [ ] 移动 `UploadFirmware` 方法
  - [ ] 移动 `DownloadFile` 方法
  - [ ] 移动 `DeleteFiles` 方法
- [ ] 创建 `internal/handler/maintenance_log_management.go`
  - [ ] 移动 `GetLogs` 方法
  - [ ] 移动 `ClearLogs` 方法
- [ ] 创建 `internal/handler/maintenance_snmp_config.go`
  - [ ] 移动 `GetSNMPConfig` 方法
  - [ ] 移动 `UpdateSNMPConfig` 方法
  - [ ] 移动 `AddSNMPCommunity` 方法
  - [ ] 移动 `GetSNMPCommunity` 方法
  - [ ] 移动 `DeleteSNMPCommunity` 方法
- [ ] 创建 `internal/handler/maintenance_snmp_trap.go`
  - [ ] 移动 `GetSNMPTrapConfig` 方法
  - [ ] 移动 `GetSNMPTrapHosts` 方法
  - [ ] 移动 `UpdateSNMPTrapConfig` 方法
  - [ ] 移动 `AddSNMPTrapHost` 方法
  - [ ] 移动 `DeleteSNMPTrapHost` 方法
  - [ ] 移动 `TestSNMPTrap` 方法

### 任务 2: 创建 Handler 子模块文件 (7-12)
- [ ] 创建 `internal/handler/maintenance_worm_protection.go`
  - [ ] 移动 `GetWormRules` 方法
  - [ ] 移动 `AddWormRule` 方法
  - [ ] 移动 `UpdateWormRule` 方法
  - [ ] 移动 `DeleteWormRules` 方法
  - [ ] 移动 `ClearWormStats` 方法
- [ ] 创建 `internal/handler/maintenance_ddos_protection.go`
  - [ ] 移动 `GetDDoSConfig` 方法
  - [ ] 移动 `UpdateDDoSConfig` 方法
- [ ] 创建 `internal/handler/maintenance_arp_protection.go`
  - [ ] 移动 `GetARPConfig` 方法
  - [ ] 移动 `UpdateARPConfig` 方法
- [ ] 创建 `internal/handler/maintenance_user_management.go`
  - [ ] 移动 `GetUsers` 方法
  - [ ] 移动 `CreateUser` 方法
  - [ ] 移动 `UpdateUser` 方法
  - [ ] 移动 `DeleteUsers` 方法
- [ ] 创建 `internal/handler/maintenance_session_management.go`
  - [ ] 移动 `GetSessions` 方法
  - [ ] 移动 `DeleteSession` 方法
- [ ] 创建 `internal/handler/maintenance_load_config.go`
  - [ ] 移动 `GetConfigFiles` 方法
  - [ ] 移动 `LoadConfig` 方法

### 任务 3: 重构 maintenance_handler.go
- [ ] 添加子模块处理器结构体引用
- [ ] 更新 `NewMaintenanceHandler()` 构造函数
- [ ] 移除已拆分的 Handler 方法（保留空结构体用于路由注册）
- [ ] 添加注释说明各子模块职责

### 任务 4: 验证编译
- [ ] 执行 `go build -mod=vendor -o dist/switch-admin.exe ./cmd/main.go`
- [ ] 修复编译错误
- [ ] 执行 `go vet ./...` 检查

## 阶段二：DataModel 层拆分

### 任务 5: 创建 DataModel 子模块文件 (1-6)
- [ ] 创建 `internal/datamodel/maintenance_reboot_save.go`
  - [ ] 移动 `GetRebootSaveContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_system_config.go`
  - [ ] 移动 `GetMaintenanceSystemConfigContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_file_management.go`
  - [ ] 移动 `GetFilesContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_log_management.go`
  - [ ] 移动 `GetLogsContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_snmp_config.go`
  - [ ] 移动 `GetSNMPContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_snmp_trap.go`
  - [ ] 移动 `GetSNMPTrapContent` 函数

### 任务 6: 创建 DataModel 子模块文件 (7-12)
- [ ] 创建 `internal/datamodel/maintenance_worm_protection.go`
  - [ ] 移动 `GetWormProtectionContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_ddos_protection.go`
  - [ ] 移动 `GetDDoSProtectionContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_arp_protection.go`
  - [ ] 移动 `GetARPProtectionContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_user_management.go`
  - [ ] 移动 `GetUsersContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_session_management.go`
  - [ ] 移动 `GetSessionsContent` 函数
- [ ] 创建 `internal/datamodel/maintenance_load_config.go`
  - [ ] 移动 `GetLoadConfigContent` 函数

### 任务 7: 重构 maintenance.go
- [ ] 保留主页面框架结构（如果有）
- [ ] 移除已拆分的页面函数
- [ ] 添加注释说明各子模块职责
- [ ] 保留通用组件和辅助函数

### 任务 8: 验证编译
- [ ] 执行 `go build -mod=vendor -o dist/switch-admin.exe ./cmd/main.go`
- [ ] 修复编译错误
- [ ] 执行 `go vet ./...` 检查

## 阶段三：测试与验证

### 任务 9: 启动服务测试
- [ ] 启动应用 `./dist/switch-admin.exe`
- [ ] 访问所有维护模块页面验证 UI
- [ ] 测试所有 API 端点功能

### 任务 10: 代码审查与清理
- [ ] 检查是否有重复代码可提取
- [ ] 统一代码风格和注释
- [ ] 执行 `go fmt ./...`
- [ ] 准备提交代码

## 依赖关系

```
任务 1 → 任务 2 → 任务 3 → 任务 4 (Handler 层完成)
                              ↓
任务 5 → 任务 6 → 任务 7 → 任务 8 (DataModel 层完成)
                              ↓
                    任务 9 → 任务 10 (测试与验证)
```

## 预计工时

- Handler 层拆分：3-4 小时
- DataModel 层拆分：3-4 小时
- 测试与验证：1-2 小时
- **总计：7-10 小时**
