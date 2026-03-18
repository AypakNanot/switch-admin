package model

// PortConfig 端口配置
type PortConfig struct {
	PortID      string `json:"port_id"`
	AdminStatus string `json:"admin_status"` // enable/disable
	LinkStatus  string `json:"link_status"`  // up/down
	SpeedDuplex string `json:"speed_duplex"` // 1000F/100F/auto
	FlowControl string `json:"flow_control"` // on/off
	Description string `json:"description"`
	Aggregation string `json:"aggregation"` // "-" 或 "Ag1" 等
}

// PortConfigListResponse 端口配置列表响应
type PortConfigListResponse struct {
	Ports []PortConfig `json:"ports"`
	Total int          `json:"total"`
}

// PortConfigRequest 端口配置请求
type PortConfigRequest struct {
	AdminStatus string `json:"admin_status"`
	Description string `json:"description"`
	SpeedDuplex string `json:"speed_duplex"`
	FlowControl string `json:"flow_control"`
}

// LinkAggregation 链路聚合组
type LinkAggregation struct {
	GroupID     int      `json:"group_id"`
	Name        string   `json:"name"`
	Mode        string   `json:"mode"`         // LACP/Static
	LoadBalance string   `json:"load_balance"` // src-dst-ip/src-dst-mac
	MemberPorts []string `json:"member_ports"`
	MinActive   int      `json:"min_active"`
	Status      string   `json:"status"` // normal/degraded/down
}

// LinkAggregationListResponse 链路聚合列表响应
type LinkAggregationListResponse struct {
	Aggregations []LinkAggregation `json:"aggregations"`
	Total        int               `json:"total"`
}

// LinkAggregationRequest 链路聚合请求
type LinkAggregationRequest struct {
	GroupID     int      `json:"group_id"`
	Mode        string   `json:"mode"`
	Description string   `json:"description"`
	LoadBalance string   `json:"load_balance"`
	MemberPorts []string `json:"member_ports"`
	MinActive   int      `json:"min_active"`
	LacpTimeout string   `json:"lacp_timeout"`
	Priority    int      `json:"priority"`
}
