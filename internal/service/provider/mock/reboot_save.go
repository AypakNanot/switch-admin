package mock

import (
	"context"
	"time"
)

// SaveConfig 保存配置
func (p *MaintenanceProvider) SaveConfig(ctx context.Context) error {
	time.Sleep(500 * time.Millisecond)
	return nil
}

// RebootSwitch 重启交换机
func (p *MaintenanceProvider) RebootSwitch(ctx context.Context, delay int) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// FactoryReset 恢复出厂配置
func (p *MaintenanceProvider) FactoryReset(ctx context.Context) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}
