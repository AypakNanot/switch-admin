package model

// RouteEntry IPv4 路由表条目
type RouteEntry struct {
	DestIP     string `json:"dest_ip"`    // 目的 IP 地址
	DestMask   string `json:"dest_mask"`  // 目的 IP 掩码
	Protocol   string `json:"protocol"`   // 协议类型：Static/OSPF/RIP/BGP/Connected/Local
	OutPort    string `json:"out_port"`   // 出端口
	NextHop    string `json:"next_hop"`   // 下一跳
	Metric     int    `json:"metric"`     // 度量值
	Preference int    `json:"preference"` // 优先级
}

// StaticRoute 静态路由配置
type StaticRoute struct {
	ID         string `json:"id"`          // 路由 ID
	DestIP     string `json:"dest_ip"`     // 目的 IP 地址
	DestMask   string `json:"dest_mask"`   // 目的 IP 掩码
	NextHop    string `json:"next_hop"`    // 下一跳
	Distance   int    `json:"distance"`    // 管理距离 (1-255)
	Status     string `json:"status"`      // 状态：active/inactive/warning
	StatusDesc string `json:"status_desc"` // 状态说明
}

// RouteTableQuery 路由表查询参数
type RouteTableQuery struct {
	DestIP   string `json:"dest_ip"`   // 目的 IP（支持模糊匹配）
	Protocol string `json:"protocol"`  // 协议类型筛选
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页条数
}

// RouteTableResponse 路由表查询响应
type RouteTableResponse struct {
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
	Items      []RouteEntry `json:"items"`
}
