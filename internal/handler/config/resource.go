package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetResource 获取资源使用情况
// GET /api/v1/config/resource
func (h *Handler) GetResource(c *gin.Context) {
	// Mock 资源使用情况数据
	result := map[string]interface{}{
		"cpu_usage":       15.5,
		"memory_usage":    42.3,
		"flash_usage":     35.0,
		"mac_table": map[string]interface{}{
			"used":  1250,
			"total": 8000,
		},
		"arp_table": map[string]interface{}{
			"used":  350,
			"total": 2000,
		},
		"routing_table": map[string]interface{}{
			"used":  120,
			"total": 1000,
		},
		"vlan_table": map[string]interface{}{
			"used":  25,
			"total": 4094,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}
