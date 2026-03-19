package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetVLANConfig 获取 VLAN 配置 (CLI)
func (p *CLIProvider) GetVLANConfig(ctx context.Context) (*model.VLANConfigListResponse, error) {
	return &model.VLANConfigListResponse{}, nil
}

// CreateVLAN 创建 VLAN (CLI)
func (p *CLIProvider) CreateVLAN(ctx context.Context, req model.VLANCreateRequest) error {
	return nil
}

// UpdateVLAN 更新 VLAN 配置 (CLI)
func (p *CLIProvider) UpdateVLAN(ctx context.Context, vlanID string, req model.VLANUpdateRequest) error {
	return nil
}

// DeleteVLAN 删除 VLAN (CLI)
func (p *CLIProvider) DeleteVLAN(ctx context.Context, vlanID string) error {
	return nil
}

// GetSTPConfig 获取 STP 配置 (CLI)
func (p *CLIProvider) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	return &model.STPConfig{}, nil
}

// GetSTPStatus 获取 STP 状态 (CLI)
func (p *CLIProvider) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	return &model.STPStatus{}, nil
}

// UpdateSTPConfig 更新 STP 配置 (CLI)
func (p *CLIProvider) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	return nil
}

// GetPoEConfig 获取 PoE 配置 (CLI)
func (p *CLIProvider) GetPoEConfig(ctx context.Context) (*model.PoEConfig, error) {
	return &model.PoEConfig{}, nil
}

// UpdatePoEConfig 更新 PoE 配置 (CLI)
func (p *CLIProvider) UpdatePoEConfig(ctx context.Context, portID string, req model.PoEPortRequest) error {
	return nil
}

// GetStormControl 获取风暴控制配置 (CLI)
func (p *CLIProvider) GetStormControl(ctx context.Context) (*model.StormControlConfig, error) {
	return &model.StormControlConfig{}, nil
}

// UpdateStormControlGlobal 更新全局风暴控制配置 (CLI)
func (p *CLIProvider) UpdateStormControlGlobal(ctx context.Context, req model.StormControlRequest) error {
	return nil
}

// UpdateStormControlPort 更新端口风暴控制配置 (CLI)
func (p *CLIProvider) UpdateStormControlPort(ctx context.Context, portID string, req model.StormControlPortRequest) error {
	return nil
}

// GetFlowControl 获取流控配置 (CLI)
func (p *CLIProvider) GetFlowControl(ctx context.Context) (*model.FlowControlConfig, error) {
	return &model.FlowControlConfig{}, nil
}

// UpdateFlowControlGlobal 更新全局流控配置 (CLI)
func (p *CLIProvider) UpdateFlowControlGlobal(ctx context.Context, req model.FlowControlRequest) error {
	return nil
}

// UpdateFlowControlPort 更新端口流控配置 (CLI)
func (p *CLIProvider) UpdateFlowControlPort(ctx context.Context, portID string, req model.FlowControlPortRequest) error {
	return nil
}

// GetPortIsolation 获取端口隔离配置 (CLI)
func (p *CLIProvider) GetPortIsolation(ctx context.Context) (*model.PortIsolationConfig, error) {
	return &model.PortIsolationConfig{}, nil
}

// UpdatePortIsolation 更新端口隔离配置 (CLI)
func (p *CLIProvider) UpdatePortIsolation(ctx context.Context, req model.PortIsolationRequest) error {
	return nil
}

// DeletePortIsolation 删除端口隔离配置 (CLI)
func (p *CLIProvider) DeletePortIsolation(ctx context.Context, groupID int) error {
	return nil
}

// GetPortMirror 获取端口镜像配置 (CLI)
func (p *CLIProvider) GetPortMirror(ctx context.Context) (*model.PortMonitorConfig, error) {
	return &model.PortMonitorConfig{}, nil
}

// UpdatePortMirror 更新端口镜像配置 (CLI)
func (p *CLIProvider) UpdatePortMirror(ctx context.Context, req model.PortMirrorRequest) error {
	return nil
}

// DeletePortMirror 删除端口镜像配置 (CLI)
func (p *CLIProvider) DeletePortMirror(ctx context.Context, sessionID int) error {
	return nil
}

// GetMacTable 获取 MAC 地址表 (CLI)
func (p *CLIProvider) GetMacTable(ctx context.Context) (*model.MacTableListResponse, error) {
	return &model.MacTableListResponse{}, nil
}

// ClearDynamicMacEntries 清除动态 MAC 表项 (CLI)
func (p *CLIProvider) ClearDynamicMacEntries(ctx context.Context) error {
	return nil
}

// GetERPSConfig 获取 ERPS 配置 (CLI)
func (p *CLIProvider) GetERPSConfig(ctx context.Context) (*model.ERPSConfig, error) {
	return &model.ERPSConfig{}, nil
}

// UpdateERPSConfig 更新 ERPS 配置 (CLI)
func (p *CLIProvider) UpdateERPSConfig(ctx context.Context, req model.ERPSRequest) error {
	return nil
}

// GetMulticastConfig 获取组播配置 (CLI)
func (p *CLIProvider) GetMulticastConfig(ctx context.Context) (*model.MulticastConfig, error) {
	return &model.MulticastConfig{}, nil
}

// UpdateMulticastConfig 更新组播配置 (CLI)
func (p *CLIProvider) UpdateMulticastConfig(ctx context.Context, req model.MulticastRequest) error {
	return nil
}

// GetResource 获取资源使用情况 (CLI)
func (p *CLIProvider) GetResource(ctx context.Context) (*model.ResourceUsage, error) {
	return &model.ResourceUsage{}, nil
}

// GetStackConfig 获取堆叠配置 (CLI)
func (p *CLIProvider) GetStackConfig(ctx context.Context) (*model.StackConfig, error) {
	return &model.StackConfig{}, nil
}

// UpdateStackMember 更新堆叠成员配置 (CLI)
func (p *CLIProvider) UpdateStackMember(ctx context.Context, req model.StackRequest) error {
	return nil
}
