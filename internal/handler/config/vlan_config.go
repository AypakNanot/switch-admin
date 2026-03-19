package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVLANConfig 获取 VLAN 配置
// GET /api/v1/config/vlan
func (h *Handler) GetVLANConfig(c *gin.Context) {
	// Mock VLAN 配置数据
	result := map[string]interface{}{
		"vlans": []map[string]interface{}{
			{"vlan_id": 1, "name": "default", "status": "active", "ports": []string{"GE1/0/1", "GE1/0/2", "GE1/0/3", "GE1/0/4"}, "type": "static"},
			{"vlan_id": 10, "name": "users", "status": "active", "ports": []string{"GE1/0/5", "GE1/0/6"}, "type": "static"},
			{"vlan_id": 20, "name": "servers", "status": "active", "ports": []string{"GE1/0/7", "GE1/0/8"}, "type": "static"},
			{"vlan_id": 100, "name": "management", "status": "active", "ports": []string{"GE1/0/24"}, "type": "static"},
		},
		"total": 4,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// CreateVLAN 创建 VLAN
// POST /api/v1/config/vlan
func (h *Handler) CreateVLAN(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// Mock 创建成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 创建成功",
		"data": gin.H{
			"vlan_id": req["vlan_id"],
		},
	})
}

// UpdateVLAN 更新 VLAN 配置
// PUT /api/v1/config/vlan/:vlan_id
func (h *Handler) UpdateVLAN(c *gin.Context) {
	vlanID := c.Param("vlan_id")
	if vlanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_VLAN_ID",
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
		"message": "VLAN 配置更新成功",
		"data": gin.H{
			"vlan_id": vlanID,
		},
	})
}

// DeleteVLAN 删除 VLAN
// DELETE /api/v1/config/vlan/:vlan_id
func (h *Handler) DeleteVLAN(c *gin.Context) {
	vlanID := c.Param("vlan_id")
	if vlanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_VLAN_ID",
		})
		return
	}

	// Mock 删除成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 删除成功",
		"data": gin.H{
			"deleted_vlan_id": vlanID,
		},
	})
}
