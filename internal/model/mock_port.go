package model

import "time"

// MockPort 端口模拟数据
type MockPort struct {
	ID          int64     `xorm:"pk autoincr" json:"id"`
	PortName    string    `xorm:"VARCHAR(32) NOT NULL UNIQUE" json:"port_name"`
	AdminStatus bool      `xorm:"DEFAULT 1" json:"admin_status"`         // 管理状态
	LinkStatus  bool      `xorm:"DEFAULT 0" json:"link_status"`          // 链路状态
	Speed       string    `xorm:"VARCHAR(16) DEFAULT '-'" json:"speed"`  // 速率
	Duplex      string    `xorm:"VARCHAR(16) DEFAULT '-'" json:"duplex"` // 双工
	Description string    `xorm:"VARCHAR(256) DEFAULT ''" json:"description"`
	RxBytes     uint64    `xorm:"BIGINT DEFAULT 0" json:"rx_bytes"`
	TxBytes     uint64    `xorm:"BIGINT DEFAULT 0" json:"tx_bytes"`
	RxPackets   uint64    `xorm:"BIGINT DEFAULT 0" json:"rx_packets"`
	TxPackets   uint64    `xorm:"BIGINT DEFAULT 0" json:"tx_packets"`
	RxErrors    uint64    `xorm:"BIGINT DEFAULT 0" json:"rx_errors"`
	TxErrors    uint64    `xorm:"BIGINT DEFAULT 0" json:"tx_errors"`
	UpdatedAt   time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 表名
func (m *MockPort) TableName() string {
	return "mock_port"
}
