package mock

import (
	"context"
	"switch-admin/internal/model"
	"time"
)

// ConfigProvider Mock 模式的 Config Provider
// 提供模拟数据用于开发和测试
type ConfigProvider struct{}

// NewConfigProvider 创建 Mock Config Provider
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{}
}

// GetPortList 获取端口列表（Mock 数据）
func (p *ConfigProvider) GetPortList(ctx context.Context) (*model.PortConfigListResponse, error) {
	time.Sleep(50 * time.Millisecond) // 模拟网络延迟

	return &model.PortConfigListResponse{
		Ports: []model.PortConfig{
			{
				PortID:      "GE1/0/1",
				AdminStatus: "enable",
				LinkStatus:  "up",
				SpeedDuplex: "1000F",
				FlowControl: "off",
				Description: "Server-A",
				Aggregation: "-",
			},
			{
				PortID:      "GE1/0/2",
				AdminStatus: "enable",
				LinkStatus:  "down",
				SpeedDuplex: "auto",
				FlowControl: "off",
				Description: "",
				Aggregation: "-",
			},
			{
				PortID:      "GE1/0/3",
				AdminStatus: "disable",
				LinkStatus:  "down",
				SpeedDuplex: "auto",
				FlowControl: "off",
				Description: "",
				Aggregation: "Ag1",
			},
			{
				PortID:      "GE1/0/4",
				AdminStatus: "enable",
				LinkStatus:  "up",
				SpeedDuplex: "100F",
				FlowControl: "on",
				Description: "AP-Floor2",
				Aggregation: "-",
			},
		},
		Total: 4,
	}, nil
}

// GetPortDetail 获取端口详情（Mock 数据）
func (p *ConfigProvider) GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error) {
	time.Sleep(30 * time.Millisecond)

	return &model.PortConfig{
		PortID:      portID,
		AdminStatus: "enable",
		LinkStatus:  "up",
		SpeedDuplex: "1000F",
		FlowControl: "off",
		Description: "Mock port configuration",
		Aggregation: "-",
	}, nil
}

// UpdatePort 更新端口配置（Mock）
func (p *ConfigProvider) UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error {
	time.Sleep(50 * time.Millisecond)
	// Mock 实现直接返回成功
	return nil
}
