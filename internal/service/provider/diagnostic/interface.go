package diagnostic

import (
	"context"

	"switch-admin/internal/model"
)

// DiagnosticProvider 诊断模块接口
type DiagnosticProvider interface {
	// Ping 诊断
	ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error)

	// Traceroute 诊断
	ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error)

	// 电缆检测
	ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error)
}
