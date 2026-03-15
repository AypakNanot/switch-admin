package model

import "time"

// AdapterConfig 适配器配置
type AdapterConfig struct {
	ID           int64     `xorm:"pk autoincr" json:"id"`
	FunctionName string    `xorm:"VARCHAR(64) NOT NULL" json:"function_name"` // 功能名称：port, vlan, system...
	AdapterType  string    `xorm:"VARCHAR(32) NOT NULL" json:"adapter_type"`  // 适配器类型：cli, netconf, rest
	Priority     int       `xorm:"DEFAULT 0" json:"priority"`                 // 优先级
	Enabled      bool      `xorm:"DEFAULT 1" json:"enabled"`                  // 是否启用
	Config       string    `xorm:"TEXT" json:"config"`                        // JSON 格式配置参数
	CreatedAt    time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 表名
func (a *AdapterConfig) TableName() string {
	return "adapter_config"
}

// AdapterConfigUniq 唯一索引
type AdapterConfigUniq struct {
	FunctionName string `xorm:"UNIQUE(function_name, adapter_type)"`
	AdapterType  string `xorm:"UNIQUE(function_name, adapter_type)"`
}
