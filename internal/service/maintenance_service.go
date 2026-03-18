package service

import (
	"context"
	"sync"

	"switch-admin/internal/model"
	"switch-admin/internal/service/mode"
	"switch-admin/internal/service/provider"
)

// MaintenanceService 维护服务
type MaintenanceService struct {
	mu           sync.RWMutex
	modeResolver *mode.ModeResolver
}

var maintenanceService *MaintenanceService
var maintenanceOnce sync.Once

// GetMaintenanceService 获取维护服务单例
func GetMaintenanceService() *MaintenanceService {
	maintenanceOnce.Do(func() {
		maintenanceService = &MaintenanceService{
			modeResolver: mode.NewModeResolver(mode.ModeResolverConfig{
				InitialMode: mode.ModeMock,
			}),
		}
	})
	return maintenanceService
}

// SetMode 设置模式（CLI 或 Mock）
func (s *MaintenanceService) SetMode(m mode.RunMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modeResolver.SwitchMode(m)
}

// getProvider 获取当前模式的 Provider
func (s *MaintenanceService) getProvider() provider.MaintenanceProvider {
	return s.modeResolver.GetMaintenanceProvider()
}

// GetSystemConfig 获取系统配置
func (s *MaintenanceService) GetSystemConfig(ctx context.Context) (*model.SystemConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSystemConfig(ctx)
}

// UpdateNetworkConfig 更新网络配置
func (s *MaintenanceService) UpdateNetworkConfig(ctx context.Context, req model.NetworkConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateNetworkConfig(ctx, req)
}

// UpdateTemperatureConfig 更新温度配置
func (s *MaintenanceService) UpdateTemperatureConfig(ctx context.Context, req model.TemperatureConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateTemperatureConfig(ctx, req)
}

// UpdateDeviceInfo 更新设备信息
func (s *MaintenanceService) UpdateDeviceInfo(ctx context.Context, req model.DeviceInfoRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateDeviceInfo(ctx, req)
}

// UpdateDateTime 更新时间日期
func (s *MaintenanceService) UpdateDateTime(ctx context.Context, req model.DateTimeRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateDateTime(ctx, req)
}

// SaveConfig 保存配置
func (s *MaintenanceService) SaveConfig(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().SaveConfig(ctx)
}

// RebootSwitch 重启交换机
func (s *MaintenanceService) RebootSwitch(ctx context.Context, delay int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().RebootSwitch(ctx, delay)
}

// FactoryReset 恢复出厂配置
func (s *MaintenanceService) FactoryReset(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().FactoryReset(ctx)
}

// GetUsers 获取用户列表
func (s *MaintenanceService) GetUsers(ctx context.Context) (*model.UserListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetUsers(ctx)
}

// CreateUser 创建用户
func (s *MaintenanceService) CreateUser(ctx context.Context, req model.UserRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().CreateUser(ctx, req)
}

// DeleteUser 删除用户
func (s *MaintenanceService) DeleteUser(ctx context.Context, username string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteUser(ctx, username)
}

// DeleteUsers 批量删除用户
func (s *MaintenanceService) DeleteUsers(ctx context.Context, usernames []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteUsers(ctx, usernames)
}

// GetSessions 获取会话列表
func (s *MaintenanceService) GetSessions(ctx context.Context) (*model.SessionListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSessions(ctx)
}

// DeleteSession 删除会话
func (s *MaintenanceService) DeleteSession(ctx context.Context, sessionID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteSession(ctx, sessionID)
}

// DeleteSessions 批量删除会话
func (s *MaintenanceService) DeleteSessions(ctx context.Context, sessionIDs []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteSessions(ctx, sessionIDs)
}

// GetLogs 获取日志列表
func (s *MaintenanceService) GetLogs(ctx context.Context) (*model.LogListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetLogs(ctx)
}

// ClearLogs 清除日志
func (s *MaintenanceService) ClearLogs(ctx context.Context, levels []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().ClearLogs(ctx, levels)
}

// GetFiles 获取文件列表
func (s *MaintenanceService) GetFiles(ctx context.Context, path string) (*model.FileListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetFiles(ctx, path)
}

// UploadFile 上传文件
func (s *MaintenanceService) UploadFile(ctx context.Context, req model.FileUploadRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UploadFile(ctx, req)
}

// DeleteFile 删除文件
func (s *MaintenanceService) DeleteFile(ctx context.Context, path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteFile(ctx, path)
}

// DeleteFiles 批量删除文件
func (s *MaintenanceService) DeleteFiles(ctx context.Context, paths []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteFiles(ctx, paths)
}

// DownloadFile 下载文件
func (s *MaintenanceService) DownloadFile(ctx context.Context, path string) ([]byte, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DownloadFile(ctx, path)
}

// GetSNMPConfig 获取 SNMP 配置
func (s *MaintenanceService) GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSNMPConfig(ctx)
}

// UpdateSNMPConfig 更新 SNMP 配置
func (s *MaintenanceService) UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateSNMPConfig(ctx, req)
}

// GetTrapHosts 获取 Trap 主机列表
func (s *MaintenanceService) GetTrapHosts(ctx context.Context) ([]model.TrapHost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetTrapHosts(ctx)
}

// AddTrapHost 添加 Trap 主机
func (s *MaintenanceService) AddTrapHost(ctx context.Context, req model.TrapHostRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddTrapHost(ctx, req)
}

// DeleteTrapHost 删除 Trap 主机
func (s *MaintenanceService) DeleteTrapHost(ctx context.Context, host string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteTrapHost(ctx, host)
}

// TestTrap 测试 Trap 发送
func (s *MaintenanceService) TestTrap(ctx context.Context, host string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().TestTrap(ctx, host)
}

// GetSNMPCommunities 获取 SNMP 团体列表
func (s *MaintenanceService) GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSNMPCommunities(ctx)
}

// AddCommunity 添加团体
func (s *MaintenanceService) AddCommunity(ctx context.Context, name, access, description string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddCommunity(ctx, name, access, description)
}

// DeleteCommunity 删除团体
func (s *MaintenanceService) DeleteCommunity(ctx context.Context, name string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteCommunity(ctx, name)
}

// GetWormRules 获取蠕虫规则列表
func (s *MaintenanceService) GetWormRules(ctx context.Context) (*model.WormRuleList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetWormRules(ctx)
}

// AddWormRule 添加蠕虫规则
func (s *MaintenanceService) AddWormRule(ctx context.Context, req model.WormRuleRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddWormRule(ctx, req)
}

// UpdateWormRule 更新蠕虫规则
func (s *MaintenanceService) UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateWormRule(ctx, id, req)
}

// DeleteWormRule 删除蠕虫规则
func (s *MaintenanceService) DeleteWormRule(ctx context.Context, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteWormRule(ctx, id)
}

// DeleteWormRules 批量删除蠕虫规则
func (s *MaintenanceService) DeleteWormRules(ctx context.Context, ids []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteWormRules(ctx, ids)
}

// ClearWormStats 清除蠕虫统计
func (s *MaintenanceService) ClearWormStats(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().ClearWormStats(ctx)
}

// GetDDoSConfig 获取 DDoS 防护配置
func (s *MaintenanceService) GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetDDoSConfig(ctx)
}

// UpdateDDoSConfig 更新 DDoS 防护配置
func (s *MaintenanceService) UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateDDoSConfig(ctx, req)
}

// GetARPConfig 获取 ARP 防护配置
func (s *MaintenanceService) GetARPConfig(ctx context.Context) (*model.ARPConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetARPConfig(ctx)
}

// UpdateARPConfig 更新 ARP 防护配置
func (s *MaintenanceService) UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateARPConfig(ctx, req)
}

// GetConfigFiles 获取配置文件列表
func (s *MaintenanceService) GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetConfigFiles(ctx)
}

// LoadConfig 加载配置
func (s *MaintenanceService) LoadConfig(ctx context.Context, configFile string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().LoadConfig(ctx, configFile)
}
