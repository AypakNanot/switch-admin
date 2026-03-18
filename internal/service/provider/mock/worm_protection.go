package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetWormRules 获取蠕虫规则列表
func (p *MaintenanceProvider) GetWormRules(ctx context.Context) (*model.WormRuleList, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.WormRuleList{
		Rules: []model.WormRule{
			{ID: "1", Name: "SQL Slammer", Protocol: "UDP", Port: "1434", Stats: 15, Enabled: true},
			{ID: "2", Name: "Blaster", Protocol: "TCP", Port: "135-139", Stats: 8, Enabled: true},
			{ID: "3", Name: "Sasser", Protocol: "TCP", Port: "445", Stats: 3, Enabled: false},
		},
	}, nil
}

// AddWormRule 添加蠕虫规则
func (p *MaintenanceProvider) AddWormRule(ctx context.Context, req model.WormRuleRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Name == "" || req.Port == "" {
		return fmt.Errorf("规则名称和端口不能为空")
	}
	return nil
}

// UpdateWormRule 更新蠕虫规则
func (p *MaintenanceProvider) UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// DeleteWormRule 删除蠕虫规则
func (p *MaintenanceProvider) DeleteWormRule(ctx context.Context, id string) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// DeleteWormRules 批量删除蠕虫规则
func (p *MaintenanceProvider) DeleteWormRules(ctx context.Context, ids []string) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// ClearWormStats 清除蠕虫统计
func (p *MaintenanceProvider) ClearWormStats(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
