package mock

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"switch-admin/internal/model"
)

// PingProvider Mock 模式的 Ping Provider
// 用于离线测试模式，生成模拟 Ping 结果
type PingProvider struct{}

// NewPingProvider 创建 Mock Ping Provider
func NewPingProvider() *PingProvider {
	return &PingProvider{}
}

// ExecutePing 生成模拟 Ping 结果
func (p *PingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// 模拟网络延迟
	time.Sleep(100 * time.Millisecond)

	results := make([]model.PingResult, 0, req.Count)
	received := 0
	var minTime, maxTime float64 = 9999, 0
	var totalTime float64 = 0

	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= req.Count; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 模拟：80% 概率成功
		success := rand.Float64() < 0.8
		if success {
			// 模拟延迟 1-50ms
			latency := 1 + rand.Float64()*49
			results = append(results, model.PingResult{
				Seq:    i,
				Time:   fmt.Sprintf("%.2fms", latency),
				TTL:    64,
				Status: "success",
			})
			received++
			totalTime += latency
			if latency < minTime {
				minTime = latency
			}
			if latency > maxTime {
				maxTime = latency
			}
		} else {
			results = append(results, model.PingResult{
				Seq:    i,
				Time:   "*",
				TTL:    0,
				Status: "timeout",
			})
		}

		// 模拟间隔
		if i < req.Count {
			time.Sleep(time.Duration(req.Interval) * time.Second)
		}
	}

	lossRate := float64(req.Count-received) / float64(req.Count) * 100
	avgTime := 0.0
	if received > 0 {
		avgTime = totalTime / float64(received)
	}

	response := &model.PingTaskResponse{
		TaskID:  "", // 由 Service 层设置
		Status:  "completed",
		Target:  req.Target,
		VrfID:   req.VrfID,
		Results: results,
		Statistics: model.PingStatistics{
			Sent:     req.Count,
			Received: received,
			LossRate: fmt.Sprintf("%.0f%%", lossRate),
			MinTime:  fmt.Sprintf("%.2fms", minTime),
			AvgTime:  fmt.Sprintf("%.2fms", avgTime),
			MaxTime:  fmt.Sprintf("%.2fms", maxTime),
		},
	}

	if received == 0 {
		response.Error = "Destination Host Unreachable"
	}

	return response, nil
}
