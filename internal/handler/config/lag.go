package config

import (
	"net/http"
	"strconv"

	"switch-admin/internal/model"

	"github.com/gin-gonic/gin"
)

// GetLinkAggregation 获取链路聚合列表
// GET /api/v1/config/link-aggregation
func (h *Handler) GetLinkAggregation(c *gin.Context) {
	result, err := h.service.GetLinkAggregationList(c.Request.Context())
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
			"aggregations": result.Aggregations,
			"total":        result.Total,
		},
	})
}

// CreateLinkAggregation 创建链路聚合组
// POST /api/v1/config/link-aggregation
func (h *Handler) CreateLinkAggregation(c *gin.Context) {
	var req model.LinkAggregationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.CreateLinkAggregation(c.Request.Context(), req); err != nil {
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

// UpdateLinkAggregation 更新链路聚合组
// PUT /api/v1/config/link-aggregation/:id
func (h *Handler) UpdateLinkAggregation(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	var req model.LinkAggregationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.UpdateLinkAggregation(c.Request.Context(), id, req); err != nil {
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

// DeleteLinkAggregation 删除链路聚合组
// DELETE /api/v1/config/link-aggregation/:id
func (h *Handler) DeleteLinkAggregation(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_ID",
		})
		return
	}

	if err := h.service.DeleteLinkAggregation(c.Request.Context(), id); err != nil {
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
