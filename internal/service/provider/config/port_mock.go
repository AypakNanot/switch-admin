package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetPortList 获取端口配置列表
func (p *MockProvider) GetPortList(ctx context.Context) (*model.PortConfigListResponse, error) {
	mockPorts, err := p.mockDataDAO.GetAllMockPorts()
	if err != nil {
		return nil, err
	}

	var ports []model.PortConfig
	for _, mockPort := range mockPorts {
		ports = append(ports, *p.mockPortToPortConfig(&mockPort))
	}

	return &model.PortConfigListResponse{
		Ports: ports,
		Total: len(ports),
	}, nil
}

// GetPortDetail 获取端口详情
func (p *MockProvider) GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error) {
	mockPort, err := p.mockDataDAO.GetMockPortByName(portID)
	if err != nil {
		return nil, err
	}

	return p.mockPortToPortConfig(mockPort), nil
}

// UpdatePort 更新端口配置
func (p *MockProvider) UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error {
	// TODO: 实现 Mock 数据更新
	return nil
}

// GetLinkAggregationList 获取链路聚合列表
func (p *MockProvider) GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error) {
	// Mock 数据
	return &model.LinkAggregationListResponse{
		Aggregations: []model.LinkAggregation{},
		Total:        0,
	}, nil
}

// CreateLinkAggregation 创建链路聚合组
func (p *MockProvider) CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error {
	return nil
}

// UpdateLinkAggregation 更新链路聚合组
func (p *MockProvider) UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error {
	return nil
}

// DeleteLinkAggregation 删除链路聚合组
func (p *MockProvider) DeleteLinkAggregation(ctx context.Context, id int) error {
	return nil
}
