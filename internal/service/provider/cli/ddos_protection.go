package cli

import (
	"context"

	"switch-admin/internal/model"
)

// GetDDoSConfig 获取 DDoS 防护配置
func (p *MaintenanceProvider) GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error) {
	// TODO: 实现调用交换机 CLI 获取 DDoS 配置
	return &model.DDoSConfig{
		Enabled:   true,
		Threshold: 1000,
		Action:    "drop",
	}, nil
}

// UpdateDDoSConfig 更新 DDoS 防护配置
func (p *MaintenanceProvider) UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error {
	// TODO: 实现调用交换机 CLI 更新 DDoS 配置
	return nil
}
