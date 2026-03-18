package cli

import (
	"context"
	"fmt"

	"switch-admin/internal/model"
)

// GetWormRules 获取蠕虫规则列表
func (p *MaintenanceProvider) GetWormRules(ctx context.Context) (*model.WormRuleList, error) {
	// TODO: 实现调用交换机 CLI 获取蠕虫规则列表
	return &model.WormRuleList{
		Rules: []model.WormRule{
			{ID: "1", Name: "SQL Slammer", Protocol: "UDP", Port: "1434", Stats: 15, Enabled: true},
		},
	}, nil
}

// AddWormRule 添加蠕虫规则
func (p *MaintenanceProvider) AddWormRule(ctx context.Context, req model.WormRuleRequest) error {
	// TODO: 实现调用交换机 CLI 添加蠕虫规则
	if req.Name == "" || req.Port == "" {
		return fmt.Errorf("规则名称和端口不能为空")
	}
	return nil
}

// UpdateWormRule 更新蠕虫规则
func (p *MaintenanceProvider) UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error {
	// TODO: 实现调用交换机 CLI 更新蠕虫规则
	return nil
}

// DeleteWormRule 删除蠕虫规则
func (p *MaintenanceProvider) DeleteWormRule(ctx context.Context, id string) error {
	// TODO: 实现调用交换机 CLI 删除蠕虫规则
	return nil
}

// DeleteWormRules 批量删除蠕虫规则
func (p *MaintenanceProvider) DeleteWormRules(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := p.DeleteWormRule(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ClearWormStats 清除蠕虫统计
func (p *MaintenanceProvider) ClearWormStats(ctx context.Context) error {
	// TODO: 实现调用交换机 CLI 清除统计
	return nil
}
