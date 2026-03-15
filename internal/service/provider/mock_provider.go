package provider

import (
	"fmt"
	"time"

	"switch-admin/internal/dao"
	"switch-admin/internal/model"
)

// MockProvider Mock 数据提供者
// 用于离线测试模式，数据存储在 SQLite 数据库中
type MockProvider struct {
	mockDataDAO *dao.MockDataDAO
}

// NewMockProvider 创建 Mock 数据提供者
func NewMockProvider() *MockProvider {
	return &MockProvider{
		mockDataDAO: dao.NewMockDataDAO(),
	}
}

// GetPortStatus 获取端口状态
func (p *MockProvider) GetPortStatus(portID string) (*PortStatus, error) {
	mockPort, err := p.mockDataDAO.GetMockPortByName(portID)
	if err != nil {
		return nil, err
	}

	return p.mockPortToPortStatus(mockPort), nil
}

// GetAllPorts 获取所有端口状态
func (p *MockProvider) GetAllPorts() ([]*PortStatus, error) {
	mockPorts, err := p.mockDataDAO.GetAllMockPorts()
	if err != nil {
		return nil, err
	}

	var ports []*PortStatus
	for _, mockPort := range mockPorts {
		ports = append(ports, p.mockPortToPortStatus(&mockPort))
	}

	return ports, nil
}

// SetPortAdminStatus 设置端口管理状态
func (p *MockProvider) SetPortAdminStatus(portID string, enabled bool) error {
	return p.mockDataDAO.UpdateMockPortAdminStatus(portID, enabled)
}

// GetSystemInfo 获取系统信息
func (p *MockProvider) GetSystemInfo() (*SystemInfo, error) {
	mockInfo, err := p.mockDataDAO.GetMockSystemInfo()
	if err != nil {
		return nil, err
	}

	return p.mockSystemInfoToSystemInfo(mockInfo), nil
}

// ClearPortStats 清零单个端口统计
func (p *MockProvider) ClearPortStats(portID string) error {
	return p.mockDataDAO.ClearMockPortStats(portID)
}

// ClearAllPortStats 清零所有端口统计
func (p *MockProvider) ClearAllPortStats() error {
	return p.mockDataDAO.ClearAllMockPortsStats()
}

// mockPortToPortStatus 将 MockPort 转换为 PortStatus
func (p *MockProvider) mockPortToPortStatus(mockPort *model.MockPort) *PortStatus {
	adminStatus := "disable"
	if mockPort.AdminStatus {
		adminStatus = "enable"
	}

	linkStatus := "down"
	if mockPort.LinkStatus {
		linkStatus = "up"
	}

	return &PortStatus{
		Name:        mockPort.PortName,
		AdminStatus: adminStatus,
		LinkStatus:  linkStatus,
		Speed:       mockPort.Speed,
		Duplex:      mockPort.Duplex,
		Description: mockPort.Description,
		RxBytes:     mockPort.RxBytes,
		TxBytes:     mockPort.TxBytes,
		RxPackets:   mockPort.RxPackets,
		TxPackets:   mockPort.TxPackets,
		RxErrors:    mockPort.RxErrors,
		TxErrors:    mockPort.TxErrors,
		UpdatedAt:   mockPort.UpdatedAt,
	}
}

// mockSystemInfoToSystemInfo 将 MockSystemInfo 转换为 SystemInfo
func (p *MockProvider) mockSystemInfoToSystemInfo(mockInfo *model.MockSystemInfo) *SystemInfo {
	// 计算运行时间字符串
	now := time.Now()
	duration := now.Sub(mockInfo.BootTime)

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	uptimeStr := fmt.Sprintf("%d 小时 %d 分钟", hours, minutes)
	if days > 0 {
		uptimeStr = fmt.Sprintf("%d 天 %d 小时 %d 分钟", days, hours, minutes)
	}

	return &SystemInfo{
		Model:           mockInfo.Model,
		SerialNumber:    mockInfo.SerialNumber,
		MACAddress:      mockInfo.MACAddress,
		SoftwareVersion: mockInfo.SoftwareVersion,
		HardwareVersion: mockInfo.HardwareVersion,
		Uptime:          uptimeStr,
		UptimeSeconds:   int64(duration.Seconds()),
		BootTime:        mockInfo.BootTime.Format(time.RFC3339),
	}
}
