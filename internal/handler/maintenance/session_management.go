package maintenance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetSessions 获取会话列表
// GET /api/v1/sessions
func (h *Handler) GetSessions(c *gin.Context) {
	sessions, err := h.service.GetSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": sessions,
	})
}

// DeleteSession 删除会话（强制踢出）
// DELETE /api/v1/sessions/:session_id
func (h *Handler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_SESSION_ID",
		})
		return
	}

	if err := h.service.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "会话已终止",
	})
}

// DeleteSessions 批量删除会话
// DELETE /api/v1/sessions
func (h *Handler) DeleteSessions(c *gin.Context) {
	var req model.SessionDeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.DeleteSessions(c.Request.Context(), req.SessionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "会话已终止",
	})
}
