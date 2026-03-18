package network

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetLAGList 获取链路聚合组列表
// GET /api/v1/network/lags
func (h *Handler) GetLAGList(c *gin.Context) {
	result, err := h.service.GetLAGList(c.Request.Context())
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
			"lags":  result.LAGs,
			"total": result.Total,
		},
	})
}

// CreateLAG 创建链路聚合组
// POST /api/v1/network/lags
func (h *Handler) CreateLAG(c *gin.Context) {
	var req model.LAGRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.CreateLAG(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组创建成功",
	})
}

// UpdateLAG 更新链路聚合组
// PUT /api/v1/network/lags/:id
func (h *Handler) UpdateLAG(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.LAGRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateLAG(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组更新成功",
	})
}

// DeleteLAG 删除链路聚合组
// DELETE /api/v1/network/lags/:id
func (h *Handler) DeleteLAG(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	if err := h.service.DeleteLAG(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "链路聚合组删除成功",
	})
}

// AddLAGPort 添加端口到聚合组
// POST /api/v1/network/lags/:id/ports
func (h *Handler) AddLAGPort(c *gin.Context) {
	idStr := c.Param("id")
	lagID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req struct {
		Port string `json:"port"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.AddLAGPort(c.Request.Context(), lagID, req.Port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口添加成功",
	})
}

// RemoveLAGPort 从聚合组移除端口
// DELETE /api/v1/network/lags/:id/ports/:port
func (h *Handler) RemoveLAGPort(c *gin.Context) {
	idStr := c.Param("id")
	lagID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	port := c.Param("port")
	if port == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PORT",
		})
		return
	}

	if err := h.service.RemoveLAGPort(c.Request.Context(), lagID, port); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "端口移除成功",
	})
}
