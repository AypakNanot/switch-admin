package maintenance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/service"
)

// Handler 维护模块主处理器
type Handler struct {
	service *service.MaintenanceService
}

// New 创建维护模块处理器
func New() *Handler {
	return &Handler{
		service: service.GetMaintenanceService(),
	}
}

// SaveConfig 保存配置
// POST /api/v1/system/save-config
func (h *Handler) SaveConfig(c *gin.Context) {
	if err := h.service.SaveConfig(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "配置保存成功",
	})
}

// RebootSwitch 重启交换机
// POST /api/v1/system/reboot
func (h *Handler) RebootSwitch(c *gin.Context) {
	var req struct {
		Delay int `json:"delay"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Delay = 30
	}

	if err := h.service.RebootSwitch(c.Request.Context(), req.Delay); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "重启命令已发送，设备将在 30 秒后重启",
	})
}

// FactoryReset 恢复出厂配置
// POST /api/v1/system/factory-reset
func (h *Handler) FactoryReset(c *gin.Context) {
	if err := h.service.FactoryReset(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "恢复出厂设置成功，设备将重启",
	})
}

// GetSystemConfig 获取系统配置
// GET /api/v1/system/config
func (h *Handler) GetSystemConfig(c *gin.Context) {
	config, err := h.service.GetSystemConfig(c.Request.Context())
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

// UpdateNetworkConfig 更新网络配置
// PUT /api/v1/system/network
func (h *Handler) UpdateNetworkConfig(c *gin.Context) {
	var req struct {
		IP      string `json:"ip"`
		Mask    string `json:"mask"`
		Gateway string `json:"gateway"`
		DNS     string `json:"dns"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateNetworkConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "网络配置保存成功",
	})
}

// UpdateTemperatureConfig 更新温度阈值配置
// PUT /api/v1/system/temperature
func (h *Handler) UpdateTemperatureConfig(c *gin.Context) {
	var req struct {
		Low  int `json:"low"`
		High int `json:"high"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateTemperatureConfig(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "温度配置保存成功",
	})
}

// UpdateDeviceInfo 更新设备信息
// PUT /api/v1/system/info
func (h *Handler) UpdateDeviceInfo(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Contact  string `json:"contact"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateDeviceInfo(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备信息保存成功",
	})
}

// UpdateDateTime 更新时间日期
// PUT /api/v1/system/datetime
func (h *Handler) UpdateDateTime(c *gin.Context) {
	var req struct {
		Timezone string `json:"timezone"`
		DateTime string `json:"datetime"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateDateTime(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "日期时间配置保存成功",
	})
}
