package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPoEConfig 获取 PoE 配置
// GET /api/v1/config/poe
func (h *Handler) GetPoEConfig(c *gin.Context) {
	// Mock PoE 配置数据
	result := map[string]interface{}{
		"global": map[string]interface{}{
			"enabled":      true,
			"power_budget": 370,
			"power_used":   125,
			"power_avail":  245,
		},
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "enabled": true, "status": "delivering", "power": 15.4, "class": 3, "voltage": 53.2, "current": 0.29},
			{"port_id": "GE1/0/2", "enabled": true, "status": "delivering", "power": 3.2, "class": 2, "voltage": 52.8, "current": 0.06},
			{"port_id": "GE1/0/3", "enabled": false, "status": "disabled", "power": 0, "class": "-", "voltage": 0, "current": 0},
			{"port_id": "GE1/0/4", "enabled": true, "status": "searching", "power": 0, "class": "-", "voltage": 0, "current": 0},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdatePoEConfig 更新 PoE 配置
// PUT /api/v1/config/poe/:port_id
func (h *Handler) UpdatePoEConfig(c *gin.Context) {
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
		"message": "PoE 配置保存成功",
		"data": gin.H{
			"port_id": portID,
		},
	})
}
