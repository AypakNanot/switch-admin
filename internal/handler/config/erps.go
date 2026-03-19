package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetERPSConfig 获取 ERPS 配置
// GET /api/v1/config/erps
func (h *Handler) GetERPSConfig(c *gin.Context) {
	// Mock ERPS 配置数据
	result := map[string]interface{}{
		"enabled":       true,
		"ring_id":       1,
		"control_vlan":  4000,
		"data_vlans":    []int{10, 20, 30},
		"role":          "auto",
		"wtr":           5,
		"ring_status":   "normal",
		"active_topology": "clockwise",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateERPSConfig 更新 ERPS 配置
// PUT /api/v1/config/erps
func (h *Handler) UpdateERPSConfig(c *gin.Context) {
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
		"message": "ERPS 配置保存成功",
	})
}
