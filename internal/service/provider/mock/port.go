package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetPortList 获取端口列表
func (p *NetworkProvider) GetPortList(ctx context.Context) (*model.NetworkPortListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.NetworkPortListResponse{
		Ports: []model.Port{
			{Name: "eth0/1", Status: "up", Speed: "1000Mbps", Duplex: "full", VLAN: 10, Type: "copper"},
			{Name: "eth0/2", Status: "up", Speed: "1000Mbps", Duplex: "full", VLAN: 10, Type: "copper"},
			{Name: "eth0/3", Status: "down", Speed: "auto", Duplex: "auto", VLAN: 20, Type: "copper"},
			{Name: "eth0/4", Status: "up", Speed: "10Gbps", Duplex: "full", VLAN: 1, Type: "fiber"},
			{Name: "eth0/5", Status: "up", Speed: "1000Mbps", Duplex: "full", VLAN: 20, Type: "copper"},
			{Name: "eth0/6", Status: "up", Speed: "1000Mbps", Duplex: "full", VLAN: 20, Type: "copper"},
		},
		Total: 6,
	}, nil
}

// GetPortDetail 获取端口详情
func (p *NetworkProvider) GetPortDetail(ctx context.Context, portName string) (*model.PortDetail, error) {
	time.Sleep(50 * time.Millisecond)
	if portName == "" {
		return nil, fmt.Errorf("端口名称不能为空")
	}

	return &model.PortDetail{
		Name:        portName,
		Status:      "up",
		Speed:       "1000Mbps",
		Duplex:      "full",
		VLAN:        10,
		Type:        "copper",
		MAC:         "00:11:22:33:44:55",
		Description: "Management Port",
		RXBytes:     1234567890,
		TXBytes:     987654321,
		RXErrors:    0,
		TXErrors:    0,
	}, nil
}

// UpdatePort 更新端口配置
func (p *NetworkProvider) UpdatePort(ctx context.Context, portName string, req model.PortUpdateRequest) error {
	time.Sleep(100 * time.Millisecond)
	if portName == "" {
		return fmt.Errorf("端口名称不能为空")
	}
	return nil
}

// ResetPort 重置端口配置
func (p *NetworkProvider) ResetPort(ctx context.Context, portName string) error {
	time.Sleep(100 * time.Millisecond)
	if portName == "" {
		return fmt.Errorf("端口名称不能为空")
	}
	return nil
}

// RestartPort 重启端口
func (p *NetworkProvider) RestartPort(ctx context.Context, portName string) error {
	time.Sleep(200 * time.Millisecond)
	if portName == "" {
		return fmt.Errorf("端口名称不能为空")
	}
	return nil
}
