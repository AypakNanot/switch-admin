package diagnostic

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
	"switch-admin/internal/service"
)

// DiagnosticHandler 网络诊断处理器
type DiagnosticHandler struct {
	service *service.DiagnosticService
}

// NewDiagnosticHandler 创建诊断处理器
func NewDiagnosticHandler() *DiagnosticHandler {
	return &DiagnosticHandler{
		service: service.GetDiagnosticService(),
	}
}

// CreatePingTask 创建 Ping 任务
// POST /api/v1/diagnostic/ping
func (h *DiagnosticHandler) CreatePingTask(c *gin.Context) {
	var req model.PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	// 默认参数
	if req.Count <= 0 {
		req.Count = 4
	}
	if req.Timeout <= 0 {
		req.Timeout = 2
	}
	if req.Interval <= 0 {
		req.Interval = 1
	}
	if req.VrfID == "" {
		req.VrfID = "mgmt vrf"
	}

	taskID, err := h.service.CreatePingTask(req)
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
			"task_id": taskID,
			"status":  "running",
		},
	})
}

// GetPingTaskResult 获取 Ping 任务结果
// GET /api/v1/diagnostic/ping/:task_id
func (h *DiagnosticHandler) GetPingTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")

	result, ok := h.service.GetPingTaskResult(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// DeletePingTask 删除 Ping 任务
// DELETE /api/v1/diagnostic/ping/:task_id
func (h *DiagnosticHandler) DeletePingTask(c *gin.Context) {
	taskID := c.Param("task_id")
	h.service.DeletePingTask(taskID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Task deleted",
	})
}

// CreateTracerouteTask 创建 Traceroute 任务
// POST /api/v1/diagnostic/traceroute
func (h *DiagnosticHandler) CreateTracerouteTask(c *gin.Context) {
	var req model.TracerouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	// 默认参数
	if req.MaxHops <= 0 {
		req.MaxHops = 30
	}
	if req.Timeout <= 0 {
		req.Timeout = 2
	}
	if req.Probes <= 0 {
		req.Probes = 3
	}
	if req.VrfID == "" {
		req.VrfID = "mgmt vrf"
	}

	taskID, err := h.service.CreateTracerouteTask(req)
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
			"task_id": taskID,
			"status":  "running",
		},
	})
}

// GetTracerouteTaskResult 获取 Traceroute 任务结果
// GET /api/v1/diagnostic/traceroute/:task_id
func (h *DiagnosticHandler) GetTracerouteTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")

	result, ok := h.service.GetTracerouteTaskResult(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// DeleteTracerouteTask 删除 Traceroute 任务
// DELETE /api/v1/diagnostic/traceroute/:task_id
func (h *DiagnosticHandler) DeleteTracerouteTask(c *gin.Context) {
	taskID := c.Param("task_id")
	h.service.DeleteTracerouteTask(taskID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Task deleted",
	})
}

// GetDetectablePorts 获取可检测端口列表
// GET /api/v1/diagnostic/cable/ports
func (h *DiagnosticHandler) GetDetectablePorts(c *gin.Context) {
	ports := h.service.GetDetectablePorts()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"ports": ports,
		},
	})
}

// ExecuteCableTest 执行电缆检测
// POST /api/v1/diagnostic/cable
func (h *DiagnosticHandler) ExecuteCableTest(c *gin.Context) {
	var req model.CableTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	result, err := h.service.ExecuteCableTest(req)
	if err != nil {
		// 判断错误类型
		errMsg := err.Error()
		code := http.StatusBadRequest
		if errMsg == "光口不支持虚拟电缆检测" {
			c.JSON(code, gin.H{
				"code":    code,
				"error":   "PORT_NOT_ELECTRICAL",
				"message": errMsg,
				"data": gin.H{
					"port_id":   req.PortID,
					"port_type": "optical",
				},
			})
		} else if errMsg == "端口已关闭，请先在端口配置中启用该端口" {
			c.JSON(code, gin.H{
				"code":    code,
				"error":   "PORT_ADMIN_DOWN",
				"message": errMsg,
				"data": gin.H{
					"port_id":      req.PortID,
					"admin_status": "down",
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"error":   "CABLE_TEST_FAILED",
				"message": errMsg,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}
