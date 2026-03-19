package diagnostic

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

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

	// 解析 ping 输出（GBK 编码）
	response := parsePingOutputGBK(output, req)

	if err != nil {
		// 检查是否是超时或无法访问的错误
		if exitErr, ok := err.(*exec.ExitError); ok {
			if len(response.Results) == 0 {
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
			response.Status = "completed"
			return response, nil
		}
		return nil, err
	}

	response.Status = "completed"
	return response, nil
}

// parsePingOutputGBK 解析 ping 命令输出（Windows GBK 编码）
func parsePingOutputGBK(output []byte, req model.PingRequest) *model.PingTaskResponse {
	response := &model.PingTaskResponse{
		TaskID:  "",
		Target:  req.Target,
		VrfID:   req.VrfID,
		Results: []model.PingResult{},
	}

	// GBK 编码字节模式
	laiZi := []byte{0xc0, 0xb4, 0xd7, 0xd4}       // 来自
	huiFu := []byte{0xb5, 0xc4, 0xbb, 0xd8, 0xb8, 0xb4} // 的回复
	shiJian := []byte{0xca, 0xb1, 0xbc, 0xe4}     // 时间
	qingQiu := []byte{0xc7, 0xeb, 0xc7, 0xf3, 0xc3, 0xac} // 请求超时
	yiFaSong := []byte{0xd2, 0xd1, 0xb7, 0xa2, 0xcb, 0xcd} // 已发送
	yiJieShou := []byte{0xd2, 0xd1, 0xbd, 0xd3, 0xca, 0xd5} // 已接收
	diuShi := []byte{0xb6, 0xaa, 0xca, 0xa7}      // 丢失
	zuiDuan := []byte{0xd7, 0xee, 0xb6, 0xcc}     // 最短
	zuiChang := []byte{0xd7, 0xee, 0xb3, 0xa4}    // 最长
	pingJun := []byte{0xc6, 0xbd, 0xbe, 0xf9}     // 平均

	lines := bytes.Split(output, []byte{'\n'})
	sent := 0
	received := 0
	var times []float64

	seq := 0
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// 检查是否响应行：包含"来自"和"的回复"
		if bytes.Contains(line, laiZi) && bytes.Contains(line, huiFu) {
			seq++
			sent++
			received++

			// 提取 TTL 值
			ttlIdx := bytes.Index(line, []byte("TTL="))
			ttl := 0
			if ttlIdx >= 0 {
				ttlStart := ttlIdx + 4
				ttlEnd := ttlStart
				for ttlEnd < len(line) && line[ttlEnd] >= '0' && line[ttlEnd] <= '9' {
					ttlEnd++
				}
				if ttlEnd > ttlStart {
					ttl, _ = strconv.Atoi(string(line[ttlStart:ttlEnd]))
				}
			}

			// 提取时间值（在"时间"之后，"ms" 之前）
			timeIdx := bytes.Index(line, shiJian)
			timeMs := 0
			if timeIdx >= 0 {
				timeStart := timeIdx + len(shiJian)
				// 跳过 < = 等符号
				for timeStart < len(line) && (line[timeStart] == '<' || line[timeStart] == '>' || line[timeStart] == '=') {
					timeStart++
				}
				timeEnd := timeStart
				for timeEnd < len(line) && line[timeEnd] >= '0' && line[timeEnd] <= '9' {
					timeEnd++
				}
				if timeEnd > timeStart {
					timeMs, _ = strconv.Atoi(string(line[timeStart:timeEnd]))
				}
			}

			times = append(times, float64(timeMs))
			response.Results = append(response.Results, model.PingResult{
				Seq:    seq,
				Time:   fmt.Sprintf("%dms", timeMs),
				TTL:    ttl,
				Status: "success",
			})
			continue
		}

		// 检查是否超时：包含"请求超时"
		if bytes.Contains(line, qingQiu) {
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

		// 解析统计行：包含"已发送"、"已接收"、"丢失"
		if bytes.Contains(line, yiFaSong) && bytes.Contains(line, yiJieShou) && bytes.Contains(line, diuShi) {
			// 提取数字
			nums := extractNumbers(line)
			if len(nums) >= 3 {
				sent = nums[0]
				received = nums[1]
				lost := nums[2]
				lossRate := float64(lost) / float64(sent) * 100
				response.Statistics = model.PingStatistics{
					Sent:     sent,
					Received: received,
					LossRate: fmt.Sprintf("%.0f%%", lossRate),
				}
			}
			continue
		}

		// 解析时间统计：包含"最短"、"最长"、"平均"
		if bytes.Contains(line, zuiDuan) && bytes.Contains(line, zuiChang) && bytes.Contains(line, pingJun) {
			nums := extractNumbers(line)
			if len(nums) >= 3 {
				response.Statistics.MinTime = fmt.Sprintf("%dms", nums[0])
				response.Statistics.MaxTime = fmt.Sprintf("%dms", nums[1])
				response.Statistics.AvgTime = fmt.Sprintf("%dms", nums[2])
			}
			continue
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

// extractNumbers 从 GBK 编码的字节切片中提取数字
func extractNumbers(data []byte) []int {
	var nums []int
	i := 0
	for i < len(data) {
		// 跳过非数字
		for i < len(data) && (data[i] < '0' || data[i] > '9') {
			i++
		}
		if i >= len(data) {
			break
		}
		// 提取数字
		start := i
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			i++
		}
		if i > start {
			num, _ := strconv.Atoi(string(data[start:i]))
			nums = append(nums, num)
		}
	}
	return nums
}

// safeString 安全地将 GBK 字节转换为字符串（用于调试，实际不使用）
func safeString(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	// GBK 编码直接转字符串会乱码，但不会影响功能
	return string(b)
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
	target := strings.ToLower(req.Target)
	isReachable := isSimulatedReachable(target)

	if !isReachable {
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
	if target == "localhost" || target == "127.0.0.1" || target == "::1" {
		return true
	}

	privatePrefixes := []string{
		"192.168.", "10.",
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

	if target == "8.8.8.8" || target == "8.8.4.4" ||
		target == "1.1.1.1" || target == "1.0.0.1" {
		return true
	}

	if strings.HasSuffix(target, ".1") && len(target) <= 15 {
		return true
	}

	return false
}

// ExecuteTraceroute 执行 Traceroute 诊断 (Mock)
func (p *MockProvider) ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error) {
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
