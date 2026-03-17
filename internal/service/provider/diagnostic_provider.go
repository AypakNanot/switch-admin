package provider

import (
	"context"

	"switch-admin/internal/model"
)

// PingProvider Ping 诊断接口
// 所有 Ping 操作必须通过此接口，实现模式隔离
type PingProvider interface {
	// ExecutePing 执行 Ping 诊断
	// - ctx: 上下文，支持取消操作
	// - req: Ping 请求参数（目标 IP、VRF、Count、Timeout 等）
	// - 返回：Ping 任务响应（包含统计信息）和错误
	ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error)
}

// TracerouteProvider Traceroute 诊断接口
type TracerouteProvider interface {
	// ExecuteTraceroute 执行 Traceroute 诊断
	ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error)
}

// CableTestProvider 电缆检测接口
type CableTestProvider interface {
	// ExecuteCableTest 执行电缆检测
	ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error)
}
