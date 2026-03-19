package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetVLANConfig 获取 VLAN 配置 (Mock)
func (p *MockProvider) GetVLANConfig(ctx context.Context) (*model.VLANConfigListResponse, error) {
	return &model.VLANConfigListResponse{}, nil
}

// CreateVLAN 创建 VLAN (Mock)
func (p *MockProvider) CreateVLAN(ctx context.Context, req model.VLANCreateRequest) error {
	return nil
}

// UpdateVLAN 更新 VLAN 配置 (Mock)
func (p *MockProvider) UpdateVLAN(ctx context.Context, vlanID string, req model.VLANUpdateRequest) error {
	return nil
}

// DeleteVLAN 删除 VLAN (Mock)
func (p *MockProvider) DeleteVLAN(ctx context.Context, vlanID string) error {
	return nil
}

// GetSTPConfig 获取 STP 配置 (Mock)
func (p *MockProvider) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	return &model.STPConfig{}, nil
}

// GetSTPStatus 获取 STP 状态 (Mock)
func (p *MockProvider) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	return &model.STPStatus{}, nil
}

// UpdateSTPConfig 更新 STP 配置 (Mock)
func (p *MockProvider) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	return nil
}

// GetPoEConfig 获取 PoE 配置 (Mock)
func (p *MockProvider) GetPoEConfig(ctx context.Context) (*model.PoEConfig, error) {
	return &model.PoEConfig{}, nil
}

// UpdatePoEConfig 更新 PoE 配置 (Mock)
func (p *MockProvider) UpdatePoEConfig(ctx context.Context, portID string, req model.PoEPortRequest) error {
	return nil
}
