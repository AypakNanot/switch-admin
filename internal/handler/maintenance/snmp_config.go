package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switch-admin/internal/model"
)

// GetSNMPConfig 获取 SNMP 配置
// GET /api/v1/snmp/config
func (h *Handler) GetSNMPConfig(c *gin.Context) {
	config, err := h.service.GetSNMPConfig(c.Request.Context())
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

// UpdateSNMPConfig 更新 SNMP 配置
// PUT /api/v1/snmp/config
func (h *Handler) UpdateSNMPConfig(c *gin.Context) {
	var req model.SNMPConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateSNMPConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SNMP 配置保存成功",
	})
}

// GetSNMPCommunity 获取 SNMP 团体列表
// GET /api/v1/snmp/communities
func (h *Handler) GetSNMPCommunity(c *gin.Context) {
	communities, err := h.service.GetSNMPCommunities(c.Request.Context())
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
			"communities": communities,
			"total":       len(communities),
		},
	})
}

// AddSNMPCommunity 添加 SNMP 团体
// POST /api/v1/snmp/communities
func (h *Handler) AddSNMPCommunity(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Access      string `json:"access"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddCommunity(c.Request.Context(), req.Name, req.Access, req.Description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "团体添加成功",
	})
}

// DeleteSNMPCommunity 删除 SNMP 团体
// DELETE /api/v1/snmp/communities/:name
func (h *Handler) DeleteSNMPCommunity(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_NAME",
		})
		return
	}

	if err := h.service.DeleteCommunity(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "团体删除成功",
	})
}
