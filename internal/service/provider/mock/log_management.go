package mock

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"switch-admin/internal/model"
)

// GetLogs 获取日志列表
func (p *MaintenanceProvider) GetLogs(ctx context.Context) (*model.LogListResponse, error) {
	time.Sleep(50 * time.Millisecond)

	logs := make([]model.LogEntry, 0, 20)
	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	sources := []string{"system", "security", "network", "user"}

	for i := 1; i <= 20; i++ {
		logs = append(logs, model.LogEntry{
			ID:        i,
			Level:     levels[rand.Intn(len(levels))],
			Message:   fmt.Sprintf("Log message #%d - System event occurred", i),
			Timestamp: time.Now().Add(time.Duration(-i) * time.Hour).Format("2006-01-02 15:04:05"),
			Source:    sources[rand.Intn(len(sources))],
		})
	}

	return &model.LogListResponse{
		Logs:  logs,
		Total: 156,
	}, nil
}

// ClearLogs 清除日志
func (p *MaintenanceProvider) ClearLogs(ctx context.Context, levels []string) error {
	time.Sleep(200 * time.Millisecond)
	return nil
}
