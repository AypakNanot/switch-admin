package model

// PingTask Ping 诊断任务
type PingTask struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"` // running/completed/failed
	Target   string `json:"target"`
	VrfID    string `json:"vrf_id"`
	Count    int    `json:"count"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}

// PingResult Ping 单次探测结果
type PingResult struct {
	Seq    int    `json:"seq"`
	Time   string `json:"time"`
	TTL    int    `json:"ttl"`
	Status string `json:"status"` // success/timeout/error
}

// PingStatistics Ping 统计信息
type PingStatistics struct {
	Sent     int    `json:"sent"`
	Received int    `json:"received"`
	LossRate string `json:"loss_rate"`
	MinTime  string `json:"min_time"`
	AvgTime  string `json:"avg_time"`
	MaxTime  string `json:"max_time"`
}

// PingTaskResponse Ping 任务响应
type PingTaskResponse struct {
	TaskID     string         `json:"task_id"`
	Status     string         `json:"status"`
	Target     string         `json:"target"`
	VrfID      string         `json:"vrf_id"`
	Results    []PingResult   `json:"results"`
	Statistics PingStatistics `json:"statistics"`
	Error      string         `json:"error,omitempty"`
}

// PingRequest Ping 请求参数
type PingRequest struct {
	VrfID  string `json:"vrf_id"`
	Target string `json:"target"`
	Count  int    `json:"count"`
	Timeout int   `json:"timeout"`
	Interval int `json:"interval"`
}
