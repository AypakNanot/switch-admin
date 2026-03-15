package adapter

import (
	"switch-admin/internal/service/provider"
)

// AdapterType 适配器类型
type AdapterType string

const (
	AdapterTypeCLI     AdapterType = "cli"
	AdapterTypeNetconf AdapterType = "netconf"
	AdapterTypeREST    AdapterType = "rest"
)

// SwitchAdapter 交换机适配器接口
// 所有交换机硬件操作必须通过此接口
type SwitchAdapter interface {
	// Name 返回适配器名称
	Name() string

	// Type 返回适配器类型
	Type() AdapterType

	// Connect 连接到交换机
	Connect() error

	// Close 关闭连接
	Close() error

	// IsConnected 检查连接状态
	IsConnected() bool

	// GetPortStatus 获取端口状态
	GetPortStatus(portID string) (*provider.PortStatus, error)

	// GetAllPorts 获取所有端口状态
	GetAllPorts() ([]*provider.PortStatus, error)

	// SetPortAdminStatus 设置端口管理状态
	SetPortAdminStatus(portID string, enabled bool) error

	// GetSystemInfo 获取系统信息
	GetSystemInfo() (*provider.SystemInfo, error)

	// ClearPortStats 清零单个端口统计
	ClearPortStats(portID string) error

	// ClearAllPortStats 清零所有端口统计
	ClearAllPortStats() error
}

// AdapterConfig 适配器配置
type AdapterConfig struct {
	FunctionName string                 `json:"function_name"` // 功能名称
	AdapterType  string                 `json:"adapter_type"`  // 适配器类型
	Priority     int                    `json:"priority"`      // 优先级
	Enabled      bool                   `json:"enabled"`       // 是否启用
	Config       map[string]interface{} `json:"config"`        // 配置参数
}

// ConnectionLog 连接日志
type ConnectionLog struct {
	AdapterType  string `json:"adapter_type"`
	FunctionName string `json:"function_name"`
	Operation    string `json:"operation"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	DurationMs   int    `json:"duration_ms"`
}
