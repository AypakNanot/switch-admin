package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetVLANList 获取 VLAN 列表
// GET /api/v1/network/vlans
func (h *Handler) GetVLANList(c *gin.Context) {
	result, err := h.service.GetVLANList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"vlans": result.VLANs,
			"total": result.Total,
		},
	})
}

// CreateVLAN 创建 VLAN
// POST /api/v1/network/vlans
func (h *Handler) CreateVLAN(c *gin.Context) {
	var req model.VLANRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.CreateVLAN(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 创建成功",
	})
}

// UpdateVLAN 更新 VLAN
// PUT /api/v1/network/vlans/:id
func (h *Handler) UpdateVLAN(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.VLANRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateVLAN(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 更新成功",
	})
}

// DeleteVLAN 删除 VLAN
// DELETE /api/v1/network/vlans/:id
func (h *Handler) DeleteVLAN(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	if err := h.service.DeleteVLAN(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 删除成功",
	})
}

// DeleteVLANs 批量删除 VLAN
// DELETE /api/v1/network/vlans
func (h *Handler) DeleteVLANs(c *gin.Context) {
	var req struct {
		Ids []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.DeleteVLANs(c.Request.Context(), req.Ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "VLAN 批量删除成功",
	})
}

// AddVLANPort 添加 VLAN 端口
// POST /api/v1/network/vlans/:id/ports
func (h *Handler) AddVLANPort(c *gin.Context) {
	idStr := c.Param("id")
	vlanID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req struct {
		Port string `json:"port"`
		Mode string `json:"mode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddVLANPort(c.Request.Context(), vlanID, req.Port, req.Mode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口添加成功",
	})
}

// RemoveVLANPort 移除 VLAN 端口
// DELETE /api/v1/network/vlans/:id/ports/:port
func (h *Handler) RemoveVLANPort(c *gin.Context) {
	idStr := c.Param("id")
	vlanID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	port := c.Param("port")
	if port == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT",
		})
		return
	}

	if err := h.service.RemoveVLANPort(c.Request.Context(), vlanID, port); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口移除成功",
	})
}
