package cli

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"switch-admin/internal/model"
)

// PingProvider CLI 模式的 Ping Provider
// 执行系统 Ping 命令并解析结果
type PingProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewPingProvider 创建 CLI Ping Provider
func NewPingProvider() *PingProvider {
	return &PingProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}

// ExecutePing 执行系统 Ping 命令并解析结果
func (p *PingProvider) ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error) {
	// 1. 构建 Ping 命令
	cmd := "ping"
	args := p.buildPingArgs(req)

	// 2. 处理 VRF 路由表（如果指定）
	if req.VrfID != "" {
		cmd, args = p.wrapVRFCommand(req.VrfID, cmd, args)
	}

	// 3. 执行命令
	startTime := time.Now()
	output, err := p.execFunc(cmd, args...)
	execTime := time.Since(startTime)

	// 4. 解析输出
	return p.parsePingOutput(output, err, req, execTime)
}

// buildPingArgs 构建 Ping 命令参数
func (p *PingProvider) buildPingArgs(req model.PingRequest) []string {
	args := []string{}

	// Count
	if runtime.GOOS == "windows" {
		args = append(args, "-n", strconv.Itoa(req.Count))
	} else {
		args = append(args, "-c", strconv.Itoa(req.Count))
	}

	// Timeout (Windows 单位毫秒，Linux 单位秒)
	if runtime.GOOS == "windows" {
		args = append(args, "-w", strconv.Itoa(req.Timeout*1000))
	} else {
		args = append(args, "-W", strconv.Itoa(req.Timeout))
	}

	// Interval (仅 Linux 支持)
	if runtime.GOOS != "windows" && req.Interval > 0 {
		args = append(args, "-i", strconv.Itoa(req.Interval))
	}

	// Target
	args = append(args, req.Target)

	return args
}

// wrapVRFCommand 包装 VRF 命令
func (p *PingProvider) wrapVRFCommand(vrfID, cmd string, args []string) (string, []string) {
	if runtime.GOOS == "windows" {
		// Windows: 使用 -S 指定源地址（模拟 VRF）
		// 实际 VRF 支持需要 Windows Server
		return cmd, args
	} else {
		// Linux: ip vrf exec <vrf> ping ...
		newArgs := append([]string{"vrf", "exec", vrfID, cmd}, args...)
		return "ip", newArgs
	}
}

// parsePingOutput 解析 Ping 命令输出
func (p *PingProvider) parsePingOutput(output []byte, err error, req model.PingRequest, execTime time.Duration) (*model.PingTaskResponse, error) {
	response := &model.PingTaskResponse{
		TaskID:     "",
		Status:     "completed",
		Target:     req.Target,
		VrfID:      req.VrfID,
		Results:    make([]model.PingResult, 0),
		Statistics: model.PingStatistics{},
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	var rttTimes []float64
	received := 0

	// 根据操作系统选择解析器
	if runtime.GOOS == "windows" {
		response.Results, rttTimes, received = p.parseWindowsOutput(lines)
	} else {
		response.Results, rttTimes, received = p.parseLinuxOutput(lines)
	}

	// 计算统计信息
	sent := req.Count
	lossRate := float64(sent-received) / float64(sent) * 100

	var minTime, maxTime, avgTime float64
	if len(rttTimes) > 0 {
		minTime = rttTimes[0]
		maxTime = rttTimes[0]
		var total float64
		for _, rtt := range rttTimes {
			total += rtt
			if rtt < minTime {
				minTime = rtt
			}
			if rtt > maxTime {
				maxTime = rtt
			}
		}
		avgTime = total / float64(len(rttTimes))
	} else if req.Count > 0 && minTime == 9999 {
		minTime = 0
	}

	response.Statistics = model.PingStatistics{
		Sent:     sent,
		Received: received,
		LossRate: fmt.Sprintf("%.0f%%", lossRate),
		MinTime:  formatRTT(minTime),
		AvgTime:  formatRTT(avgTime),
		MaxTime:  formatRTT(maxTime),
	}

	// 检查是否有错误
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 命令执行失败（如网络不可达）
			response.Error = fmt.Sprintf("网络不可达：%v", exitErr)
			response.Status = "failed"
		}
	}

	// 检查是否所有请求都超时
	if received == 0 {
		response.Error = "Destination Host Unreachable"
	}

	return response, nil
}

// parseWindowsOutput 解析 Windows Ping 输出
func (p *PingProvider) parseWindowsOutput(lines []string) ([]model.PingResult, []float64, int) {
	results := make([]model.PingResult, 0)
	var rttTimes []float64
	received := 0
	seq := 0

	// Windows 输出正则
	// 来自 192.168.1.1 的回复：字节=32 时间=1ms TTL=64
	replyRegex := regexp.MustCompile(`来自 (.+?) 的回复：字节= (\d+) 时间= ([0-9.]+)ms TTL= (\d+)`)
	// 请求超时。
	timeoutRegex := regexp.MustCompile(`请求超时。`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if matches := replyRegex.FindStringSubmatch(line); matches != nil {
			seq++
			rtt, _ := strconv.ParseFloat(matches[3], 64)
			results = append(results, model.PingResult{
				Seq:    seq,
				Time:   fmt.Sprintf("%.2fms", rtt),
				TTL:    64,
				Status: "success",
			})
			rttTimes = append(rttTimes, rtt)
			received++
		} else if timeoutRegex.MatchString(line) {
			seq++
			results = append(results, model.PingResult{
				Seq:    seq,
				Time:   "*",
				TTL:    0,
				Status: "timeout",
			})
		}
	}

	return results, rttTimes, received
}

// parseLinuxOutput 解析 Linux Ping 输出
func (p *PingProvider) parseLinuxOutput(lines []string) ([]model.PingResult, []float64, int) {
	results := make([]model.PingResult, 0)
	var rttTimes []float64
	received := 0

	// Linux 输出正则
	// 64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=1.23 ms
	replyRegex := regexp.MustCompile(`(\d+) bytes from (.+?): icmp_seq= (\d+) ttl= (\d+) time= ([0-9.]+) ms`)
	// Request timeout for icmp_seq 1
	timeoutRegex := regexp.MustCompile(`Request timeout for icmp_seq (\d+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if matches := replyRegex.FindStringSubmatch(line); matches != nil {
			seq, _ := strconv.Atoi(matches[3])
			rtt, _ := strconv.ParseFloat(matches[5], 64)
			results = append(results, model.PingResult{
				Seq:    seq,
				Time:   fmt.Sprintf("%.2fms", rtt),
				TTL:    64,
				Status: "success",
			})
			rttTimes = append(rttTimes, rtt)
			received++
		} else if matches := timeoutRegex.FindStringSubmatch(line); matches != nil {
			seq, _ := strconv.Atoi(matches[1])
			results = append(results, model.PingResult{
				Seq:    seq,
				Time:   "*",
				TTL:    0,
				Status: "timeout",
			})
		}
	}

	return results, rttTimes, received
}

// formatRTT 格式化 RTT 时间
func formatRTT(rtt float64) string {
	if rtt == 0 {
		return "0ms"
	}
	return fmt.Sprintf("%.2fms", rtt)
}
