package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPortMirror 获取端口镜像配置
// GET /api/v1/config/port-mirror
func (h *Handler) GetPortMirror(c *gin.Context) {
	// Mock 端口镜像配置数据
	result := map[string]interface{}{
		"mirrors": []map[string]interface{}{
			{"id": 1, "name": "Mirror1", "type": "local", "source_ports": []string{"GE1/0/1", "GE1/0/2"}, "dest_port": "GE1/0/24", "direction": "both"},
			{"id": 2, "name": "Mirror2", "type": "local", "source_ports": []string{"GE1/0/5"}, "dest_port": "GE1/0/23", "direction": "rx"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdatePortMirror 更新端口镜像配置
// PUT /api/v1/config/port-mirror
func (h *Handler) UpdatePortMirror(c *gin.Context) {
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
		"message": "端口镜像配置保存成功",
	})
}
