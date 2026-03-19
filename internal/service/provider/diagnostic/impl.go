package diagnostic

import (
	"context"

	"switch-admin/internal/model"
)

// ExecutePing 执行 Ping 诊断 (CLI)
func (p *CLIProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// TODO: 实现 CLI 模式的 Ping 诊断
	return &model.PingTaskResponse{}, nil
}

// ExecuteTraceroute 执行 Traceroute 诊断 (CLI)
func (p *CLIProvider) ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error) {
	// TODO: 实现 CLI 模式的 Traceroute 诊断
	return &model.TracerouteResponse{}, nil
}

// ExecuteCableTest 执行电缆检测 (CLI)
func (p *CLIProvider) ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error) {
	// TODO: 实现 CLI 模式的电缆检测
	return &model.CableTestResult{}, nil
}

// ExecutePing 执行 Ping 诊断 (Mock)
func (p *MockProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// Mock 数据返回
	return &model.PingTaskResponse{
		TaskID: "mock-task-id",
		Status: "completed",
		Target: req.Target,
		VrfID:  req.VrfID,
		Results: []model.PingResult{
			{Seq: 1, Time: "1.2ms", TTL: 64, Status: "success"},
			{Seq: 2, Time: "2.3ms", TTL: 64, Status: "success"},
			{Seq: 3, Time: "1.8ms", TTL: 64, Status: "success"},
			{Seq: 4, Time: "4.1ms", TTL: 64, Status: "success"},
		},
		Statistics: model.PingStatistics{
			Sent:     4,
			Received: 4,
			LossRate: "0%",
			MinTime:  "1.2ms",
			AvgTime:  "2.4ms",
			MaxTime:  "4.1ms",
		},
	}, nil
}

// ExecuteTraceroute 执行 Traceroute 诊断 (Mock)
func (p *MockProvider) ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error) {
	// Mock 数据返回
	return &model.TracerouteResponse{
		TaskID: "mock-task-id",
		Status: "completed",
		Target: req.Target,
		VrfID:  req.VrfID,
		Hops: []model.HopInfo{
			{Hop: 1, IP: "192.168.1.1", Times: []string{"1.2ms", "1.1ms", "1.3ms"}, Status: "ok"},
			{Hop: 2, IP: "10.0.0.1", Times: []string{"5.3ms", "5.1ms", "5.5ms"}, Status: "ok"},
			{Hop: 3, IP: req.Target, Times: []string{"10.1ms", "10.3ms", "10.2ms"}, Status: "ok"},
		},
	}, nil
}

// ExecuteCableTest 执行电缆检测 (Mock)
func (p *MockProvider) ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error) {
	// Mock 数据返回
	return &model.CableTestResult{
		TaskID:      "mock-task-id",
		PortID:      req.PortID,
		Status:      "completed",
		CableStatus: "normal",
		Pairs: model.CablePairs{
			PairA: model.CablePair{Status: "ok", FaultDistance: ""},
			PairB: model.CablePair{Status: "ok", FaultDistance: ""},
			PairC: model.CablePair{Status: "ok", FaultDistance: ""},
			PairD: model.CablePair{Status: "ok", FaultDistance: ""},
		},
	}, nil
}
