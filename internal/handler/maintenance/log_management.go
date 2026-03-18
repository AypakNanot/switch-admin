package maintenance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLogs 获取日志列表
// GET /api/v1/logs
func (h *Handler) GetLogs(c *gin.Context) {
	level := c.Query("level")

	logs, err := h.service.GetLogs(c.Request.Context())
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
			"logs":  logs.Logs,
			"total": logs.Total,
		},
	})

	_ = level
}

// ClearLogs 清除日志
// DELETE /api/v1/logs
func (h *Handler) ClearLogs(c *gin.Context) {
	var req struct {
		Levels []string `json:"levels"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Levels = nil
	}

	if err := h.service.ClearLogs(c.Request.Context(), req.Levels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "日志已清除",
	})
}
