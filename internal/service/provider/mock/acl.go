package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetACLList 获取 ACL 列表
func (p *NetworkProvider) GetACLList(ctx context.Context) (*model.ACLListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.ACLListResponse{
		ACLs: []model.ACL{
			{ID: 1, Name: "ACL-INBOUND", Type: "standard", Rules: 5, Status: "active"},
			{ID: 2, Name: "ACL-OUTBOUND", Type: "extended", Rules: 10, Status: "active"},
		},
		Total: 2,
	}, nil
}

// CreateACL 创建 ACL
func (p *NetworkProvider) CreateACL(ctx context.Context, req model.ACLRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Name == "" {
		return fmt.Errorf("ACL 名称不能为空")
	}
	return nil
}

// UpdateACL 更新 ACL
func (p *NetworkProvider) UpdateACL(ctx context.Context, id int, req model.ACLRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// DeleteACL 删除 ACL
func (p *NetworkProvider) DeleteACL(ctx context.Context, id int) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}

// GetACLRules 获取 ACL 规则列表
func (p *NetworkProvider) GetACLRules(ctx context.Context, aclID int) (*model.ACLRuleListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.ACLRuleListResponse{
		Rules: []model.ACLRule{
			{ID: 1, Action: "permit", Source: "192.168.1.0/24", Destination: "any", Port: "any"},
			{ID: 2, Action: "deny", Source: "any", Destination: "10.0.0.0/8", Port: "any"},
			{ID: 3, Action: "permit", Source: "172.16.0.0/16", Destination: "192.168.1.100", Port: "80"},
		},
		Total: 3,
	}, nil
}

// AddACLRule 添加 ACL 规则
func (p *NetworkProvider) AddACLRule(ctx context.Context, aclID int, req model.ACLRuleRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Action == "" {
		return fmt.Errorf("规则动作不能为空")
	}
	return nil
}

// UpdateACLRule 更新 ACL 规则
func (p *NetworkProvider) UpdateACLRule(ctx context.Context, aclID int, ruleID int, req model.ACLRuleRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// DeleteACLRule 删除 ACL 规则
func (p *NetworkProvider) DeleteACLRule(ctx context.Context, aclID int, ruleID int) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}
