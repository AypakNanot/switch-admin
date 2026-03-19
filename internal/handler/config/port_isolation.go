package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPortIsolation 获取端口隔离配置
// GET /api/v1/config/port-isolation
func (h *Handler) GetPortIsolation(c *gin.Context) {
	// Mock 端口隔离配置数据
	result := map[string]interface{}{
		"global": map[string]interface{}{
			"enabled": true,
			"mode":    "vlan",
		},
		"ports": []map[string]interface{}{
			{"port_id": "GE1/0/1", "isolation_group": 1, "isolated_ports": []string{"GE1/0/2", "GE1/0/3"}, "uplink_ports": []string{"GE1/0/24"}},
			{"port_id": "GE1/0/2", "isolation_group": 1, "isolated_ports": []string{"GE1/0/1", "GE1/0/3"}, "uplink_ports": []string{"GE1/0/24"}},
			{"port_id": "GE1/0/3", "isolation_group": 2, "isolated_ports": []string{}, "uplink_ports": []string{"GE1/0/24"}},
			{"port_id": "GE1/0/4", "isolation_group": 0, "isolated_ports": []string{}, "uplink_ports": []string{}},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdatePortIsolation 更新端口隔离配置
// PUT /api/v1/config/port-isolation/:port_id
func (h *Handler) UpdatePortIsolation(c *gin.Context) {
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
		"message": "端口隔离配置保存成功",
		"data": gin.H{
			"port_id": portID,
		},
	})
}
