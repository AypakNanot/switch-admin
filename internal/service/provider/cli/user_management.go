package cli

import (
	"context"
	"fmt"

	"switch-admin/internal/model"
)

// GetUsers 获取用户列表
func (p *MaintenanceProvider) GetUsers(ctx context.Context) (*model.UserListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取用户列表
	return &model.UserListResponse{
		Users: []model.UserConfig{
			{Username: "admin", Role: 0, RoleName: "超级管理员", CreatedAt: "2024-01-01 10:00:00"},
			{Username: "operator1", Role: 2, RoleName: "操作员", CreatedAt: "2024-02-15 14:30:00"},
		},
	}, nil
}

// CreateUser 创建用户
func (p *MaintenanceProvider) CreateUser(ctx context.Context, req model.UserRequest) error {
	// TODO: 实现调用交换机 CLI 创建用户
	// 示例：_, err := p.execFunc("add", "user", req.Username, "-p", req.Password, "-r", strconv.Itoa(req.Role))
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
	// TODO: 实现调用交换机 CLI 删除用户
	if username == "admin" {
		return fmt.Errorf("不能删除超级管理员")
	}
	return nil
}

// DeleteUsers 批量删除用户
func (p *MaintenanceProvider) DeleteUsers(ctx context.Context, usernames []string) error {
	for _, username := range usernames {
		if err := p.DeleteUser(ctx, username); err != nil {
			return err
		}
	}
	return nil
}
