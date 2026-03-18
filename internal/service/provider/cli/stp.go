package cli

import (
	"context"

	"switch-admin/internal/model"
)

// GetSTPConfig 获取 STP 配置
func (p *NetworkProvider) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	// TODO: 实现调用交换机 CLI 获取 STP 配置

	return &model.STPConfig{
		Enabled: true,
		Mode:    "RSTP",
		Priority: 32768,
		RootBridge: "00:11:22:33:44:55",
	}, nil
}

// UpdateSTPConfig 更新 STP 配置
func (p *NetworkProvider) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	// TODO: 实现调用交换机 CLI 更新 STP 配置
	return nil
}

// GetSTPStatus 获取 STP 状态
func (p *NetworkProvider) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	// TODO: 实现调用交换机 CLI 获取 STP 状态

	return &model.STPStatus{
		Enabled: true,
		Mode:    "RSTP",
		RootBridge: "00:11:22:33:44:55",
		RootPort: "eth0/24",
		PortStates: []model.STPPortState{
			{Port: "eth0/1", State: "forwarding", Role: "designated"},
			{Port: "eth0/24", State: "forwarding", Role: "root"},
		},
	}, nil
}
