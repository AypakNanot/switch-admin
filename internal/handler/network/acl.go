package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetACLList 获取 ACL 列表
// GET /api/v1/network/acls
func (h *Handler) GetACLList(c *gin.Context) {
	result, err := h.service.GetACLList(c.Request.Context())
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
			"acls":  result.ACLs,
			"total": result.Total,
		},
	})
}

// CreateACL 创建 ACL
// POST /api/v1/network/acls
func (h *Handler) CreateACL(c *gin.Context) {
	var req model.ACLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.CreateACL(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 创建成功",
	})
}

// UpdateACL 更新 ACL
// PUT /api/v1/network/acls/:id
func (h *Handler) UpdateACL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.ACLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateACL(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 更新成功",
	})
}

// DeleteACL 删除 ACL
// DELETE /api/v1/network/acls/:id
func (h *Handler) DeleteACL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	if err := h.service.DeleteACL(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 删除成功",
	})
}

// GetACLRules 获取 ACL 规则列表
// GET /api/v1/network/acls/:id/rules
func (h *Handler) GetACLRules(c *gin.Context) {
	idStr := c.Param("id")
	aclID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	result, err := h.service.GetACLRules(c.Request.Context(), aclID)
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
			"total": result.Total,
		},
	})
}

// AddACLRule 添加 ACL 规则
// POST /api/v1/network/acls/:id/rules
func (h *Handler) AddACLRule(c *gin.Context) {
	idStr := c.Param("id")
	aclID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.ACLRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddACLRule(c.Request.Context(), aclID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 规则添加成功",
	})
}

// UpdateACLRule 更新 ACL 规则
// PUT /api/v1/network/acls/:id/rules/:ruleId
func (h *Handler) UpdateACLRule(c *gin.Context) {
	idStr := c.Param("id")
	aclID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	ruleIdStr := c.Param("ruleId")
	ruleID, err := strconv.Atoi(ruleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_RULE_ID",
		})
		return
	}

	var req model.ACLRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateACLRule(c.Request.Context(), aclID, ruleID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 规则更新成功",
	})
}

// DeleteACLRule 删除 ACL 规则
// DELETE /api/v1/network/acls/:id/rules/:ruleId
func (h *Handler) DeleteACLRule(c *gin.Context) {
	idStr := c.Param("id")
	aclID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	ruleIdStr := c.Param("ruleId")
	ruleID, err := strconv.Atoi(ruleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_RULE_ID",
		})
		return
	}

	if err := h.service.DeleteACLRule(c.Request.Context(), aclID, ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ACL 规则删除成功",
	})
}
