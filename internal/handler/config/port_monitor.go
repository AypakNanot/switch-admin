package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPortMonitor 获取端口监测配置
// GET /api/v1/config/port-monitor
func (h *Handler) GetPortMonitor(c *gin.Context) {
	// Mock 端口监测配置数据
	result := map[string]interface{}{
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "enabled": true, "threshold": 80, "current_usage": 45.5, "status": "normal"},
			{"port_id": "GE1/0/2", "enabled": true, "threshold": 80, "current_usage": 82.3, "status": "warning"},
			{"port_id": "GE1/0/3", "enabled": false, "threshold": 80, "current_usage": 0, "status": "disabled"},
			{"port_id": "GE1/0/4", "enabled": true, "threshold": 90, "current_usage": 25.0, "status": "normal"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdatePortMonitor 更新端口监测配置
// PUT /api/v1/config/port-monitor/:port_id
func (h *Handler) UpdatePortMonitor(c *gin.Context) {
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
		"message": "端口监测配置保存成功",
		"data": gin.H{
			"port_id": portID,
		},
	})
}
