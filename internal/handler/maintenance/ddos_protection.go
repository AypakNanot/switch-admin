package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switch-admin/internal/model"
)

// GetDDoSConfig 获取 DDoS 防护配置
// GET /api/v1/security/ddos/config
func (h *Handler) GetDDoSConfig(c *gin.Context) {
	config, err := h.service.GetDDoSConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": config,
	})
}

// UpdateDDoSConfig 更新 DDoS 防护配置
// PUT /api/v1/security/ddos/config
func (h *Handler) UpdateDDoSConfig(c *gin.Context) {
	var req model.DDoSConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateDDoSConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "DDoS 防护配置保存成功",
	})
}
