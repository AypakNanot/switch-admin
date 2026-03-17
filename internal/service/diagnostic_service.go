package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"switch-admin/internal/model"
	"switch-admin/internal/service/mode"
)

// DiagnosticService 网络诊断服务
type DiagnosticService struct {
	pingTasks       map[string]*model.PingTask
	pingResults     map[string]*model.PingTaskResponse
	traceTasks      map[string]*model.TracerouteTask
	traceResults    map[string]*model.TracerouteResponse
	cableTasks      map[string]string
	cableResults    map[string]*model.CableTestResult
	mu              sync.RWMutex
	modeResolver    *mode.ModeResolver
}

var diagnosticService *DiagnosticService
var diagnosticOnce sync.Once

// GetDiagnosticService 获取诊断服务单例
func GetDiagnosticService() *DiagnosticService {
	diagnosticOnce.Do(func() {
		diagnosticService = &DiagnosticService{
			pingTasks:    make(map[string]*model.PingTask),
			pingResults:  make(map[string]*model.PingTaskResponse),
			traceTasks:   make(map[string]*model.TracerouteTask),
			traceResults: make(map[string]*model.TracerouteResponse),
			cableTasks:   make(map[string]string),
			cableResults: make(map[string]*model.CableTestResult),
			modeResolver: mode.NewModeResolver(mode.ModeResolverConfig{
				InitialMode: mode.ModeMock,
			}),
		}
		go diagnosticService.cleanupRoutine()
	})
	return diagnosticService
}

// cleanupRoutine 定期清理超时任务
func (s *DiagnosticService) cleanupRoutine() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		// 清理完成的 Ping 任务（超过 5 分钟）
		for taskID, task := range s.pingTasks {
			if task.Status == "completed" || task.Status == "failed" {
				// 简单实现：实际应记录创建时间
				delete(s.pingTasks, taskID)
				delete(s.pingResults, taskID)
			}
		}
		// 清理完成的 Traceroute 任务
		for taskID, task := range s.traceTasks {
			if task.Status == "completed" || task.Status == "failed" {
				delete(s.traceTasks, taskID)
				delete(s.traceResults, taskID)
			}
		}
		s.mu.Unlock()
	}
}

// CreatePingTask 创建 Ping 任务
func (s *DiagnosticService) CreatePingTask(req model.PingRequest) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskID := fmt.Sprintf("ping_%d", time.Now().UnixNano())
	task := &model.PingTask{
		TaskID:   taskID,
		Status:   "running",
		Target:   req.Target,
		VrfID:    req.VrfID,
		Count:    req.Count,
		Timeout:  req.Timeout,
		Interval: req.Interval,
	}
	s.pingTasks[taskID] = task

	// 异步执行 Ping
	go s.executePing(taskID, req)

	return taskID, nil
}

// executePing 执行 Ping（使用 Provider 模式）
func (s *DiagnosticService) executePing(taskID string, req model.PingRequest) {
	// 通过 ModeResolver 获取 Provider
	pingProvider := s.modeResolver.GetPingProvider()

	// 执行 Ping
	response, err := pingProvider.ExecutePing(context.Background(), req)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		// 执行失败
		s.pingResults[taskID] = &model.PingTaskResponse{
			TaskID:  taskID,
			Status:  "failed",
			Target:  req.Target,
			VrfID:   req.VrfID,
			Error:   err.Error(),
			Results: []model.PingResult{},
			Statistics: model.PingStatistics{
				Sent:     req.Count,
				Received: 0,
				LossRate: "100%",
			},
		}
		if task, ok := s.pingTasks[taskID]; ok {
			task.Status = "failed"
		}
		return
	}

	// 设置 TaskID 并存储结果
	response.TaskID = taskID
	s.pingResults[taskID] = response

	if task, ok := s.pingTasks[taskID]; ok {
		task.Status = "completed"
	}
}

// GetPingTaskResult 获取 Ping 任务结果
func (s *DiagnosticService) GetPingTaskResult(taskID string) (*model.PingTaskResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.pingResults[taskID]
	return result, ok
}

// DeletePingTask 删除 Ping 任务
func (s *DiagnosticService) DeletePingTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pingTasks, taskID)
	delete(s.pingResults, taskID)
}

// CreateTracerouteTask 创建 Traceroute 任务
func (s *DiagnosticService) CreateTracerouteTask(req model.TracerouteRequest) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskID := fmt.Sprintf("trace_%d", time.Now().UnixNano())
	task := &model.TracerouteTask{
		TaskID:  taskID,
		Status:  "running",
		Target:  req.Target,
		VrfID:   req.VrfID,
		MaxHops: req.MaxHops,
		Timeout: req.Timeout,
		Probes:  req.Probes,
	}
	s.traceTasks[taskID] = task

	// 异步执行 Traceroute
	go s.executeTraceroute(taskID, req)

	return taskID, nil
}

// executeTraceroute 执行 Traceroute（Mock 实现）
func (s *DiagnosticService) executeTraceroute(taskID string, req model.TracerouteRequest) {
	time.Sleep(200 * time.Millisecond)

	hops := make([]model.HopInfo, 0, req.MaxHops)
	rand.Seed(time.Now().UnixNano())

	// 模拟 4-10 跳
	totalHops := 4 + rand.Intn(7)
	reached := false

	for i := 1; i <= req.MaxHops && !reached; i++ {
		if i > totalHops {
			reached = true
			break
		}

		times := make([]string, req.Probes)
		status := "ok"
		var ip string

		// 模拟：某些跳超时
		if rand.Float64() < 0.1 {
			for j := range times {
				times[j] = "*"
			}
			status = "timeout"
			ip = "*"
		} else {
			// 生成模拟 IP
			ip = fmt.Sprintf("10.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256))
			for j := range times {
				latency := 1 + rand.Float64()*50
				times[j] = fmt.Sprintf("%.2fms", latency)
			}
		}

		if i == totalHops {
			status = "destination"
			ip = req.Target
			reached = true
		}

		hops = append(hops, model.HopInfo{
			Hop:    i,
			IP:     ip,
			Times:  times,
			Status: status,
		})

		time.Sleep(100 * time.Millisecond)
	}

	response := &model.TracerouteResponse{
		TaskID:    taskID,
		Status:    "completed",
		Target:    req.Target,
		VrfID:     req.VrfID,
		TotalHops: len(hops),
		Hops:      hops,
	}

	s.mu.Lock()
	s.traceResults[taskID] = response
	if task, ok := s.traceTasks[taskID]; ok {
		task.Status = "completed"
	}
	s.mu.Unlock()
}

// GetTracerouteTaskResult 获取 Traceroute 任务结果
func (s *DiagnosticService) GetTracerouteTaskResult(taskID string) (*model.TracerouteResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.traceResults[taskID]
	return result, ok
}

// DeleteTracerouteTask 删除 Traceroute 任务
func (s *DiagnosticService) DeleteTracerouteTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.traceTasks, taskID)
	delete(s.traceResults, taskID)
}

// GetDetectablePorts 获取可检测端口列表
func (s *DiagnosticService) GetDetectablePorts() []model.PortInfo {
	// Mock 端口列表
	return []model.PortInfo{
		{
			PortID:      "eth0/1",
			Name:        "eth0/1",
			Type:        "electrical",
			AdminStatus: "up",
			LinkStatus:  "up",
			Label:       "eth0/1 (电口，Admin UP, Link UP)",
			Detectable:  true,
		},
		{
			PortID:      "eth0/2",
			Name:        "eth0/2",
			Type:        "electrical",
			AdminStatus: "up",
			LinkStatus:  "down",
			Label:       "eth0/2 (电口，Admin UP, Link DOWN) - 可检测断线",
			Detectable:  true,
			Hint:        "Link DOWN，可检测断线故障",
		},
		{
			PortID:      "eth0/3",
			Name:        "eth0/3",
			Type:        "electrical",
			AdminStatus: "down",
			LinkStatus:  "down",
			Label:       "eth0/3 (电口，Admin DOWN)",
			Detectable:  false,
			Hint:        "端口已关闭，需先启用",
		},
		{
			PortID:      "eth0/4",
			Name:        "eth0/4",
			Type:        "optical",
			AdminStatus: "up",
			LinkStatus:  "up",
			Label:       "eth0/4 (光口)",
			Detectable:  false,
			Hint:        "光口不支持虚拟电缆检测",
		},
	}
}

// ExecuteCableTest 执行电缆检测
func (s *DiagnosticService) ExecuteCableTest(req model.CableTestRequest) (*model.CableTestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查端口
	ports := s.GetDetectablePorts()
	var targetPort *model.PortInfo
	for _, port := range ports {
		if port.PortID == req.PortID {
			targetPort = &port
			break
		}
	}

	if targetPort == nil {
		return nil, fmt.Errorf("端口不存在")
	}

	if !targetPort.Detectable {
		if targetPort.Type == "optical" {
			return nil, fmt.Errorf("光口不支持虚拟电缆检测")
		}
		if targetPort.AdminStatus == "down" {
			return nil, fmt.Errorf("端口已关闭，请先在端口配置中启用该端口")
		}
	}

	taskID := fmt.Sprintf("cable_%d", time.Now().UnixNano())
	s.cableTasks[req.PortID] = taskID

	// 模拟检测结果
	result := &model.CableTestResult{
		TaskID:      taskID,
		PortID:      req.PortID,
		Status:      "completed",
		AdminStatus: targetPort.AdminStatus,
		LinkStatus:  targetPort.LinkStatus,
	}

	rand.Seed(time.Now().UnixNano())

	// 根据链路状态模拟结果
	if targetPort.LinkStatus == "down" {
		// 模拟断线
		faultDistance := 20 + rand.Intn(50)
		result.CableStatus = "open"
		result.FaultDesc = "断路"
		result.FaultDistance = fmt.Sprintf("约 %dm 处断开", faultDistance)
		result.Pairs = model.CablePairs{
			PairA: model.CablePair{Status: "open", FaultDistance: fmt.Sprintf("%dm", faultDistance)},
			PairB: model.CablePair{Status: "open", FaultDistance: fmt.Sprintf("%dm", faultDistance+2)},
			PairC: model.CablePair{Status: "open", FaultDistance: fmt.Sprintf("%dm", faultDistance-2)},
			PairD: model.CablePair{Status: "open", FaultDistance: fmt.Sprintf("%dm", faultDistance+1)},
		}
	} else {
		// 正常
		length := 30 + rand.Intn(70)
		result.CableStatus = "normal"
		result.FaultDesc = "正常"
		result.CableLength = fmt.Sprintf("约 %dm", length)
		result.Pairs = model.CablePairs{
			PairA: model.CablePair{Status: "ok", FaultDistance: fmt.Sprintf("%dm", length)},
			PairB: model.CablePair{Status: "ok", FaultDistance: fmt.Sprintf("%dm", length)},
			PairC: model.CablePair{Status: "ok", FaultDistance: fmt.Sprintf("%dm", length)},
			PairD: model.CablePair{Status: "ok", FaultDistance: fmt.Sprintf("%dm", length)},
		}
	}

	s.cableResults[taskID] = result
	return result, nil
}

// GetCableTestResult 获取电缆检测结果
func (s *DiagnosticService) GetCableTestResult(taskID string) (*model.CableTestResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.cableResults[taskID]
	return result, ok
}
