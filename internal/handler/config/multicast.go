package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMulticastConfig 获取组播配置
// GET /api/v1/config/multicast
func (h *Handler) GetMulticastConfig(c *gin.Context) {
	// Mock 组播配置数据
	result := map[string]interface{}{
		"igmp_snooping":       true,
		"igmp_version":        3,
		"fast_leave":          true,
		"unknown_multicast":   "discard",
		"multicast_vlan":      100,
		"static_groups": []map[string]interface{}{
			{"vlan_id": 10, "group": "239.1.1.1", "ports": []string{"GE1/0/1", "GE1/0/2"}},
			{"vlan_id": 20, "group": "239.2.2.2", "ports": []string{"GE1/0/5"}},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateMulticastConfig 更新组播配置
// PUT /api/v1/config/multicast
func (h *Handler) UpdateMulticastConfig(c *gin.Context) {
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
		"message": "组播配置保存成功",
	})
}
