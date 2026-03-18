package service

import (
	"context"
	"sync"

	"switch-admin/internal/model"
	"switch-admin/internal/service/mode"
	"switch-admin/internal/service/provider"
)

// NetworkService 网络服务
type NetworkService struct {
	mu           sync.RWMutex
	modeResolver *mode.ModeResolver
}

var networkService *NetworkService
var networkOnce sync.Once

// GetNetworkService 获取网络服务单例
func GetNetworkService() *NetworkService {
	networkOnce.Do(func() {
		networkService = &NetworkService{
			modeResolver: mode.NewModeResolver(mode.ModeResolverConfig{
				InitialMode: mode.ModeMock,
			}),
		}
	})
	return networkService
}

// SetMode 设置模式（CLI 或 Mock）
func (s *NetworkService) SetMode(m mode.RunMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modeResolver.SwitchMode(m)
}

// getProvider 获取当前模式的 Provider
func (s *NetworkService) getProvider() provider.NetworkProvider {
	return s.modeResolver.GetNetworkProvider()
}

// GetVLANList 获取 VLAN 列表
func (s *NetworkService) GetVLANList(ctx context.Context) (*model.VLANListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetVLANList(ctx)
}

// CreateVLAN 创建 VLAN
func (s *NetworkService) CreateVLAN(ctx context.Context, req model.VLANRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().CreateVLAN(ctx, req)
}

// UpdateVLAN 更新 VLAN
func (s *NetworkService) UpdateVLAN(ctx context.Context, id int, req model.VLANRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateVLAN(ctx, id, req)
}

// DeleteVLAN 删除 VLAN
func (s *NetworkService) DeleteVLAN(ctx context.Context, id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteVLAN(ctx, id)
}

// DeleteVLANs 批量删除 VLAN
func (s *NetworkService) DeleteVLANs(ctx context.Context, ids []int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteVLANs(ctx, ids)
}

// AddVLANPort 添加 VLAN 端口
func (s *NetworkService) AddVLANPort(ctx context.Context, vlanID int, port string, mode string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddVLANPort(ctx, vlanID, port, mode)
}

// RemoveVLANPort 移除 VLAN 端口
func (s *NetworkService) RemoveVLANPort(ctx context.Context, vlanID int, port string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().RemoveVLANPort(ctx, vlanID, port)
}

// GetPortList 获取端口列表
func (s *NetworkService) GetPortList(ctx context.Context) (*model.NetworkPortListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetPortList(ctx)
}

// GetPortDetail 获取端口详情
func (s *NetworkService) GetPortDetail(ctx context.Context, portName string) (*model.PortDetail, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetPortDetail(ctx, portName)
}

// UpdatePort 更新端口配置
func (s *NetworkService) UpdatePort(ctx context.Context, portName string, req model.PortUpdateRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdatePort(ctx, portName, req)
}

// ResetPort 重置端口配置
func (s *NetworkService) ResetPort(ctx context.Context, portName string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().ResetPort(ctx, portName)
}

// RestartPort 重启端口
func (s *NetworkService) RestartPort(ctx context.Context, portName string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().RestartPort(ctx, portName)
}

// GetLAGList 获取链路聚合组列表
func (s *NetworkService) GetLAGList(ctx context.Context) (*model.LAGListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetLAGList(ctx)
}

// CreateLAG 创建链路聚合组
func (s *NetworkService) CreateLAG(ctx context.Context, req model.LAGRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().CreateLAG(ctx, req)
}

// UpdateLAG 更新链路聚合组
func (s *NetworkService) UpdateLAG(ctx context.Context, id int, req model.LAGRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateLAG(ctx, id, req)
}

// DeleteLAG 删除链路聚合组
func (s *NetworkService) DeleteLAG(ctx context.Context, id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteLAG(ctx, id)
}

// AddLAGPort 添加端口到聚合组
func (s *NetworkService) AddLAGPort(ctx context.Context, lagID int, port string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddLAGPort(ctx, lagID, port)
}

// RemoveLAGPort 从聚合组移除端口
func (s *NetworkService) RemoveLAGPort(ctx context.Context, lagID int, port string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().RemoveLAGPort(ctx, lagID, port)
}

// GetSTPConfig 获取 STP 配置
func (s *NetworkService) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSTPConfig(ctx)
}

// UpdateSTPConfig 更新 STP 配置
func (s *NetworkService) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateSTPConfig(ctx, req)
}

// GetSTPStatus 获取 STP 状态
func (s *NetworkService) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetSTPStatus(ctx)
}

// GetACLList 获取 ACL 列表
func (s *NetworkService) GetACLList(ctx context.Context) (*model.ACLListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetACLList(ctx)
}

// CreateACL 创建 ACL
func (s *NetworkService) CreateACL(ctx context.Context, req model.ACLRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().CreateACL(ctx, req)
}

// UpdateACL 更新 ACL
func (s *NetworkService) UpdateACL(ctx context.Context, id int, req model.ACLRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateACL(ctx, id, req)
}

// DeleteACL 删除 ACL
func (s *NetworkService) DeleteACL(ctx context.Context, id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteACL(ctx, id)
}

// GetACLRules 获取 ACL 规则列表
func (s *NetworkService) GetACLRules(ctx context.Context, aclID int) (*model.ACLRuleListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetACLRules(ctx, aclID)
}

// AddACLRule 添加 ACL 规则
func (s *NetworkService) AddACLRule(ctx context.Context, aclID int, req model.ACLRuleRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().AddACLRule(ctx, aclID, req)
}

// UpdateACLRule 更新 ACL 规则
func (s *NetworkService) UpdateACLRule(ctx context.Context, aclID int, ruleID int, req model.ACLRuleRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateACLRule(ctx, aclID, ruleID, req)
}

// DeleteACLRule 删除 ACL 规则
func (s *NetworkService) DeleteACLRule(ctx context.Context, aclID int, ruleID int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteACLRule(ctx, aclID, ruleID)
}
