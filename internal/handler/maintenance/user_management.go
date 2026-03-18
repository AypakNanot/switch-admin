package maintenance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetUsers 获取用户列表
// GET /api/v1/users
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
	})
}

// CreateUser 创建用户
// POST /api/v1/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req model.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.CreateUser(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户创建成功",
	})
}

// UpdateUser 更新用户
// PUT /api/v1/users/:username
func (h *Handler) UpdateUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_USERNAME",
		})
		return
	}

	var req struct {
		Password string `json:"password"`
		Role     int    `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	// TODO: 调用服务层更新用户

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户更新成功",
	})
}

// DeleteUsers 批量删除用户
// DELETE /api/v1/users
func (h *Handler) DeleteUsers(c *gin.Context) {
	var req model.UserDeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_REQUEST",
		})
		return
	}

	if err := h.service.DeleteUsers(c.Request.Context(), req.Usernames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户删除成功",
	})
}
