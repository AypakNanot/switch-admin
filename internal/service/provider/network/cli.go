package network

import (
	"context"

	"switch-admin/internal/model"
)

// GetVLANList 获取 VLAN 列表 (CLI)
func (p *CLIProvider) GetVLANList(ctx context.Context) (*model.VLANListResponse, error) {
	return &model.VLANListResponse{}, nil
}

// CreateVLAN 创建 VLAN (CLI)
func (p *CLIProvider) CreateVLAN(ctx context.Context, req model.VLANRequest) error {
	return nil
}

// UpdateVLAN 更新 VLAN (CLI)
func (p *CLIProvider) UpdateVLAN(ctx context.Context, id int, req model.VLANRequest) error {
	return nil
}

// DeleteVLAN 删除 VLAN (CLI)
func (p *CLIProvider) DeleteVLAN(ctx context.Context, id int) error {
	return nil
}

// DeleteVLANs 批量删除 VLAN (CLI)
func (p *CLIProvider) DeleteVLANs(ctx context.Context, ids []int) error {
	return nil
}

// AddVLANPort 添加 VLAN 端口 (CLI)
func (p *CLIProvider) AddVLANPort(ctx context.Context, vlanID int, port string, mode string) error {
	return nil
}

// RemoveVLANPort 移除 VLAN 端口 (CLI)
func (p *CLIProvider) RemoveVLANPort(ctx context.Context, vlanID int, port string) error {
	return nil
}

// GetPortList 获取端口列表 (CLI)
func (p *CLIProvider) GetPortList(ctx context.Context) (*model.NetworkPortListResponse, error) {
	return &model.NetworkPortListResponse{}, nil
}

// GetPortDetail 获取端口详情 (CLI)
func (p *CLIProvider) GetPortDetail(ctx context.Context, portName string) (*model.PortDetail, error) {
	return &model.PortDetail{}, nil
}

// UpdatePort 更新端口 (CLI)
func (p *CLIProvider) UpdatePort(ctx context.Context, portName string, req model.PortUpdateRequest) error {
	return nil
}

// ResetPort 重置端口 (CLI)
func (p *CLIProvider) ResetPort(ctx context.Context, portName string) error {
	return nil
}

// RestartPort 重启端口 (CLI)
func (p *CLIProvider) RestartPort(ctx context.Context, portName string) error {
	return nil
}

// GetLAGList 获取链路聚合列表 (CLI)
func (p *CLIProvider) GetLAGList(ctx context.Context) (*model.LAGListResponse, error) {
	return &model.LAGListResponse{}, nil
}

// CreateLAG 创建链路聚合 (CLI)
func (p *CLIProvider) CreateLAG(ctx context.Context, req model.LAGRequest) error {
	return nil
}

// UpdateLAG 更新链路聚合 (CLI)
func (p *CLIProvider) UpdateLAG(ctx context.Context, id int, req model.LAGRequest) error {
	return nil
}

// DeleteLAG 删除链路聚合 (CLI)
func (p *CLIProvider) DeleteLAG(ctx context.Context, id int) error {
	return nil
}

// AddLAGPort 添加链路聚合端口 (CLI)
func (p *CLIProvider) AddLAGPort(ctx context.Context, lagID int, port string) error {
	return nil
}

// RemoveLAGPort 移除链路聚合端口 (CLI)
func (p *CLIProvider) RemoveLAGPort(ctx context.Context, lagID int, port string) error {
	return nil
}

// GetSTPConfig 获取 STP 配置 (CLI)
func (p *CLIProvider) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	return &model.STPConfig{}, nil
}

// UpdateSTPConfig 更新 STP 配置 (CLI)
func (p *CLIProvider) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	return nil
}

// GetSTPStatus 获取 STP 状态 (CLI)
func (p *CLIProvider) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	return &model.STPStatus{}, nil
}

// GetACLList 获取 ACL 列表 (CLI)
func (p *CLIProvider) GetACLList(ctx context.Context) (*model.ACLListResponse, error) {
	return &model.ACLListResponse{}, nil
}

// CreateACL 创建 ACL (CLI)
func (p *CLIProvider) CreateACL(ctx context.Context, req model.ACLRequest) error {
	return nil
}

// UpdateACL 更新 ACL (CLI)
func (p *CLIProvider) UpdateACL(ctx context.Context, id int, req model.ACLRequest) error {
	return nil
}

// DeleteACL 删除 ACL (CLI)
func (p *CLIProvider) DeleteACL(ctx context.Context, id int) error {
	return nil
}

// GetACLRules 获取 ACL 规则列表 (CLI)
func (p *CLIProvider) GetACLRules(ctx context.Context, aclID int) (*model.ACLRuleListResponse, error) {
	return &model.ACLRuleListResponse{}, nil
}

// AddACLRule 添加 ACL 规则 (CLI)
func (p *CLIProvider) AddACLRule(ctx context.Context, aclID int, req model.ACLRuleRequest) error {
	return nil
}

// UpdateACLRule 更新 ACL 规则 (CLI)
func (p *CLIProvider) UpdateACLRule(ctx context.Context, aclID int, ruleID int, req model.ACLRuleRequest) error {
	return nil
}

// DeleteACLRule 删除 ACL 规则 (CLI)
func (p *CLIProvider) DeleteACLRule(ctx context.Context, aclID int, ruleID int) error {
	return nil
}
