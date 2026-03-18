package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetSNMPConfig 获取 SNMP 配置
func (p *MaintenanceProvider) GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.SNMPConfig{
		Version:     "v2c",
		Community:   "public",
		Contact:     "admin@example.com",
		Location:    "机房 A",
		Enabled:     true,
		TrapEnabled: true,
		TrapHosts: []model.TrapHost{
			{ID: 1, Host: "192.168.1.100", Port: 162, Version: "v2c", Enabled: true},
			{ID: 2, Host: "192.168.1.101", Port: 162, Version: "v2c", Enabled: false},
		},
		Communities: []model.SNMPCommunity{
			{Name: "public", Access: "read", Description: "只读团体"},
			{Name: "private", Access: "write", Description: "读写团体"},
		},
	}, nil
}

// UpdateSNMPConfig 更新 SNMP 配置
func (p *MaintenanceProvider) UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Community == "" {
		return fmt.Errorf("团体名不能为空")
	}
	return nil
}

// GetTrapHosts 获取 Trap 主机列表
func (p *MaintenanceProvider) GetTrapHosts(ctx context.Context) ([]model.TrapHost, error) {
	time.Sleep(50 * time.Millisecond)

	return []model.TrapHost{
		{ID: 1, Host: "192.168.1.100", Port: 162, Version: "v2c", Enabled: true},
		{ID: 2, Host: "192.168.1.101", Port: 162, Version: "v2c", Enabled: false},
	}, nil
}

// AddTrapHost 添加 Trap 主机
func (p *MaintenanceProvider) AddTrapHost(ctx context.Context, req model.TrapHostRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Host == "" {
		return fmt.Errorf("主机地址不能为空")
	}
	return nil
}

// DeleteTrapHost 删除 Trap 主机
func (p *MaintenanceProvider) DeleteTrapHost(ctx context.Context, host string) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// TestTrap 测试 Trap 发送
func (p *MaintenanceProvider) TestTrap(ctx context.Context, host string) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// GetSNMPCommunities 获取 SNMP 团体列表
func (p *MaintenanceProvider) GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error) {
	time.Sleep(50 * time.Millisecond)

	return []model.SNMPCommunity{
		{Name: "public", Access: "read", Description: "只读团体"},
		{Name: "private", Access: "write", Description: "读写团体"},
	}, nil
}

// AddCommunity 添加团体
func (p *MaintenanceProvider) AddCommunity(ctx context.Context, name, access, description string) error {
	time.Sleep(100 * time.Millisecond)
	if name == "" {
		return fmt.Errorf("团体名不能为空")
	}
	return nil
}

// DeleteCommunity 删除团体
func (p *MaintenanceProvider) DeleteCommunity(ctx context.Context, name string) error {
	time.Sleep(100 * time.Millisecond)
	if name == "public" {
		return fmt.Errorf("不能删除默认团体")
	}
	return nil
}
