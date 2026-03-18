package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switch-admin/internal/model"
)

// GetARPConfig 获取 ARP 防护配置
// GET /api/v1/security/arp/config
func (h *Handler) GetARPConfig(c *gin.Context) {
	config, err := h.service.GetARPConfig(c.Request.Context())
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

// UpdateARPConfig 更新 ARP 防护配置
// PUT /api/v1/security/arp/config
func (h *Handler) UpdateARPConfig(c *gin.Context) {
	var req model.ARPConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateARPConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ARP 防护配置保存成功",
	})
}
