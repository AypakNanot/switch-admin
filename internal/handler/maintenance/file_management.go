package maintenance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"switch-admin/internal/model"
)

// GetFiles 获取文件列表
// GET /api/v1/files
func (h *Handler) GetFiles(c *gin.Context) {
	path := c.Query("path")

	files, err := h.service.GetFiles(c.Request.Context(), path)
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
			"files": files.Files,
			"total": files.Total,
		},
	})
}

// UploadFile 上传文件
// POST /api/v1/files/upload
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_FILE",
		})
		return
	}

	// 读取文件内容
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	defer fileContent.Close()

	content := make([]byte, file.Size)
	_, err = fileContent.Read(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	if err := h.service.UploadFile(c.Request.Context(), model.FileUploadRequest{
		Path:      "flash:/" + file.Filename,
		Content:   content,
		Overwrite: true,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
	})
}

// UploadFirmware 上传固件
// POST /api/v1/files/firmware
func (h *Handler) UploadFirmware(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_FILE",
		})
		return
	}

	_ = file
	// TODO: 保存固件到交换机

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "固件上传成功",
	})
}

// DownloadFile 下载文件
// GET /api/v1/files/download
func (h *Handler) DownloadFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "INVALID_PATH",
		})
		return
	}

	content, contentType, err := h.service.DownloadFile(c.Request.Context(), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=\""+filePath+"\"")
	c.Data(http.StatusOK, contentType, content)
}

// DeleteFiles 删除文件
// DELETE /api/v1/files
func (h *Handler) DeleteFiles(c *gin.Context) {
	var req model.FileDeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// 兼容旧版本：从查询参数获取
		filePath := c.Query("path")
		if filePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  400,
				"error": "INVALID_PATH",
			})
			return
		}
		req.Paths = []string{filePath}
	}

	if err := h.service.DeleteFiles(c.Request.Context(), req.Paths); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件删除成功",
	})
}
