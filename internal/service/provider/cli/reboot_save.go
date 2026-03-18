package cli

import (
	"context"
	"time"
)

// SaveConfig 保存配置
func (p *MaintenanceProvider) SaveConfig(ctx context.Context) error {
	// TODO: 实现调用交换机 CLI 保存配置
	// 示例：_, err := p.execFunc("write", "memory")
	time.Sleep(500 * time.Millisecond)
	return nil
}

// RebootSwitch 重启交换机
func (p *MaintenanceProvider) RebootSwitch(ctx context.Context, delay int) error {
	// TODO: 实现调用交换机 CLI 重启
	// 示例：_, err := p.execFunc("reboot", "-d", strconv.Itoa(delay))
	time.Sleep(200 * time.Millisecond)
	return nil
}

// FactoryReset 恢复出厂配置
func (p *MaintenanceProvider) FactoryReset(ctx context.Context) error {
	// TODO: 实现调用交换机 CLI 恢复出厂配置
	// 示例：_, err := p.execFunc("erase", "startup-config")
	time.Sleep(200 * time.Millisecond)
	return nil
}
