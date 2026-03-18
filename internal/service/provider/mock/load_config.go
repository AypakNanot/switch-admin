package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetConfigFiles 获取配置文件列表
func (p *MaintenanceProvider) GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.LoadConfigListResponse{
		Files: []model.LoadConfigFile{
			{Name: "startup-config.bin", Path: "/flash/startup-config.bin", Size: 2048, CreatedAt: "2024-03-01 10:00:00"},
			{Name: "backup-20240315.bin", Path: "/flash/backup/backup-20240315.bin", Size: 2048, CreatedAt: "2024-03-15 15:30:00"},
			{Name: "backup-20240310.bin", Path: "/flash/backup/backup-20240310.bin", Size: 2048, CreatedAt: "2024-03-10 09:00:00"},
		},
	}, nil
}

// LoadConfig 加载配置
func (p *MaintenanceProvider) LoadConfig(ctx context.Context, configFile string) error {
	time.Sleep(500 * time.Millisecond)
	if configFile == "" {
		return fmt.Errorf("配置文件名不能为空")
	}
	return nil
}
