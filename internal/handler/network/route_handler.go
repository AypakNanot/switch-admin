package network

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
	"switch-admin/internal/service"
)

// RouteHandler 路由管理处理器
type RouteHandler struct {
	diagnosticService *service.DiagnosticService
}

// NewRouteHandler 创建路由处理器
func NewRouteHandler() *RouteHandler {
	return &RouteHandler{
		diagnosticService: service.GetDiagnosticService(),
	}
}

// GetRouteTable 获取 IPv4 路由表
// GET /api/v1/routes/table
func (h *RouteHandler) GetRouteTable(c *gin.Context) {
	// 获取查询参数
	destIP := c.DefaultQuery("dest_ip", "")
	protocol := c.DefaultQuery("protocol", "")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	// Mock 路由表数据
	allRoutes := []model.RouteEntry{
		{DestIP: "0.0.0.0", DestMask: "0.0.0.0", Protocol: "Static", OutPort: "vlan1", NextHop: "192.168.1.1", Metric: 1, Preference: 60},
		{DestIP: "10.10.1.0", DestMask: "255.255.255.0", Protocol: "Connected", OutPort: "eth0/1", NextHop: "0.0.0.0", Metric: 0, Preference: 0},
		{DestIP: "10.10.2.0", DestMask: "255.255.255.0", Protocol: "Connected", OutPort: "eth0/2", NextHop: "0.0.0.0", Metric: 0, Preference: 0},
		{DestIP: "10.10.3.0", DestMask: "255.255.255.0", Protocol: "OSPF", OutPort: "eth0/3", NextHop: "10.10.1.1", Metric: 10, Preference: 110},
		{DestIP: "10.10.4.0", DestMask: "255.255.255.0", Protocol: "OSPF", OutPort: "eth0/3", NextHop: "10.10.1.1", Metric: 20, Preference: 110},
		{DestIP: "10.10.5.0", DestMask: "255.255.255.0", Protocol: "Static", OutPort: "vlan10", NextHop: "192.168.10.1", Metric: 1, Preference: 60},
		{DestIP: "172.16.0.0", DestMask: "255.255.0.0", Protocol: "RIP", OutPort: "eth0/4", NextHop: "10.10.2.1", Metric: 2, Preference: 120},
		{DestIP: "192.168.0.0", DestMask: "255.255.0.0", Protocol: "Static", OutPort: "vlan1", NextHop: "192.168.1.1", Metric: 1, Preference: 60},
		{DestIP: "192.168.1.0", DestMask: "255.255.255.0", Protocol: "Connected", OutPort: "vlan1", NextHop: "0.0.0.0", Metric: 0, Preference: 0},
		{DestIP: "192.168.100.0", DestMask: "255.255.255.0", Protocol: "BGP", OutPort: "eth0/1", NextHop: "10.10.1.254", Metric: 100, Preference: 255},
	}

	// 过滤
	filteredRoutes := make([]model.RouteEntry, 0)
	for _, route := range allRoutes {
		// 协议过滤
		if protocol != "" && protocol != "全部" && route.Protocol != protocol {
			continue
		}
		// IP 过滤（支持模糊匹配）
		if destIP != "" {
			if !strings.HasPrefix(route.DestIP, destIP) && !strings.HasPrefix(route.DestIP, destIP+".") {
				continue
			}
		}
		filteredRoutes = append(filteredRoutes, route)
	}

	total := len(filteredRoutes)
	totalPages := (total + pageSize - 1) / pageSize

	// 分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedRoutes := filteredRoutes[start:end]

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"items":       pagedRoutes,
		},
	})
}

// GetStaticRoutes 获取静态路由列表
// GET /api/v1/routes/static
func (h *RouteHandler) GetStaticRoutes(c *gin.Context) {
	// Mock 静态路由数据
	staticRoutes := []model.StaticRoute{
		{ID: "1", DestIP: "0.0.0.0", DestMask: "0.0.0.0", NextHop: "192.168.1.1", Distance: 1, Status: "active", StatusDesc: ""},
		{ID: "2", DestIP: "10.10.5.0", DestMask: "255.255.255.0", NextHop: "192.168.10.1", Distance: 1, Status: "active", StatusDesc: ""},
		{ID: "3", DestIP: "192.168.0.0", DestMask: "255.255.0.0", NextHop: "192.168.1.1", Distance: 5, Status: "active", StatusDesc: ""},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": staticRoutes,
		},
	})
}

// CreateStaticRoute 创建静态路由
// POST /api/v1/routes/static
func (h *RouteHandler) CreateStaticRoute(c *gin.Context) {
	var req model.StaticRoute
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid JSON format",
		})
		return
	}

	// 基本校验
	if req.DestIP == "" || req.DestMask == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "目的 IP 和掩码不能为空",
		})
		return
	}

	// Mock 创建成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "静态路由配置成功",
		"data": gin.H{
			"id":          "4",
			"dest_ip":     req.DestIP,
			"dest_mask":   req.DestMask,
			"next_hop":    req.NextHop,
			"distance":    req.Distance,
			"status":      "active",
			"status_desc": "",
		},
	})
}

// DeleteStaticRoute 删除静态路由
// DELETE /api/v1/routes/static/:id
func (h *RouteHandler) DeleteStaticRoute(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"deleted_id": id,
		},
	})
}

// UpdateStaticRoute 更新静态路由
// PUT /api/v1/routes/static/:id
func (h *RouteHandler) UpdateStaticRoute(c *gin.Context) {
	id := c.Param("id")
	var req model.StaticRoute
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid JSON format",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data": gin.H{
			"id":        id,
			"dest_ip":   req.DestIP,
			"dest_mask": req.DestMask,
			"next_hop":  req.NextHop,
			"distance":  req.Distance,
		},
	})
}

// GetStaticRoute 获取单条静态路由
// GET /api/v1/routes/static/:id
func (h *RouteHandler) GetStaticRoute(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": model.StaticRoute{
			ID:       id,
			DestIP:   "0.0.0.0",
			DestMask: "0.0.0.0",
			NextHop:  "192.168.1.1",
			Distance: 1,
			Status:   "active",
		},
	})
}
