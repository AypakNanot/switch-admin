package config

import (
	"switch-admin/internal/dao"
	"switch-admin/internal/model"
)

// MockProvider Mock 模式的 Config Provider
// 使用本地 Mock 数据实现配置功能，用于离线测试
type MockProvider struct {
	mockDataDAO *dao.MockDataDAO
}

// NewMockProvider 创建 Mock Config Provider
func NewMockProvider() *MockProvider {
	return &MockProvider{
		mockDataDAO: dao.NewMockDataDAO(),
	}
}

// mockPortToPortConfig 将 MockPort 转换为 PortConfig
func (p *MockProvider) mockPortToPortConfig(mockPort *model.MockPort) *model.PortConfig {
	adminStatus := "disable"
	if mockPort.AdminStatus {
		adminStatus = "enable"
	}

	linkStatus := "down"
	if mockPort.LinkStatus {
		linkStatus = "up"
	}

	return &model.PortConfig{
		PortID:      mockPort.PortName,
		AdminStatus: adminStatus,
		LinkStatus:  linkStatus,
		SpeedDuplex: mockPort.Speed + mockPort.Duplex,
		FlowControl: "off",
		Description: mockPort.Description,
		Aggregation: "-",
	}
}
