package cli

import (
	"context"
	"fmt"

	"switch-admin/internal/model"
)

// GetLAGList 获取链路聚合组列表
func (p *NetworkProvider) GetLAGList(ctx context.Context) (*model.LAGListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取链路聚合组列表

	return &model.LAGListResponse{
		LAGs: []model.LAG{
			{ID: 1, Name: "lag-1", Ports: []string{"eth0/1", "eth0/2"}, Status: "active", Mode: "LACP"},
			{ID: 2, Name: "lag-2", Ports: []string{"eth0/3", "eth0/4"}, Status: "active", Mode: "LACP"},
		},
		Total: 2,
	}, nil
}

// CreateLAG 创建链路聚合组
func (p *NetworkProvider) CreateLAG(ctx context.Context, req model.LAGRequest) error {
	// TODO: 实现调用交换机 CLI 创建链路聚合组
	if req.Name == "" || len(req.Ports) == 0 {
		return fmt.Errorf("聚合组名称和端口不能为空")
	}
	return nil
}

// UpdateLAG 更新链路聚合组
func (p *NetworkProvider) UpdateLAG(ctx context.Context, id int, req model.LAGRequest) error {
	// TODO: 实现调用交换机 CLI 更新链路聚合组
	if req.Name == "" {
		return fmt.Errorf("聚合组名称不能为空")
	}
	return nil
}

// DeleteLAG 删除链路聚合组
func (p *NetworkProvider) DeleteLAG(ctx context.Context, id int) error {
	// TODO: 实现调用交换机 CLI 删除链路聚合组
	return nil
}

// AddLAGPort 添加端口到聚合组
func (p *NetworkProvider) AddLAGPort(ctx context.Context, lagID int, port string) error {
	// TODO: 实现调用交换机 CLI 添加端口到聚合组
	if port == "" {
		return fmt.Errorf("端口不能为空")
	}
	return nil
}

// RemoveLAGPort 从聚合组移除端口
func (p *NetworkProvider) RemoveLAGPort(ctx context.Context, lagID int, port string) error {
	// TODO: 实现调用交换机 CLI 从聚合组移除端口
	if port == "" {
		return fmt.Errorf("端口不能为空")
	}
	return nil
}
