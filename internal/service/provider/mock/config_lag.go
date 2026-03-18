package mock

import (
	"context"
	"switch-admin/internal/model"
	"time"
)

// GetLinkAggregationList 获取链路聚合列表（Mock 数据）
func (p *ConfigProvider) GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.LinkAggregationListResponse{
		Aggregations: []model.LinkAggregation{
			{
				GroupID:     1,
				Name:        "Ag1",
				Mode:        "LACP",
				LoadBalance: "src-dst-ip",
				MemberPorts: []string{"GE1/0/1", "GE1/0/2"},
				MinActive:   1,
				Status:      "normal",
			},
			{
				GroupID:     2,
				Name:        "Ag2",
				Mode:        "Static",
				LoadBalance: "src-dst-mac",
				MemberPorts: []string{"GE1/0/5", "GE1/0/6"},
				MinActive:   1,
				Status:      "normal",
			},
			{
				GroupID:     3,
				Name:        "Ag3",
				Mode:        "LACP",
				LoadBalance: "src-dst-mac",
				MemberPorts: []string{"GE1/0/9", "GE1/0/10", "GE1/0/11"},
				MinActive:   2,
				Status:      "degraded",
			},
		},
		Total: 3,
	}, nil
}

// CreateLinkAggregation 创建链路聚合组（Mock）
func (p *ConfigProvider) CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}

// UpdateLinkAggregation 更新链路聚合组（Mock）
func (p *ConfigProvider) UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}

// DeleteLinkAggregation 删除链路聚合组（Mock）
func (p *ConfigProvider) DeleteLinkAggregation(ctx context.Context, id int) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}
