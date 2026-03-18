package mock

import (
	"context"
	"time"

	"switch-admin/internal/model"
)

// GetDDoSConfig 获取 DDoS 防护配置
func (p *MaintenanceProvider) GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.DDoSConfig{
		Enabled:   true,
		Threshold: 1000,
		Action:    "drop",
	}, nil
}

// UpdateDDoSConfig 更新 DDoS 防护配置
func (p *MaintenanceProvider) UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
