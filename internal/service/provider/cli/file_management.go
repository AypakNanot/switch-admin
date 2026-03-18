package cli

import (
	"context"

	"switch-admin/internal/model"
)

// GetFiles 获取文件列表
func (p *MaintenanceProvider) GetFiles(ctx context.Context, path string) (*model.FileListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取文件列表
	// 示例：output, err := p.execFunc("dir", path)
	// 然后解析输出

	return &model.FileListResponse{
		Files: []model.FileInfo{
			{Name: "config.bin", Path: "/flash/config.bin", Size: 2048, UpdatedAt: "2024-03-15 10:00:00", Type: "config"},
		},
		Total: 1,
	}, nil
}

// UploadFile 上传文件
func (p *MaintenanceProvider) UploadFile(ctx context.Context, req model.FileUploadRequest) error {
	// TODO: 实现文件上传到交换机
	// 可以通过 SCP、SFTP 或 TFTP 实现
	return nil
}

// DeleteFile 删除文件
func (p *MaintenanceProvider) DeleteFile(ctx context.Context, path string) error {
	// TODO: 实现调用交换机 CLI 删除文件
	// 示例：_, err := p.execFunc("delete", path)
	return nil
}

// DeleteFiles 批量删除文件
func (p *MaintenanceProvider) DeleteFiles(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := p.DeleteFile(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFile 下载文件
func (p *MaintenanceProvider) DownloadFile(ctx context.Context, path string) ([]byte, string, error) {
	// TODO: 实现从交换机下载文件
	// 可以通过 SCP、SFTP 或 TFTP 实现
	return []byte{}, "application/octet-stream", nil
}
