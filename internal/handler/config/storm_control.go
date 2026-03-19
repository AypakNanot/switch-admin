package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStormControl 获取风暴控制配置
// GET /api/v1/config/storm-control
func (h *Handler) GetStormControl(c *gin.Context) {
	// Mock 风暴控制配置数据
	result := map[string]interface{}{
		"global": map[string]interface{}{
			"enabled":     true,
			"mode":        "kbps",
			"storm_type":  "broadcast",
			"max_rate":    10000,
			"interval":    1,
			"action":      "drop",
		},
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "enabled": true, "storm_type": "broadcast", "max_rate": 5000, "current_rate": 1200, "status": "normal"},
			{"port_id": "GE1/0/2", "enabled": true, "storm_type": "multicast", "max_rate": 3000, "current_rate": 2800, "status": "warning"},
			{"port_id": "GE1/0/3", "enabled": false, "storm_type": "unknown-unicast", "max_rate": 2000, "current_rate": 0, "status": "disabled"},
			{"port_id": "GE1/0/4", "enabled": true, "storm_type": "broadcast", "max_rate": 5000, "current_rate": 150, "status": "normal"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateStormControlGlobal 更新全局风暴控制配置
// PUT /api/v1/config/storm-control/global
func (h *Handler) UpdateStormControlGlobal(c *gin.Context) {
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
		"message": "全局风暴控制配置保存成功",
	})
}

// UpdateStormControlPort 更新端口风暴控制配置
// PUT /api/v1/config/storm-control/:port_id
func (h *Handler) UpdateStormControlPort(c *gin.Context) {
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
		"message": "端口风暴控制配置保存成功",
		"data": gin.H{
			"port_id": portID,
		},
	})
}
