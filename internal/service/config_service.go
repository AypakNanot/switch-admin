package service

import (
	"context"
	"sync"

	"switch-admin/internal/model"
	"switch-admin/internal/service/mode"
	"switch-admin/internal/service/provider"
)

// ConfigService 配置服务
type ConfigService struct {
	mu           sync.RWMutex
	modeResolver *mode.ModeResolver
}

var configService *ConfigService
var configOnce sync.Once

// GetConfigService 获取配置服务单例
func GetConfigService() *ConfigService {
	configOnce.Do(func() {
		configService = &ConfigService{
			modeResolver: mode.NewModeResolver(mode.ModeResolverConfig{
				InitialMode: mode.ModeMock,
			}),
		}
	})
	return configService
}

// SetMode 设置模式（CLI 或 Mock）
func (s *ConfigService) SetMode(m mode.RunMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modeResolver.SwitchMode(m)
}

// getProvider 获取当前模式的 Provider
func (s *ConfigService) getProvider() provider.ConfigProvider {
	return s.modeResolver.GetConfigProvider()
}

// GetPortList 获取端口列表
func (s *ConfigService) GetPortList(ctx context.Context) (*model.PortConfigListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetPortList(ctx)
}

// GetPortDetail 获取端口详情
func (s *ConfigService) GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetPortDetail(ctx, portID)
}

// UpdatePort 更新端口配置
func (s *ConfigService) UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdatePort(ctx, portID, req)
}

// GetLinkAggregationList 获取链路聚合列表
func (s *ConfigService) GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().GetLinkAggregationList(ctx)
}

// CreateLinkAggregation 创建链路聚合组
func (s *ConfigService) CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().CreateLinkAggregation(ctx, req)
}

// UpdateLinkAggregation 更新链路聚合组
func (s *ConfigService) UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().UpdateLinkAggregation(ctx, id, req)
}

// DeleteLinkAggregation 删除链路聚合组
func (s *ConfigService) DeleteLinkAggregation(ctx context.Context, id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProvider().DeleteLinkAggregation(ctx, id)
}
