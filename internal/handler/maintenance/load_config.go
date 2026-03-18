package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetConfigFiles 获取配置文件列表
// GET /api/v1/config/files
func (h *Handler) GetConfigFiles(c *gin.Context) {
	result, err := h.service.GetConfigFiles(c.Request.Context())
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
			"files": result.Files,
			"total": len(result.Files),
		},
	})
}

// LoadConfig 加载配置文件
// POST /api/v1/config/load
func (h *Handler) LoadConfig(c *gin.Context) {
	var req struct {
		ConfigFile string `json:"config_file"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.LoadConfig(c.Request.Context(), req.ConfigFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "配置加载成功，部分配置可能需要重启生效",
	})
}
