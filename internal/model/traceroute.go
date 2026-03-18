package model

// TracerouteTask Traceroute 诊断任务
type TracerouteTask struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"` // running/completed/failed
	Target  string `json:"target"`
	VrfID   string `json:"vrf_id"`
	MaxHops int    `json:"max_hops"`
	Timeout int    `json:"timeout"`
	Probes  int    `json:"probes"`
}

// HopInfo 单跳信息
type HopInfo struct {
	Hop    int      `json:"hop"`
	IP     string   `json:"ip"`
	Times  []string `json:"times"`
	Status string   `json:"status"` // ok/timeout/destination/error
}

// TracerouteResponse Traceroute 任务响应
type TracerouteResponse struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	Target    string    `json:"target"`
	VrfID     string    `json:"vrf_id"`
	TotalHops int       `json:"total_hops"`
	Hops      []HopInfo `json:"hops"`
	Error     string    `json:"error,omitempty"`
}

// TracerouteRequest Traceroute 请求参数
type TracerouteRequest struct {
	VrfID   string `json:"vrf_id"`
	Target  string `json:"target"`
	MaxHops int    `json:"max_hops"`
	Timeout int    `json:"timeout"`
	Probes  int    `json:"probes"`
}
