package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetFiles 获取文件列表
func (p *MaintenanceProvider) GetFiles(ctx context.Context, path string) (*model.FileListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.FileListResponse{
		Files: []model.FileInfo{
			{Name: "config.bin", Path: "/flash/config.bin", Size: 2048, UpdatedAt: "2024-03-15 10:00:00", Type: "config"},
			{Name: "backup_20240310.bin", Path: "/flash/backup/backup_20240310.bin", Size: 2048, UpdatedAt: "2024-03-10 15:30:00", Type: "backup"},
			{Name: "log.txt", Path: "/flash/log/log.txt", Size: 10240, UpdatedAt: time.Now().Format("2006-01-02 15:04:05"), Type: "log"},
		},
		Total: 3,
	}, nil
}

// UploadFile 上传文件
func (p *MaintenanceProvider) UploadFile(ctx context.Context, req model.FileUploadRequest) error {
	time.Sleep(500 * time.Millisecond)
	if len(req.Content) == 0 {
		return fmt.Errorf("文件内容不能为空")
	}
	return nil
}

// DeleteFile 删除文件
func (p *MaintenanceProvider) DeleteFile(ctx context.Context, path string) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

// DeleteFiles 批量删除文件
func (p *MaintenanceProvider) DeleteFiles(ctx context.Context, paths []string) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// DownloadFile 下载文件
func (p *MaintenanceProvider) DownloadFile(ctx context.Context, path string) ([]byte, string, error) {
	time.Sleep(200 * time.Millisecond)
	return []byte("mock file content"), "application/octet-stream", nil
}
