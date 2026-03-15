package model

// CableTestResult 虚拟电缆检测结果
type CableTestResult struct {
	TaskID          string         `json:"task_id"`
	PortID          string         `json:"port_id"`
	Status          string         `json:"status"` // completed/failed
	AdminStatus     string         `json:"admin_status"`
	LinkStatus      string         `json:"link_status"`
	CableStatus     string         `json:"cable_status"` // normal/open/short/cross/impedance
	FaultDesc       string         `json:"fault_description"`
	FaultDistance   string         `json:"fault_distance"`
	CableLength     string         `json:"cable_length"`
	Pairs           CablePairs     `json:"pairs"`
}

// CablePairs 线对状态
type CablePairs struct {
	PairA CablePair `json:"pair_a"` // 线对 A (针脚 1,2)
	PairB CablePair `json:"pair_b"` // 线对 B (针脚 3,6)
	PairC CablePair `json:"pair_c"` // 线对 C (针脚 4,5)
	PairD CablePair `json:"pair_d"` // 线对 D (针脚 7,8)
}

// CablePair 单一线对状态
type CablePair struct {
	Status        string `json:"status"`         // ok/open/short/cross
	FaultDistance string `json:"fault_distance"` // 故障距离（米）
}

// PortInfo 端口信息（用于电缆检测选择）
type PortInfo struct {
	PortID      string `json:"port_id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // electrical/optical
	AdminStatus string `json:"admin_status"`
	LinkStatus  string `json:"link_status"`
	Label       string `json:"label"`       // 显示标签
	Detectable  bool   `json:"detectable"`  // 是否可检测
	Hint        string `json:"hint"`        // 提示信息
}

// CableTestRequest 电缆检测请求
type CableTestRequest struct {
	PortID string `json:"port_id"`
}

// CableTestResponse 电缆检测响应
type CableTestResponse struct {
	Code    int              `json:"code"`
	Data    CableTestResult  `json:"data"`
	Message string           `json:"message,omitempty"`
}

// PortListResponse 端口列表响应
type PortListResponse struct {
	Code  int        `json:"code"`
	Data  PortListData `json:"data"`
}

// PortListData 端口列表数据
type PortListData struct {
	Ports []PortInfo `json:"ports"`
}
