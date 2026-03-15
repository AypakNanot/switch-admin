package provider

import (
	"time"
)

// PortStatus 端口状态
type PortStatus struct {
	Name        string `json:"name"`         // 端口名称
	AdminStatus string `json:"admin_status"` // 管理状态：enable/disable
	LinkStatus  string `json:"link_status"`  // 链路状态：up/down
	Speed       string `json:"speed"`        // 速率：10M/100M/1000M/10G
	Duplex      string `json:"duplex"`       // 双工模式：Full/Half
	Description string `json:"description"`  // 端口描述

	// 流量统计（64 位计数器）
	RxBytes    uint64 `json:"rx_bytes"`
	TxBytes    uint64 `json:"tx_bytes"`
	RxPackets  uint64 `json:"rx_packets"`
	TxPackets  uint64 `json:"tx_packets"`
	RxErrors   uint64 `json:"rx_errors"`
	TxErrors   uint64 `json:"tx_errors"`

	UpdatedAt time.Time `json:"updated_at"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	Model           string `json:"model"`            // 产品型号
	SerialNumber    string `json:"serial_number"`    // 序列号
	MACAddress      string `json:"mac_address"`      // MAC 地址
	SoftwareVersion string `json:"software_version"` // 软件版本
	HardwareVersion string `json:"hardware_version"` // 硬件版本
	Uptime          string `json:"uptime"`           // 运行时间
	UptimeSeconds   int64  `json:"uptime_seconds"`   // 运行时间（秒）
	BootTime        string `json:"boot_time"`        // 启动时间
}

// DataProvider 数据提供者接口
// 所有数据访问必须通过此接口，实现模式隔离
type DataProvider interface {
	// GetPortStatus 获取单个端口状态
	GetPortStatus(portID string) (*PortStatus, error)

	// GetAllPorts 获取所有端口状态
	GetAllPorts() ([]*PortStatus, error)

	// SetPortAdminStatus 设置端口管理状态
	SetPortAdminStatus(portID string, enabled bool) error

	// GetSystemInfo 获取系统信息
	GetSystemInfo() (*SystemInfo, error)

	// ClearPortStats 清零单个端口统计
	ClearPortStats(portID string) error

	// ClearAllPortStats 清零所有端口统计
	ClearAllPortStats() error
}
