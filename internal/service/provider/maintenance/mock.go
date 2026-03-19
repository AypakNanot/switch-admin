package maintenance

import (
	"context"

	"switch-admin/internal/model"
)

// GetSystemConfig 获取系统配置 (Mock)
func (p *MockProvider) GetSystemConfig(ctx context.Context) (*model.SystemConfig, error) {
	return &model.SystemConfig{}, nil
}

// UpdateNetworkConfig 更新网络配置 (Mock)
func (p *MockProvider) UpdateNetworkConfig(ctx context.Context, req model.NetworkConfigRequest) error {
	return nil
}

// UpdateTemperatureConfig 更新温度配置 (Mock)
func (p *MockProvider) UpdateTemperatureConfig(ctx context.Context, req model.TemperatureConfigRequest) error {
	return nil
}

// UpdateDeviceInfo 更新设备信息 (Mock)
func (p *MockProvider) UpdateDeviceInfo(ctx context.Context, req model.DeviceInfoRequest) error {
	return nil
}

// UpdateDateTime 更新时间日期 (Mock)
func (p *MockProvider) UpdateDateTime(ctx context.Context, req model.DateTimeRequest) error {
	return nil
}

// SaveConfig 保存配置 (Mock)
func (p *MockProvider) SaveConfig(ctx context.Context) error {
	return nil
}

// RebootSwitch 重启交换机 (Mock)
func (p *MockProvider) RebootSwitch(ctx context.Context, delay int) error {
	return nil
}

// FactoryReset 恢复出厂设置 (Mock)
func (p *MockProvider) FactoryReset(ctx context.Context) error {
	return nil
}

// GetUsers 获取用户列表 (Mock)
func (p *MockProvider) GetUsers(ctx context.Context) (*model.UserListResponse, error) {
	return &model.UserListResponse{}, nil
}

// CreateUser 创建用户 (Mock)
func (p *MockProvider) CreateUser(ctx context.Context, req model.UserRequest) error {
	return nil
}

// DeleteUser 删除用户 (Mock)
func (p *MockProvider) DeleteUser(ctx context.Context, username string) error {
	return nil
}

// DeleteUsers 批量删除用户 (Mock)
func (p *MockProvider) DeleteUsers(ctx context.Context, usernames []string) error {
	return nil
}

// GetSessions 获取会话列表 (Mock)
func (p *MockProvider) GetSessions(ctx context.Context) (*model.SessionListResponse, error) {
	return &model.SessionListResponse{}, nil
}

// DeleteSession 删除会话 (Mock)
func (p *MockProvider) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

// DeleteSessions 批量删除会话 (Mock)
func (p *MockProvider) DeleteSessions(ctx context.Context, sessionIDs []string) error {
	return nil
}

// GetLogs 获取日志列表 (Mock)
func (p *MockProvider) GetLogs(ctx context.Context) (*model.LogListResponse, error) {
	return &model.LogListResponse{}, nil
}

// ClearLogs 清除日志 (Mock)
func (p *MockProvider) ClearLogs(ctx context.Context, levels []string) error {
	return nil
}

// GetFiles 获取文件列表 (Mock)
func (p *MockProvider) GetFiles(ctx context.Context, path string) (*model.FileListResponse, error) {
	return &model.FileListResponse{}, nil
}

// UploadFile 上传文件 (Mock)
func (p *MockProvider) UploadFile(ctx context.Context, req model.FileUploadRequest) error {
	return nil
}

// DeleteFile 删除文件 (Mock)
func (p *MockProvider) DeleteFile(ctx context.Context, path string) error {
	return nil
}

// DeleteFiles 批量删除文件 (Mock)
func (p *MockProvider) DeleteFiles(ctx context.Context, paths []string) error {
	return nil
}

// DownloadFile 下载文件 (Mock)
func (p *MockProvider) DownloadFile(ctx context.Context, path string) ([]byte, string, error) {
	return nil, "", nil
}

// GetSNMPConfig 获取 SNMP 配置 (Mock)
func (p *MockProvider) GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error) {
	return &model.SNMPConfig{}, nil
}

// UpdateSNMPConfig 更新 SNMP 配置 (Mock)
func (p *MockProvider) UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error {
	return nil
}

// GetTrapHosts 获取 Trap 主机列表 (Mock)
func (p *MockProvider) GetTrapHosts(ctx context.Context) ([]model.TrapHost, error) {
	return []model.TrapHost{}, nil
}

// AddTrapHost 添加 Trap 主机 (Mock)
func (p *MockProvider) AddTrapHost(ctx context.Context, req model.TrapHostRequest) error {
	return nil
}

// DeleteTrapHost 删除 Trap 主机 (Mock)
func (p *MockProvider) DeleteTrapHost(ctx context.Context, host string) error {
	return nil
}

// TestTrap 测试 Trap (Mock)
func (p *MockProvider) TestTrap(ctx context.Context, host string) error {
	return nil
}

// GetSNMPCommunities 获取 SNMP 团体 (Mock)
func (p *MockProvider) GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error) {
	return []model.SNMPCommunity{}, nil
}

// AddCommunity 添加 SNMP 团体 (Mock)
func (p *MockProvider) AddCommunity(ctx context.Context, name, access, description string) error {
	return nil
}

// DeleteCommunity 删除 SNMP 团体 (Mock)
func (p *MockProvider) DeleteCommunity(ctx context.Context, name string) error {
	return nil
}

// GetWormRules 获取蠕虫规则列表 (Mock)
func (p *MockProvider) GetWormRules(ctx context.Context) (*model.WormRuleList, error) {
	return &model.WormRuleList{}, nil
}

// AddWormRule 添加蠕虫规则 (Mock)
func (p *MockProvider) AddWormRule(ctx context.Context, req model.WormRuleRequest) error {
	return nil
}

// UpdateWormRule 更新蠕虫规则 (Mock)
func (p *MockProvider) UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error {
	return nil
}

// DeleteWormRule 删除蠕虫规则 (Mock)
func (p *MockProvider) DeleteWormRule(ctx context.Context, id string) error {
	return nil
}

// DeleteWormRules 批量删除蠕虫规则 (Mock)
func (p *MockProvider) DeleteWormRules(ctx context.Context, ids []string) error {
	return nil
}

// ClearWormStats 清除蠕虫统计 (Mock)
func (p *MockProvider) ClearWormStats(ctx context.Context) error {
	return nil
}

// GetDDoSConfig 获取 DDoS 配置 (Mock)
func (p *MockProvider) GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error) {
	return &model.DDoSConfig{}, nil
}

// UpdateDDoSConfig 更新 DDoS 配置 (Mock)
func (p *MockProvider) UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error {
	return nil
}

// GetARPConfig 获取 ARP 配置 (Mock)
func (p *MockProvider) GetARPConfig(ctx context.Context) (*model.ARPConfig, error) {
	return &model.ARPConfig{}, nil
}

// UpdateARPConfig 更新 ARP 配置 (Mock)
func (p *MockProvider) UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error {
	return nil
}

// GetConfigFiles 获取配置文件列表 (Mock)
func (p *MockProvider) GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error) {
	return &model.LoadConfigListResponse{}, nil
}

// LoadConfig 加载配置 (Mock)
func (p *MockProvider) LoadConfig(ctx context.Context, configFile string) error {
	return nil
}
