package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetSystemConfig 获取系统配置
func (p *MaintenanceProvider) GetSystemConfig(ctx context.Context) (*model.SystemConfig, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.SystemConfig{
		Network: model.NetworkConfig{
			IP:      "192.168.1.1",
			Mask:    "255.255.255.0",
			Gateway: "192.168.1.254",
			DNS:     "8.8.8.8",
		},
		Temperature: model.TemperatureConfig{
			Low:  -10,
			High: 60,
		},
		DeviceInfo: model.DeviceInfo{
			Name:     "Switch-Core-01",
			Location: "机房 A-机柜 3-位置 15",
			Contact:  "张三 -13800138000",
		},
		DateTime: model.DateTimeConfig{
			Timezone: "UTC+8",
			DateTime: time.Now().Format("2006-01-02T15:04:05"),
		},
	}, nil
}

// UpdateNetworkConfig 更新网络配置
func (p *MaintenanceProvider) UpdateNetworkConfig(ctx context.Context, req model.NetworkConfigRequest) error {
	time.Sleep(100 * time.Millisecond)
	// 模拟验证
	if req.IP == "" {
		return fmt.Errorf("IP 地址不能为空")
	}
	return nil
}

// UpdateTemperatureConfig 更新温度配置
func (p *MaintenanceProvider) UpdateTemperatureConfig(ctx context.Context, req model.TemperatureConfigRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Low >= req.High {
		return fmt.Errorf("低温阈值必须小于高温阈值")
	}
	return nil
}

// UpdateDeviceInfo 更新设备信息
func (p *MaintenanceProvider) UpdateDeviceInfo(ctx context.Context, req model.DeviceInfoRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Name == "" {
		return fmt.Errorf("设备名称不能为空")
	}
	return nil
}

// UpdateDateTime 更新时间日期
func (p *MaintenanceProvider) UpdateDateTime(ctx context.Context, req model.DateTimeRequest) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
