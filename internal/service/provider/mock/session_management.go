package mock

import (
	"context"
	"time"

	"switch-admin/internal/model"
)

// GetSessions 获取会话列表
func (p *MaintenanceProvider) GetSessions(ctx context.Context) (*model.SessionListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	return &model.SessionListResponse{
		Sessions: []model.SessionConfig{
			{SessionID: "sess_001", Username: "admin", LoginTime: "2024-03-18 08:00:00", LastActivity: time.Now().Format("2006-01-02 15:04:05"), IPAddress: "192.168.1.100"},
			{SessionID: "sess_002", Username: "operator1", LoginTime: "2024-03-18 09:30:00", LastActivity: time.Now().Format("2006-01-02 15:04:05"), IPAddress: "192.168.1.101"},
		},
	}, nil
}

// DeleteSession 删除会话
func (p *MaintenanceProvider) DeleteSession(ctx context.Context, sessionID string) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}

// DeleteSessions 批量删除会话
func (p *MaintenanceProvider) DeleteSessions(ctx context.Context, sessionIDs []string) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
