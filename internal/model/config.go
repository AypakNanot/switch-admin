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

// StormControlConfig 风暴控制配置
type StormControlConfig struct {
	Enabled  bool              `json:"enabled"`
	Mode     string            `json:"mode"` // kbps/pps
	StormType string          `json:"storm_type"` // broadcast/multicast/unknown-unicast
	MaxRate  int               `json:"max_rate"`
	Interval int               `json:"interval"`
	Action   string            `json:"action"` // drop/shutdown
	Ports    []StormControlPort `json:"ports"`
}

// StormControlPort 端口风暴控制配置
type StormControlPort struct {
	PortID      string `json:"port_id"`
	Enabled     bool   `json:"enabled"`
	StormType   string `json:"storm_type"`
	MaxRate     int    `json:"max_rate"`
	CurrentRate int    `json:"current_rate"`
	Status      string `json:"status"` // normal/warning/exceeded/disabled
}

// StormControlRequest 风暴控制请求
type StormControlRequest struct {
	Enabled  bool   `json:"enabled"`
	Mode     string `json:"mode"`
	StormType string `json:"storm_type"`
	MaxRate  int    `json:"max_rate"`
	Interval int    `json:"interval"`
	Action   string `json:"action"`
}

// StormControlPortRequest 端口风暴控制请求
type StormControlPortRequest struct {
	Enabled   bool   `json:"enabled"`
	StormType string `json:"storm_type"`
	MaxRate   int    `json:"max_rate"`
}

// FlowControlConfig 流控配置
type FlowControlConfig struct {
	Enabled      bool               `json:"enabled"`
	Mode         string             `json:"mode"` // auto/manual
	Backpressure bool               `json:"backpressure"`
	PauseType    string             `json:"pause_type"` // symmetric/asymmetric/none
	Ports        []FlowControlPort  `json:"ports"`
}

// FlowControlPort 端口流控配置
type FlowControlPort struct {
	PortID         string `json:"port_id"`
	Enabled        bool   `json:"enabled"`
	Status         string `json:"status"` // up/down
	Negotiation    string `json:"negotiation"`
	PauseDirection string `json:"pause_direction"` // both/tx/rx/none/backpressure
}

// FlowControlRequest 流控请求
type FlowControlRequest struct {
	Enabled      bool   `json:"enabled"`
	Mode         string `json:"mode"`
	Backpressure bool   `json:"backpressure"`
	PauseType    string `json:"pause_type"`
}

// FlowControlPortRequest 端口流控请求
type FlowControlPortRequest struct {
	Enabled bool   `json:"enabled"`
	Mode    string `json:"mode"`
}

// PortIsolationConfig 端口隔离配置
type PortIsolationConfig struct {
	Enabled    bool                `json:"enabled"`
	IsolationGroups []PortIsolationGroup `json:"isolation_groups"`
}

// PortIsolationGroup 端口隔离组
type PortIsolationGroup struct {
	GroupID int      `json:"group_id"`
	Name    string   `json:"name"`
	Ports   []string `json:"ports"`
	IsolationMode string `json:"isolation_mode"` // all/l2
}

// PortIsolationRequest 端口隔离请求
type PortIsolationRequest struct {
	GroupID       int      `json:"group_id"`
	Name          string   `json:"name"`
	Ports         []string `json:"ports"`
	IsolationMode string   `json:"isolation_mode"`
}

// PortMonitorConfig 端口镜像配置
type PortMonitorConfig struct {
	Sessions []PortMirrorSession `json:"sessions"`
}

// PortMirrorSession 端口镜像会话
type PortMirrorSession struct {
	SessionID     int      `json:"session_id"`
	Name          string   `json:"name"`
	MonitorPort   string   `json:"monitor_port"`
	SourcePorts   []string `json:"source_ports"`
	Direction     string   `json:"direction"` // ingress/egress/both
	Enabled       bool     `json:"enabled"`
}

// PortMirrorRequest 端口镜像请求
type PortMirrorRequest struct {
	SessionID   int      `json:"session_id"`
	Name        string   `json:"name"`
	MonitorPort string   `json:"monitor_port"`
	SourcePorts []string `json:"source_ports"`
	Direction   string   `json:"direction"`
	Enabled     bool     `json:"enabled"`
}

// MacTableEntry MAC 表项
type MacTableEntry struct {
	VLANID    int    `json:"vlan_id"`
	MACAddress string `json:"mac_address"`
	PortID    string `json:"port_id"`
	Type      string `json:"type"` // dynamic/static
	AgingTime int    `json:"aging_time"`
}

// MacTableListResponse MAC 表列表响应
type MacTableListResponse struct {
	Entries []MacTableEntry `json:"entries"`
	Total   int             `json:"total"`
}

// ERPSConfig ERPS 配置
type ERPSConfig struct {
	Enabled        bool     `json:"enabled"`
	RingID         int      `json:"ring_id"`
	ControlVLAN    int      `json:"control_vlan"`
	DataVLANs      []int    `json:"data_vlans"`
	Role           string   `json:"role"` // auto/master/slave
	WTR            int      `json:"wtr"` // Wait To Restore time in minutes
	RingStatus     string   `json:"ring_status"` // normal/forced/switched
	ActiveTopology string   `json:"active_topology"` // clockwise/counter-clockwise
}

// ERPSRequest ERPS 请求
type ERPSRequest struct {
	Enabled     bool   `json:"enabled"`
	RingID      int    `json:"ring_id"`
	ControlVLAN int    `json:"control_vlan"`
	DataVLANs   []int  `json:"data_vlans"`
	Role        string `json:"role"`
	WTR         int    `json:"wtr"`
}

// MulticastConfig 组播配置
type MulticastConfig struct {
	Enabled       bool     `json:"enabled"`
	Mode          string   `json:"mode"` // igmp-snoop/mvr
	RouterPorts   []string `json:"router_ports"`
	HostPorts     []string `json:"host_ports"`
	FastLeave     bool     `json:"fast_leave"`
}

// MulticastRequest 组播请求
type MulticastRequest struct {
	Enabled   bool     `json:"enabled"`
	Mode      string   `json:"mode"`
	RouterPorts []string `json:"router_ports"`
	HostPorts []string   `json:"host_ports"`
	FastLeave bool     `json:"fast_leave"`
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage       int    `json:"cpu_usage"` // percentage
	MemoryUsage    int    `json:"memory_usage"` // percentage
	Temperature    int    `json:"temperature"` // Celsius
	FanStatus      string `json:"fan_status"` // normal/abnormal
	PowerStatus    string `json:"power_status"` // normal/abnormal
	Uptime         string `json:"uptime"`
	FlashUsage     int    `json:"flash_usage"` // percentage
	DRAMSize       int    `json:"dram_size"` // MB
	FlashSize      int    `json:"flash_size"` // MB
}

// StackConfig 堆叠配置
type StackConfig struct {
	Enabled      bool         `json:"enabled"`
	MasterID     int          `json:"master_id"`
	MemberCount  int          `json:"member_count"`
	Members      []StackMember `json:"members"`
	Topology     string       `json:"topology"` // chain/ring
}

// StackMember 堆叠成员
type StackMember struct {
	MemberID   int    `json:"member_id"`
	MACAddress string `json:"mac_address"`
	Priority   int    `json:"priority"`
	Role       string `json:"role"` // master/slave
	Status     string `json:"status"` // active/inactive/offline
}

// StackRequest 堆叠请求
type StackRequest struct {
	MemberID int `json:"member_id"`
	Priority int `json:"priority"`
}
