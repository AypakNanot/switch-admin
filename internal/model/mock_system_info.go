package model

import (
	"fmt"
	"time"
)

// MockSystemInfo 系统信息模拟数据
type MockSystemInfo struct {
	ID               int64     `xorm:"pk autoincr" json:"id"`
	Model            string    `xorm:"VARCHAR(64) DEFAULT 'BroadEdge-S3652'" json:"model"`
	SerialNumber     string    `xorm:"VARCHAR(64) DEFAULT 'E605MT252088'" json:"serial_number"`
	MACAddress       string    `xorm:"VARCHAR(32) DEFAULT '00:07:30:D2:35:67'" json:"mac_address"`
	SoftwareVersion  string    `xorm:"VARCHAR(64) DEFAULT 'OPTEL v7.0.5.15'" json:"software_version"`
	HardwareVersion  string    `xorm:"VARCHAR(32) DEFAULT '3.0'" json:"hardware_version"`
	UptimeSeconds    int64     `xorm:"DEFAULT 0" json:"uptime_seconds"`
	BootTime         time.Time `xorm:"DATETIME" json:"boot_time"`
	LastCalculatedAt time.Time `xorm:"DATETIME" json:"-"` // 上次计算时间（用于计算运行时间）
}

// TableName 表名
func (m *MockSystemInfo) TableName() string {
	return "mock_system_info"
}

// CalculateUptime 计算运行时间字符串
func (m *MockSystemInfo) CalculateUptime() string {
	now := time.Now()
	duration := now.Sub(m.BootTime)

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d 天 %d 小时 %d 分钟", days, hours, minutes)
	}
	return fmt.Sprintf("%d 小时 %d 分钟", hours, minutes)
}
