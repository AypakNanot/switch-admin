package model

import "time"

// SysConfig 系统配置
type SysConfig struct {
	ID          int64     `xorm:"pk autoincr" json:"id"`
	ConfigKey   string    `xorm:"VARCHAR(64) NOT NULL UNIQUE" json:"config_key"`
	ConfigValue string    `xorm:"TEXT NOT NULL" json:"config_value"`
	Description string    `xorm:"TEXT" json:"description"`
	CreatedAt   time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 表名
func (s *SysConfig) TableName() string {
	return "sys_config"
}
