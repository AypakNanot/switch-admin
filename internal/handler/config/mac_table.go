package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMacTable 获取 MAC 地址表
// GET /api/v1/config/mac-table
func (h *Handler) GetMacTable(c *gin.Context) {
	// Mock MAC 地址表数据
	result := map[string]interface{}{
		"total": 5,
		"items": []map[string]interface{}{
			{"mac_address": "00:1A:2B:3C:4D:5E", "vlan_id": 10, "port": "GE1/0/1", "type": "dynamic", "age_time": "280s"},
			{"mac_address": "00:1A:2B:3C:4D:5F", "vlan_id": 10, "port": "GE1/0/2", "type": "dynamic", "age_time": "150s"},
			{"mac_address": "AA:BB:CC:DD:EE:01", "vlan_id": 20, "port": "GE1/0/5", "type": "static", "age_time": "-"},
			{"mac_address": "00:00:00:00:00:00", "vlan_id": 1, "port": "-", "type": "blackhole", "age_time": "-"},
			{"mac_address": "11:22:33:44:55:66", "vlan_id": 100, "port": "GE1/0/24", "type": "dynamic", "age_time": "300s"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// AddMacTableEntry 添加静态 MAC 地址
// POST /api/v1/config/mac-table
func (h *Handler) AddMacTableEntry(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// Mock 添加成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "静态 MAC 地址添加成功",
	})
}

// DeleteMacTableEntry 删除静态 MAC 地址
// DELETE /api/v1/config/mac-table/:id
func (h *Handler) DeleteMacTableEntry(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	// Mock 删除成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "静态 MAC 地址删除成功",
		"data": gin.H{
			"deleted_id": id,
		},
	})
}
