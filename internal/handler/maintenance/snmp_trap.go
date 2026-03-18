package maintenance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switch-admin/internal/model"
)

// GetSNMPTrapConfig 获取 SNMP Trap 配置
// GET /api/v1/snmp/trap/config
func (h *Handler) GetSNMPTrapConfig(c *gin.Context) {
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
		"data": gin.H{
			"enabled_traps": []string{"Linkup", "Linkdown"},
			"trap_enabled":  config.TrapEnabled,
		},
	})
}

// UpdateSNMPTrapConfig 更新 SNMP Trap 配置
// PUT /api/v1/snmp/trap/config
func (h *Handler) UpdateSNMPTrapConfig(c *gin.Context) {
	var req struct {
		EnabledTraps []string `json:"enabled_traps"`
		TrapEnabled  bool     `json:"trap_enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// TODO: 调用 service 层更新 Trap 配置

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Trap 配置保存成功",
	})
}

// GetSNMPTrapHosts 获取 SNMP Trap 目标主机列表
// GET /api/v1/snmp/trap/hosts
func (h *Handler) GetSNMPTrapHosts(c *gin.Context) {
	hosts, err := h.service.GetTrapHosts(c.Request.Context())
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
			"hosts": hosts,
			"total": len(hosts),
		},
	})
}

// AddSNMPTrapHost 添加 SNMP Trap 目标主机
// POST /api/v1/snmp/trap/hosts
func (h *Handler) AddSNMPTrapHost(c *gin.Context) {
	var req model.TrapHostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddTrapHost(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Trap 目标添加成功",
	})
}

// DeleteSNMPTrapHost 删除 SNMP Trap 目标主机
// DELETE /api/v1/snmp/trap/hosts/:id
func (h *Handler) DeleteSNMPTrapHost(c *gin.Context) {
	host := c.Param("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_HOST",
		})
		return
	}

	if err := h.service.DeleteTrapHost(c.Request.Context(), host); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Trap 目标删除成功",
	})
}

// TestSNMPTrap 发送测试 Trap
// POST /api/v1/snmp/trap/hosts/:host/test
func (h *Handler) TestSNMPTrap(c *gin.Context) {
	host := c.Param("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_HOST",
		})
		return
	}

	if err := h.service.TestTrap(c.Request.Context(), host); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试 Trap 已发送",
	})
}
