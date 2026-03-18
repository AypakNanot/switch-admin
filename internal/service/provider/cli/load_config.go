package cli

import (
	"context"
	"fmt"

	"switch-admin/internal/model"
)

// GetConfigFiles 获取配置文件列表
func (p *MaintenanceProvider) GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取配置文件列表
	return &model.LoadConfigListResponse{
		Files: []model.LoadConfigFile{
			{Name: "startup-config.bin", Path: "/flash/startup-config.bin", Size: 2048, CreatedAt: "2024-03-01 10:00:00"},
		},
	}, nil
}

// LoadConfig 加载配置
func (p *MaintenanceProvider) LoadConfig(ctx context.Context, configFile string) error {
	// TODO: 实现调用交换机 CLI 加载配置
	// 示例：_, err := p.execFunc("copy", configFile, "running-config")
	if configFile == "" {
		return fmt.Errorf("配置文件名不能为空")
	}
	return nil
}
