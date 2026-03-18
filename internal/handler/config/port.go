package config

import (
	"net/http"

	"switch-admin/internal/model"

	"github.com/gin-gonic/gin"
)

// GetPorts 获取端口列表
// GET /api/v1/config/ports
func (h *Handler) GetPorts(c *gin.Context) {
	result, err := h.service.GetPortList(c.Request.Context())
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
			"ports": result.Ports,
			"total": result.Total,
		},
	})
}

// GetPortDetail 获取端口详情
// GET /api/v1/config/ports/:port_id
func (h *Handler) GetPortDetail(c *gin.Context) {
	portID := c.Param("port_id")
	if portID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_ID",
		})
		return
	}

	result, err := h.service.GetPortDetail(c.Request.Context(), portID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// UpdatePort 更新端口配置
// PUT /api/v1/config/ports/:port_id
func (h *Handler) UpdatePort(c *gin.Context) {
	portID := c.Param("port_id")
	if portID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_ID",
		})
		return
	}

	var req model.PortConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdatePort(c.Request.Context(), portID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口配置保存成功",
	})
}
