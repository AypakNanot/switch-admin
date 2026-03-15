package handler

import (
	"net/http"
	"switch-admin/internal/dao"

	"github.com/gin-gonic/gin"
)

// ConfigHandler 配置模块 API 处理器
type ConfigHandler struct {
	configDAO *dao.ConfigDAO
}

// NewConfigHandler 创建配置模块处理器
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		configDAO: dao.NewConfigDAO(),
	}
}

// GetPorts 获取端口列表
func (h *ConfigHandler) GetPorts(c *gin.Context) {
	// TODO: 实现从 DAO 获取端口数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"ports": []gin.H{
				{
					"port_id":      "GE1/0/1",
					"admin_status": "enable",
					"link_status":  "up",
					"speed_duplex": "1000F",
					"flow_control": "off",
					"description":  "Server-A",
					"aggregation":  "-",
				},
				{
					"port_id":      "GE1/0/2",
					"admin_status": "enable",
					"link_status":  "down",
					"speed_duplex": "auto",
					"flow_control": "off",
					"description":  "",
					"aggregation":  "-",
				},
				{
					"port_id":      "GE1/0/3",
					"admin_status": "disable",
					"link_status":  "down",
					"speed_duplex": "auto",
					"flow_control": "off",
					"description":  "",
					"aggregation":  "Ag1",
				},
				{
					"port_id":      "GE1/0/4",
					"admin_status": "enable",
					"link_status":  "up",
					"speed_duplex": "100F",
					"flow_control": "on",
					"description":  "AP-Floor2",
					"aggregation":  "-",
				},
			},
		},
	})
}

// UpdatePort 更新端口配置
func (h *ConfigHandler) UpdatePort(c *gin.Context) {
	portId := c.Param("port_id")

	var data struct {
		AdminStatus string `json:"admin_status"`
		Description string `json:"description"`
		SpeedDuplex string `json:"speed_duplex"`
		FlowControl string `json:"flow_control"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// TODO: 调用 DAO 更新端口配置
	_ = portId
	_ = data

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口配置保存成功",
	})
}

// GetLinkAggregation 获取链路聚合列表
func (h *ConfigHandler) GetLinkAggregation(c *gin.Context) {
	// TODO: 实现从 DAO 获取聚合组数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"aggregations": []gin.H{
				{
					"group_id":     1,
					"name":         "Ag1",
					"mode":         "LACP",
					"member_ports": []string{"GE1/0/1", "GE1/0/2", "GE1/0/3", "GE1/0/4"},
					"load_balance": "src-dst-ip",
					"min_active":   2,
					"status":       "normal",
				},
				{
					"group_id":     2,
					"name":         "Ag2",
					"mode":         "Static",
					"member_ports": []string{"GE1/0/5", "GE1/0/6"},
					"load_balance": "src-dst-mac",
					"min_active":   1,
					"status":       "normal",
				},
				{
					"group_id":     3,
					"name":         "Ag3",
					"mode":         "LACP",
					"member_ports": []string{"GE1/0/9", "GE1/0/10", "GE1/0/11", "GE1/0/12"},
					"load_balance": "src-dst-mac",
					"min_active":   2,
					"status":       "degraded",
				},
			},
		},
	})
}

// CreateLinkAggregation 创建链路聚合组
func (h *ConfigHandler) CreateLinkAggregation(c *gin.Context) {
	var data struct {
		GroupID     int      `json:"group_id"`
		Mode        string   `json:"mode"`
		Description string   `json:"description"`
		LoadBalance string   `json:"load_balance"`
		MemberPorts []string `json:"member_ports"`
		MinActive   int      `json:"min_active"`
		LacpTimeout string   `json:"lacp_timeout"`
		Priority    int      `json:"priority"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// TODO: 调用 DAO 创建聚合组
	_ = data

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组创建成功",
	})
}

// UpdateLinkAggregation 更新链路聚合组
func (h *ConfigHandler) UpdateLinkAggregation(c *gin.Context) {
	id := c.Param("id")

	var data struct {
		Mode        string   `json:"mode"`
		Description string   `json:"description"`
		LoadBalance string   `json:"load_balance"`
		MemberPorts []string `json:"member_ports"`
		MinActive   int      `json:"min_active"`
		LacpTimeout string   `json:"lacp_timeout"`
		Priority    int      `json:"priority"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// TODO: 调用 DAO 更新聚合组
	_ = id
	_ = data

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组更新成功",
	})
}

// DeleteLinkAggregation 删除链路聚合组
func (h *ConfigHandler) DeleteLinkAggregation(c *gin.Context) {
	id := c.Param("id")

	// TODO: 调用 DAO 删除聚合组
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组删除成功",
	})
}
