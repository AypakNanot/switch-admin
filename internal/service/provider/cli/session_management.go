package cli

import (
	"context"
	"time"

	"switch-admin/internal/model"
)

// GetSessions 获取会话列表
func (p *MaintenanceProvider) GetSessions(ctx context.Context) (*model.SessionListResponse, error) {
	// TODO: 实现调用交换机 CLI 获取会话列表
	return &model.SessionListResponse{
		Sessions: []model.SessionConfig{
			{SessionID: "sess_001", Username: "admin", LoginTime: "2024-03-18 08:00:00", LastActivity: time.Now().Format("2006-01-02 15:04:05"), IPAddress: "192.168.1.100"},
		},
	}, nil
}

// DeleteSession 删除会话
func (p *MaintenanceProvider) DeleteSession(ctx context.Context, sessionID string) error {
	// TODO: 实现调用交换机 CLI 删除会话
	return nil
}

// DeleteSessions 批量删除会话
func (p *MaintenanceProvider) DeleteSessions(ctx context.Context, sessionIDs []string) error {
	for _, sessionID := range sessionIDs {
		if err := p.DeleteSession(ctx, sessionID); err != nil {
			return err
		}
	}
	return nil
}
