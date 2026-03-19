package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSTPConfig 获取 STP 配置
// GET /api/v1/config/stp
func (h *Handler) GetSTPConfig(c *gin.Context) {
	// Mock STP 配置数据
	result := map[string]interface{}{
		"enabled":       true,
		"mode":          "RSTP",
		"priority":      32768,
		"hello_time":    2,
		"max_age":       20,
		"forward_delay": 15,
		"root_bridge": map[string]interface{}{
			"mac":      "00:1A:2B:3C:4D:5E",
			"priority": 4096,
			"cost":     0,
		},
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "enabled": true, "priority": 128, "cost": 20000, "state": "forwarding", "role": "designated"},
			{"port_id": "GE1/0/2", "enabled": true, "priority": 128, "cost": 20000, "state": "blocking", "role": "alternate"},
			{"port_id": "GE1/0/3", "enabled": false, "priority": 128, "cost": 20000, "state": "disabled", "role": "-"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateSTPConfig 更新 STP 配置
// PUT /api/v1/config/stp
func (h *Handler) UpdateSTPConfig(c *gin.Context) {
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
		"message": "STP 配置保存成功",
	})
}

// GetSTPStatus 获取 STP 状态
// GET /api/v1/config/stp/status
func (h *Handler) GetSTPStatus(c *gin.Context) {
	// Mock STP 状态数据
	result := map[string]interface{}{
		"root_mac":      "00:1A:2B:3C:4D:5E",
		"root_priority": 4096,
		"root_cost":     0,
		"this_bridge": map[string]interface{}{
			"mac":      "00:1A:2B:3C:4D:5F",
			"priority": 32768,
		},
		"topology_changes": 5,
		"last_change":      "2024-01-15 10:30:00",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}
