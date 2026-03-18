package mock

import (
	"context"
	"fmt"
	"time"

	"switch-admin/internal/model"
)

// GetUsers 获取用户列表
func (p *MaintenanceProvider) GetUsers(ctx context.Context) (*model.UserListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.UserListResponse{
		Users: []model.UserConfig{
			{Username: "admin", Role: 0, RoleName: "超级管理员", CreatedAt: "2024-01-01 10:00:00"},
			{Username: "operator1", Role: 2, RoleName: "操作员", CreatedAt: "2024-02-15 14:30:00"},
			{Username: "readonly1", Role: 3, RoleName: "只读用户", CreatedAt: "2024-03-01 09:00:00"},
		},
	}, nil
}

// CreateUser 创建用户
func (p *MaintenanceProvider) CreateUser(ctx context.Context, req model.UserRequest) error {
	time.Sleep(100 * time.Millisecond)
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("用户名和密码不能为空")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("密码至少 6 位")
	}
	return nil
}

// DeleteUser 删除用户
func (p *MaintenanceProvider) DeleteUser(ctx context.Context, username string) error {
	time.Sleep(50 * time.Millisecond)
	if username == "admin" {
		return fmt.Errorf("不能删除超级管理员")
	}
	return nil
}

// DeleteUsers 批量删除用户
func (p *MaintenanceProvider) DeleteUsers(ctx context.Context, usernames []string) error {
	time.Sleep(100 * time.Millisecond)
	for _, username := range usernames {
		if username == "admin" {
			return fmt.Errorf("不能删除超级管理员")
		}
	}
	return nil
}
