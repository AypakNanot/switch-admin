package maintenance

import (
	"context"

	"switch-admin/internal/model"
)

// GetFiles 获取文件列表 (CLI)
func (p *CLIProvider) GetFiles(ctx context.Context, path string) (*model.FileListResponse, error) {
	return &model.FileListResponse{}, nil
}

// UploadFile 上传文件 (CLI)
func (p *CLIProvider) UploadFile(ctx context.Context, req model.FileUploadRequest) error {
	return nil
}

// DeleteFile 删除文件 (CLI)
func (p *CLIProvider) DeleteFile(ctx context.Context, path string) error {
	return nil
}

// DeleteFiles 批量删除文件 (CLI)
func (p *CLIProvider) DeleteFiles(ctx context.Context, paths []string) error {
	return nil
}

// DownloadFile 下载文件 (CLI)
func (p *CLIProvider) DownloadFile(ctx context.Context, path string) ([]byte, string, error) {
	return nil, "", nil
}

// GetSNMPConfig 获取 SNMP 配置 (CLI)
func (p *CLIProvider) GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error) {
	return &model.SNMPConfig{}, nil
}

// UpdateSNMPConfig 更新 SNMP 配置 (CLI)
func (p *CLIProvider) UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error {
	return nil
}

// GetTrapHosts 获取 Trap 主机列表 (CLI)
func (p *CLIProvider) GetTrapHosts(ctx context.Context) ([]model.TrapHost, error) {
	return []model.TrapHost{}, nil
}

// AddTrapHost 添加 Trap 主机 (CLI)
func (p *CLIProvider) AddTrapHost(ctx context.Context, req model.TrapHostRequest) error {
	return nil
}

// DeleteTrapHost 删除 Trap 主机 (CLI)
func (p *CLIProvider) DeleteTrapHost(ctx context.Context, host string) error {
	return nil
}

// TestTrap 测试 Trap (CLI)
func (p *CLIProvider) TestTrap(ctx context.Context, host string) error {
	return nil
}

// GetSNMPCommunities 获取 SNMP 团体 (CLI)
func (p *CLIProvider) GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error) {
	return []model.SNMPCommunity{}, nil
}

// AddCommunity 添加 SNMP 团体 (CLI)
func (p *CLIProvider) AddCommunity(ctx context.Context, name, access, description string) error {
	return nil
}

// DeleteCommunity 删除 SNMP 团体 (CLI)
func (p *CLIProvider) DeleteCommunity(ctx context.Context, name string) error {
	return nil
}

// GetWormRules 获取蠕虫规则列表 (CLI)
func (p *CLIProvider) GetWormRules(ctx context.Context) (*model.WormRuleList, error) {
	return &model.WormRuleList{}, nil
}

// AddWormRule 添加蠕虫规则 (CLI)
func (p *CLIProvider) AddWormRule(ctx context.Context, req model.WormRuleRequest) error {
	return nil
}

// UpdateWormRule 更新蠕虫规则 (CLI)
func (p *CLIProvider) UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error {
	return nil
}

// DeleteWormRule 删除蠕虫规则 (CLI)
func (p *CLIProvider) DeleteWormRule(ctx context.Context, id string) error {
	return nil
}

// DeleteWormRules 批量删除蠕虫规则 (CLI)
func (p *CLIProvider) DeleteWormRules(ctx context.Context, ids []string) error {
	return nil
}

// ClearWormStats 清除蠕虫统计 (CLI)
func (p *CLIProvider) ClearWormStats(ctx context.Context) error {
	return nil
}

// GetDDoSConfig 获取 DDoS 配置 (CLI)
func (p *CLIProvider) GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error) {
	return &model.DDoSConfig{}, nil
}

// UpdateDDoSConfig 更新 DDoS 配置 (CLI)
func (p *CLIProvider) UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error {
	return nil
}

// GetARPConfig 获取 ARP 配置 (CLI)
func (p *CLIProvider) GetARPConfig(ctx context.Context) (*model.ARPConfig, error) {
	return &model.ARPConfig{}, nil
}

// UpdateARPConfig 更新 ARP 配置 (CLI)
func (p *CLIProvider) UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error {
	return nil
}

// GetConfigFiles 获取配置文件列表 (CLI)
func (p *CLIProvider) GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error) {
	return &model.LoadConfigListResponse{}, nil
}

// LoadConfig 加载配置 (CLI)
func (p *CLIProvider) LoadConfig(ctx context.Context, configFile string) error {
	return nil
}
