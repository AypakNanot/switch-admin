package diagnostic

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"switch-admin/internal/model"
)

// ExecutePing 执行 Ping 诊断 (CLI)
func (p *CLIProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// 构建 ping 命令参数
	args := []string{"-n", strconv.Itoa(req.Count)}
	if req.Timeout > 0 {
		args = append(args, "-w", strconv.Itoa(req.Timeout*1000)) // Windows 超时单位是毫秒
	}
	args = append(args, req.Target)

	// 执行 ping 命令
	cmd := exec.Command("ping", args...)
	output, err := cmd.CombinedOutput()

	// 解析 ping 输出
	response := parsePingOutput(string(output), req)

	if err != nil {
		// 检查是否是超时或无法访问的错误
		if exitErr, ok := err.(*exec.ExitError); ok {
			// ping 命令返回非 0 退出码，但可能仍有部分输出
			if len(response.Results) == 0 {
				// 没有任何成功响应，返回失败
				return &model.PingTaskResponse{
					TaskID:  "",
					Status:  "failed",
					Target:  req.Target,
					VrfID:   req.VrfID,
					Error:   exitErr.Error(),
					Results: []model.PingResult{},
					Statistics: model.PingStatistics{
						Sent:     req.Count,
						Received: 0,
						LossRate: "100%",
					},
				}, nil
			}
			// 有部分成功，返回部分结果
			response.Status = "completed"
			return response, nil
		}
		// 其他错误（如命令不存在）
		return nil, err
	}

	response.Status = "completed"
	return response, nil
}

// parsePingOutput 解析 ping 命令输出（Windows 格式）
func parsePingOutput(output string, req model.PingRequest) *model.PingTaskResponse {
	response := &model.PingTaskResponse{
		TaskID:  "",
		Target:  req.Target,
		VrfID:   req.VrfID,
		Results: []model.PingResult{},
	}

	lines := strings.Split(output, "\n")
	sent := 0
	received := 0
	var times []float64

	// 正则表达式匹配响应行
	// Windows 格式：来自 192.168.1.1 的回复：字节=32 时间=1ms TTL=64
	replyRegex := regexp.MustCompile(`来自\s+([\d\.]+)\s+的回复.*?字节=(\d+).*?时间[<>=]?(\d+)ms.*?TTL=(\d+)`)
	// 或：Reply from 192.168.1.1: bytes=32 time=1ms TTL=64
	replyRegexEn := regexp.MustCompile(`Reply from\s+([\d\.]+):.*?bytes=(\d+).*?time[<>=]?(\d+)ms.*?TTL=(\d+)`)
	// 超时：请求超时。或 Request timed out.
	timeoutRegex := regexp.MustCompile(`请求超时。|Request timed out.`)
	// 统计行：Ping 统计数据...
	statsRegex := regexp.MustCompile(`已发送\s*=\s*(\d+)，\s*已接收\s*=\s*(\d+)，\s*丢失\s*=\s*(\d+)`)
	statsRegexEn := regexp.MustCompile(`Packets:\s*Sent\s*=\s*(\d+),\s*Received\s*=\s*(\d+),\s*Lost\s*=\s*(\d+)`)
	timeRegex := regexp.MustCompile(`最短\s*=\s*(\d+)ms,\s*最长\s*=\s*(\d+)ms,\s*平均\s*=\s*(\d+)ms`)
	timeRegexEn := regexp.MustCompile(`Minimum\s*=\s*(\d+)ms,\s*Maximum\s*=\s*(\d+)ms,\s*Average\s*=\s*(\d+)ms`)

	seq := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 检查是否响应
		matches := replyRegex.FindStringSubmatch(line)
		if matches == nil {
			matches = replyRegexEn.FindStringSubmatch(line)
		}
		if matches != nil {
			seq++
			sent++
			received++
			timeMs, _ := strconv.Atoi(matches[3])
			ttl, _ := strconv.Atoi(matches[4])
			times = append(times, float64(timeMs))
			response.Results = append(response.Results, model.PingResult{
				Seq:    seq,
				Time:   fmt.Sprintf("%dms", timeMs),
				TTL:    ttl,
				Status: "success",
			})
			continue
		}

		// 检查是否超时
		if timeoutRegex.MatchString(line) {
			seq++
			sent++
			response.Results = append(response.Results, model.PingResult{
				Seq:    seq,
				Time:   "*",
				TTL:    0,
				Status: "timeout",
			})
			continue
		}

		// 解析统计
		matches = statsRegex.FindStringSubmatch(line)
		if matches == nil {
			matches = statsRegexEn.FindStringSubmatch(line)
		}
		if matches != nil {
			sent, _ = strconv.Atoi(matches[1])
			received, _ = strconv.Atoi(matches[2])
			lost, _ := strconv.Atoi(matches[3])
			lossRate := float64(lost) / float64(sent) * 100
			response.Statistics = model.PingStatistics{
				Sent:     sent,
				Received: received,
				LossRate: fmt.Sprintf("%.0f%%", lossRate),
			}
			continue
		}

		// 解析时间统计
		matches = timeRegex.FindStringSubmatch(line)
		if matches == nil {
			matches = timeRegexEn.FindStringSubmatch(line)
		}
		if matches != nil && len(times) > 0 {
			minTime, _ := strconv.Atoi(matches[1])
			maxTime, _ := strconv.Atoi(matches[2])
			avgTime, _ := strconv.Atoi(matches[3])
			response.Statistics.MinTime = fmt.Sprintf("%dms", minTime)
			response.Statistics.MaxTime = fmt.Sprintf("%dms", maxTime)
			response.Statistics.AvgTime = fmt.Sprintf("%dms", avgTime)
		}
	}

	// 如果没有从统计行获取，手动计算
	if response.Statistics.Sent == 0 && len(response.Results) > 0 {
		sent = len(response.Results)
		received = 0
		for _, r := range response.Results {
			if r.Status == "success" {
				received++
			}
		}
		lossRate := float64(sent-received) / float64(sent) * 100
		response.Statistics.Sent = sent
		response.Statistics.Received = received
		response.Statistics.LossRate = fmt.Sprintf("%.0f%%", lossRate)

		if len(times) > 0 {
			minTime := times[0]
			maxTime := times[0]
			sum := 0.0
			for _, t := range times {
				if t < minTime {
					minTime = t
				}
				if t > maxTime {
					maxTime = t
				}
				sum += t
			}
			response.Statistics.MinTime = fmt.Sprintf("%.1fms", minTime)
			response.Statistics.MaxTime = fmt.Sprintf("%.1fms", maxTime)
			response.Statistics.AvgTime = fmt.Sprintf("%.1fms", sum/float64(len(times)))
		}
	}

	return response
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
// 模拟真实的 Ping 行为：可达地址返回成功，不可达地址返回失败
func (p *MockProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// 根据目标地址模拟不同的 Ping 结果
	// 以下地址模拟为"可达"：
	//   - 192.168.x.x, 10.x.x.x, 172.16-31.x.x (私有地址)
	//   - 8.8.8.8, 1.1.1.1 (公共 DNS)
	//   - localhost, 127.0.0.1
	// 其他地址模拟为"不可达"

	target := strings.ToLower(req.Target)
	isReachable := isSimulatedReachable(target)

	if !isReachable {
		// 模拟不可达的 Ping 结果
		return &model.PingTaskResponse{
			TaskID: "mock-task-id",
			Status: "completed",
			Target: req.Target,
			VrfID:  req.VrfID,
			Results: []model.PingResult{
				{Seq: 1, Time: "*", TTL: 0, Status: "timeout"},
				{Seq: 2, Time: "*", TTL: 0, Status: "timeout"},
				{Seq: 3, Time: "*", TTL: 0, Status: "timeout"},
				{Seq: 4, Time: "*", TTL: 0, Status: "timeout"},
			},
			Statistics: model.PingStatistics{
				Sent:     req.Count,
				Received: 0,
				LossRate: "100%",
				MinTime:  "0ms",
				AvgTime:  "0ms",
				MaxTime:  "0ms",
			},
			Error: "请求超时 - 目标主机不可达",
		}, nil
	}

	// 模拟成功的 Ping 结果
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

// isSimulatedReachable 判断目标地址是否应该模拟为"可达"
func isSimulatedReachable(target string) bool {
	// 本地地址
	if target == "localhost" || target == "127.0.0.1" || target == "::1" {
		return true
	}

	// 私有地址段
	privatePrefixes := []string{
		"192.168.",
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
	}
	for _, prefix := range privatePrefixes {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}

	// 公共 DNS 服务器
	if target == "8.8.8.8" || target == "8.8.4.4" ||
		target == "1.1.1.1" || target == "1.0.0.1" {
		return true
	}

	// 默认网关常见地址
	if strings.HasSuffix(target, ".1") && len(target) <= 15 {
		return true
	}

	return false
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
