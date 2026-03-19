package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFlowControl 获取流量控制配置
// GET /api/v1/config/flow-control
func (h *Handler) GetFlowControl(c *gin.Context) {
	// Mock 流量控制配置数据
	result := map[string]interface{}{
		"global": map[string]interface{}{
			"enabled":     true,
			"mode":        "auto",
			"backpressure": false,
			"pause_type":  "symmetric",
		},
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "enabled": true, "status": "up", "negotiation": "Full/On", "pause_direction": "both"},
			{"port_id": "GE1/0/2", "enabled": false, "status": "down", "negotiation": "-", "pause_direction": "none"},
			{"port_id": "GE1/0/3", "enabled": true, "status": "up", "negotiation": "Full/Off", "pause_direction": "none"},
			{"port_id": "GE1/0/4", "enabled": true, "status": "up", "negotiation": "Half/Backpressure", "pause_direction": "backpressure"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateFlowControlGlobal 更新全局流量控制配置
// PUT /api/v1/config/flow-control/global
func (h *Handler) UpdateFlowControlGlobal(c *gin.Context) {
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
		"message": "全局流量控制配置保存成功",
	})
}

// UpdateFlowControlPort 更新端口流量控制配置
// PUT /api/v1/config/flow-control/:port_id
func (h *Handler) UpdateFlowControlPort(c *gin.Context) {
	portID := c.Param("port_id")
	if portID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_ID",
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
		"message": "端口流量控制配置保存成功",
		"data": gin.H{
			"port_id": portID,
		},
	})
}
