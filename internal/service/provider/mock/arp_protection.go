package mock

import (
	"context"
	"time"

	"switch-admin/internal/model"
)

// GetARPConfig 获取 ARP 防护配置
func (p *MaintenanceProvider) GetARPConfig(ctx context.Context) (*model.ARPConfig, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.ARPConfig{
		Enabled:        true,
		InspectEnabled: true,
		TrustPorts:     []string{"eth0/1", "eth0/2"},
	}, nil
}

// UpdateARPConfig 更新 ARP 防护配置
func (p *MaintenanceProvider) UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
