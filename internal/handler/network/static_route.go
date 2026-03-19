package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetStaticRoutes 获取静态路由列表
// GET /api/v1/routes/static
func (h *Handler) GetStaticRoutes(c *gin.Context) {
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

// GetStaticRoute 获取单条静态路由
// GET /api/v1/routes/static/:id
func (h *Handler) GetStaticRoute(c *gin.Context) {
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

// CreateStaticRoute 创建静态路由
// POST /api/v1/routes/static
func (h *Handler) CreateStaticRoute(c *gin.Context) {
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

// UpdateStaticRoute 更新静态路由
// PUT /api/v1/routes/static/:id
func (h *Handler) UpdateStaticRoute(c *gin.Context) {
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

// DeleteStaticRoute 删除静态路由
// DELETE /api/v1/routes/static/:id
func (h *Handler) DeleteStaticRoute(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"deleted_id": id,
		},
	})
}
