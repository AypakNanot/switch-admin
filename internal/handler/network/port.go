package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetPortList 获取端口列表
// GET /api/v1/network/ports
func (h *Handler) GetPortList(c *gin.Context) {
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
// GET /api/v1/network/ports/:name
func (h *Handler) GetPortDetail(c *gin.Context) {
	portName := c.Param("name")
	if portName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_NAME",
		})
		return
	}

	result, err := h.service.GetPortDetail(c.Request.Context(), portName)
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
// PUT /api/v1/network/ports/:name
func (h *Handler) UpdatePort(c *gin.Context) {
	portName := c.Param("name")
	if portName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_NAME",
		})
		return
	}

	var req model.PortUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdatePort(c.Request.Context(), portName, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口配置更新成功",
	})
}

// ResetPort 重置端口配置
// POST /api/v1/network/ports/:name/reset
func (h *Handler) ResetPort(c *gin.Context) {
	portName := c.Param("name")
	if portName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_NAME",
		})
		return
	}

	if err := h.service.ResetPort(c.Request.Context(), portName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口配置已重置",
	})
}

// RestartPort 重启端口
// POST /api/v1/network/ports/:name/restart
func (h *Handler) RestartPort(c *gin.Context) {
	portName := c.Param("name")
	if portName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT_NAME",
		})
		return
	}

	if err := h.service.RestartPort(c.Request.Context(), portName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口已重启",
	})
}
