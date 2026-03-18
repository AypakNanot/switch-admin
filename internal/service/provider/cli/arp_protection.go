package cli

import (
	"context"

	"switch-admin/internal/model"
)

// GetARPConfig 获取 ARP 防护配置
func (p *MaintenanceProvider) GetARPConfig(ctx context.Context) (*model.ARPConfig, error) {
	// TODO: 实现调用交换机 CLI 获取 ARP 配置
	return &model.ARPConfig{
		Enabled:        true,
		InspectEnabled: true,
		TrustPorts:     []string{"eth0/1", "eth0/2"},
	}, nil
}

// UpdateARPConfig 更新 ARP 防护配置
func (p *MaintenanceProvider) UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error {
	// TODO: 实现调用交换机 CLI 更新 ARP 配置
	return nil
}
