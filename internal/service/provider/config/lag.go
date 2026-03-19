package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetLinkAggregationList 获取链路聚合列表
func (p *CLIProvider) GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取链路聚合组列表
	// 示例命令：show link-aggregation verbose

	return &model.LinkAggregationListResponse{
		Aggregations: []model.LinkAggregation{},
		Total:        0,
	}, nil
}

// CreateLinkAggregation 创建链路聚合组
func (p *CLIProvider) CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error {
	// TODO: 实现调用交换机 CLI 创建链路聚合组
	// 示例命令序列：
	// configure terminal
	// link-aggregation group <group-id> mode <mode>
	// interface range <ports>
	// port link-aggregation group <group-id>

	return nil
}

// UpdateLinkAggregation 更新链路聚合组
func (p *CLIProvider) UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error {
	// TODO: 实现调用交换机 CLI 更新链路聚合组配置
	// 示例命令序列：
	// configure terminal
	// link-aggregation group <id> load-balance <mode>
	// interface range <ports>
	// port link-aggregation group <id>

	return nil
}

// DeleteLinkAggregation 删除链路聚合组
func (p *CLIProvider) DeleteLinkAggregation(ctx context.Context, id int) error {
	// TODO: 实现调用交换机 CLI 删除链路聚合组
	// 示例命令序列：
	// configure terminal
	// undo link-aggregation group <id>

	return nil
}
