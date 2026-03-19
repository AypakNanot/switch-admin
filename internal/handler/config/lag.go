package config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLinkAggregation 获取链路聚合列表
// GET /api/v1/config/link-aggregation
func (h *Handler) GetLinkAggregation(c *gin.Context) {
	// Mock 链路聚合数据
	result := map[string]interface{}{
		"aggregations": []map[string]interface{}{
			{"group_id": 1, "name": "LAG1", "mode": "LACP", "load_balance": "src-dst-ip", "member_ports": []string{"GE1/0/1", "GE1/0/2"}, "min_active": 1, "status": "normal"},
			{"group_id": 2, "name": "LAG2", "mode": "Static", "load_balance": "src-dst-mac", "member_ports": []string{"GE1/0/3", "GE1/0/4"}, "min_active": 1, "status": "normal"},
			{"group_id": 3, "name": "LAG3", "mode": "LACP", "load_balance": "src-dst-ip", "member_ports": []string{"GE1/0/5"}, "min_active": 1, "status": "degraded"},
		},
		"total": 3,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// CreateLinkAggregation 创建链路聚合组
// POST /api/v1/config/link-aggregation
func (h *Handler) CreateLinkAggregation(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// Mock 创建成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组创建成功",
	})
}

// UpdateLinkAggregation 更新链路聚合组
// PUT /api/v1/config/link-aggregation/:id
func (h *Handler) UpdateLinkAggregation(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// Mock 更新成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组更新成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// DeleteLinkAggregation 删除链路聚合组
// DELETE /api/v1/config/link-aggregation/:id
func (h *Handler) DeleteLinkAggregation(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	// Mock 删除成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组删除成功",
		"data": gin.H{
			"deleted_id": id,
		},
	})
}
