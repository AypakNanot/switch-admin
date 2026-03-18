package model

// VLAN VLAN 配置
type VLAN struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Ports []string `json:"ports"`
	Status string  `json:"status"`
}

// VLANListResponse VLAN 列表响应
type VLANListResponse struct {
	VLANs []VLAN `json:"vlans"`
	Total int    `json:"total"`
}

// VLANRequest VLAN 请求
type VLANRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Port 端口配置
type Port struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Speed  string `json:"speed"`
	Duplex string `json:"duplex"`
	VLAN   int    `json:"vlan"`
	Type   string `json:"type"`
}

// NetworkPortListResponse 端口列表响应
type NetworkPortListResponse struct {
	Ports []Port `json:"ports"`
	Total int    `json:"total"`
}

// PortDetail 端口详情
type PortDetail struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Speed       string `json:"speed"`
	Duplex      string `json:"duplex"`
	VLAN        int    `json:"vlan"`
	Type        string `json:"type"`
	MAC         string `json:"mac"`
	Description string `json:"description"`
	RXBytes     int64  `json:"rx_bytes"`
	TXBytes     int64  `json:"tx_bytes"`
	RXErrors    int64  `json:"rx_errors"`
	TXErrors    int64  `json:"tx_errors"`
}

// PortUpdateRequest 端口更新请求
type PortUpdateRequest struct {
	Description string `json:"description"`
	Speed       string `json:"speed"`
	Duplex      string `json:"duplex"`
	VLAN        int    `json:"vlan"`
	Enabled     bool   `json:"enabled"`
}

// LAG 链路聚合组
type LAG struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Ports  []string `json:"ports"`
	Status string   `json:"status"`
	Mode   string   `json:"mode"`
}

// LAGListResponse 链路聚合组列表响应
type LAGListResponse struct {
	LAGs  []LAG `json:"lags"`
	Total int   `json:"total"`
}

// LAGRequest 链路聚合组请求
type LAGRequest struct {
	Name  string   `json:"name"`
	Ports []string `json:"ports"`
	Mode  string   `json:"mode"`
}

// STPConfig STP 配置
type STPConfig struct {
	Enabled    bool   `json:"enabled"`
	Mode       string `json:"mode"`
	Priority   int    `json:"priority"`
	RootBridge string `json:"root_bridge"`
}

// STPConfigRequest STP 配置请求
type STPConfigRequest struct {
	Enabled  bool   `json:"enabled"`
	Mode     string `json:"mode"`
	Priority int    `json:"priority"`
}

// STPStatus STP 状态
type STPStatus struct {
	Enabled    bool           `json:"enabled"`
	Mode       string         `json:"mode"`
	RootBridge string         `json:"root_bridge"`
	RootPort   string         `json:"root_port"`
	PortStates []STPPortState `json:"port_states"`
}

// STPPortState STP 端口状态
type STPPortState struct {
	Port  string `json:"port"`
	State string `json:"state"`
	Role  string `json:"role"`
}

// ACL ACL 配置
type ACL struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Rules int    `json:"rules"`
	Status string `json:"status"`
}

// ACLListResponse ACL 列表响应
type ACLListResponse struct {
	ACLs  []ACL `json:"acls"`
	Total int   `json:"total"`
}

// ACLRequest ACL 请求
type ACLRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ACLRule ACL 规则
type ACLRule struct {
	ID          int    `json:"id"`
	Action      string `json:"action"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Port        string `json:"port"`
	Protocol    string `json:"protocol"`
}

// ACLRuleListResponse ACL 规则列表响应
type ACLRuleListResponse struct {
	Rules []ACLRule `json:"rules"`
	Total int       `json:"total"`
}

// ACLRuleRequest ACL 规则请求
type ACLRuleRequest struct {
	Action      string `json:"action"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Port        string `json:"port"`
	Protocol    string `json:"protocol"`
}
