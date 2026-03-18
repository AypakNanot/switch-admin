package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetSTPConfig 获取 STP 配置
// GET /api/v1/network/stp/config
func (h *Handler) GetSTPConfig(c *gin.Context) {
	result, err := h.service.GetSTPConfig(c.Request.Context())
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

// UpdateSTPConfig 更新 STP 配置
// PUT /api/v1/network/stp/config
func (h *Handler) UpdateSTPConfig(c *gin.Context) {
	var req model.STPConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateSTPConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "STP 配置更新成功",
	})
}

// GetSTPStatus 获取 STP 状态
// GET /api/v1/network/stp/status
func (h *Handler) GetSTPStatus(c *gin.Context) {
	result, err := h.service.GetSTPStatus(c.Request.Context())
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
