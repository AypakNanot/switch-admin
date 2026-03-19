package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStackConfig 获取堆叠配置
// GET /api/v1/config/stack
func (h *Handler) GetStackConfig(c *gin.Context) {
	// Mock 堆叠配置数据
	result := map[string]interface{}{
		"stack_enabled": true,
		"stack_id":      1,
		"priority":      100,
		"members": []map[string]interface{}{
			{"member_id": 1, "mac": "00:1A:2B:3C:4D:5E", "priority": 100, "status": "master"},
			{"member_id": 2, "mac": "00:1A:2B:3C:4D:5F", "priority": 80, "status": "slave"},
			{"member_id": 3, "mac": "00:1A:2B:3C:4D:60", "priority": 60, "status": "slave"},
		},
		"stack_ports": []string{"stack-port 1/1", "stack-port 1/2"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdateStackConfig 更新堆叠配置
// PUT /api/v1/config/stack
func (h *Handler) UpdateStackConfig(c *gin.Context) {
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
		"message": "堆叠配置保存成功",
	})
}
