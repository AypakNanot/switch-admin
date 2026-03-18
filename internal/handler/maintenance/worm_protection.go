package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switch-admin/internal/model"
)

// GetWormRules 获取蠕虫规则列表
// GET /api/v1/security/worm/rules
func (h *Handler) GetWormRules(c *gin.Context) {
	result, err := h.service.GetWormRules(c.Request.Context())
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
			"rules": result.Rules,
			"total": len(result.Rules),
		},
	})
}

// AddWormRule 添加蠕虫规则
// POST /api/v1/security/worm/rules
func (h *Handler) AddWormRule(c *gin.Context) {
	var req model.WormRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddWormRule(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "规则添加成功",
	})
}

// UpdateWormRule 更新蠕虫规则
// PUT /api/v1/security/worm/rules/:id
func (h *Handler) UpdateWormRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.WormRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateWormRule(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "规则更新成功",
	})
}

// DeleteWormRules 批量删除蠕虫规则
// DELETE /api/v1/security/worm/rules
func (h *Handler) DeleteWormRules(c *gin.Context) {
	var req struct {
		Ids []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.DeleteWormRules(c.Request.Context(), req.Ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "规则删除成功",
	})
}

// ClearWormStats 清除蠕虫统计
// POST /api/v1/security/worm/clear-stats
func (h *Handler) ClearWormStats(c *gin.Context) {
	if err := h.service.ClearWormStats(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "统计已清除",
	})
}
