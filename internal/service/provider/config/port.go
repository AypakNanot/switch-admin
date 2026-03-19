package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetPortList 获取端口配置列表
func (p *CLIProvider) GetPortList(ctx context.Context) (*model.PortConfigListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取端口配置列表
	// 示例命令：show interfaces status

	return &model.PortConfigListResponse{
		Ports: []model.PortConfig{},
		Total: 0,
	}, nil
}

// GetPortDetail 获取端口详情
func (p *CLIProvider) GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error) {
	// TODO: 实现调用交换机 CLI 获取单个端口配置
	// 示例命令：show interface <port-id>

	return &model.PortConfig{
		PortID:      portID,
		AdminStatus: "enable",
		LinkStatus:  "up",
		SpeedDuplex: "1000F",
		FlowControl: "off",
		Description: "",
		Aggregation: "-",
	}, nil
}

// UpdatePort 更新端口配置
func (p *CLIProvider) UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error {
	// TODO: 实现调用交换机 CLI 更新端口配置
	// 示例命令序列：
	// configure terminal
	// interface <port-id>
	// [shutdown | no shutdown]
	// speed-duplex <speed>
	// flow-control [on | off]
	// description <desc>
	return nil
}
