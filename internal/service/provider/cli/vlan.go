package cli

import (
	"context"
	"fmt"

	"switch-admin/internal/model"
)

// GetVLANList 获取 VLAN 列表
func (p *NetworkProvider) GetVLANList(ctx context.Context) (*model.VLANListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取 VLAN 列表
	// 示例：output, err := p.execFunc("show", "vlan")

	return &model.VLANListResponse{
		VLANs: []model.VLAN{
			{ID: 1, Name: "default", Ports: []string{"eth0/1-24"}, Status: "active"},
			{ID: 10, Name: "management", Ports: []string{"eth0/1-4"}, Status: "active"},
			{ID: 20, Name: "user", Ports: []string{"eth0/5-20"}, Status: "active"},
		},
		Total: 3,
	}, nil
}

// CreateVLAN 创建 VLAN
func (p *NetworkProvider) CreateVLAN(ctx context.Context, req model.VLANRequest) error {
	// TODO: 实现调用交换机 CLI 创建 VLAN
	if req.ID == 0 || req.Name == "" {
		return fmt.Errorf("VLAN ID 和名称不能为空")
	}
	return nil
}

// UpdateVLAN 更新 VLAN
func (p *NetworkProvider) UpdateVLAN(ctx context.Context, id int, req model.VLANRequest) error {
	// TODO: 实现调用交换机 CLI 更新 VLAN
	if req.Name == "" {
		return fmt.Errorf("VLAN 名称不能为空")
	}
	return nil
}

// DeleteVLAN 删除 VLAN
func (p *NetworkProvider) DeleteVLAN(ctx context.Context, id int) error {
	// TODO: 实现调用交换机 CLI 删除 VLAN
	if id == 1 {
		return fmt.Errorf("不能删除默认 VLAN")
	}
	return nil
}

// DeleteVLANs 批量删除 VLAN
func (p *NetworkProvider) DeleteVLANs(ctx context.Context, ids []int) error {
	for _, id := range ids {
		if err := p.DeleteVLAN(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// AddVLANPort 添加 VLAN 端口
func (p *NetworkProvider) AddVLANPort(ctx context.Context, vlanID int, port string, mode string) error {
	// TODO: 实现调用交换机 CLI 添加 VLAN 端口
	if port == "" {
		return fmt.Errorf("端口不能为空")
	}
	return nil
}

// RemoveVLANPort 移除 VLAN 端口
func (p *NetworkProvider) RemoveVLANPort(ctx context.Context, vlanID int, port string) error {
	// TODO: 实现调用交换机 CLI 移除 VLAN 端口
	if port == "" {
		return fmt.Errorf("端口不能为空")
	}
	return nil
}
