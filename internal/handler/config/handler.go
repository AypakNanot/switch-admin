package config

import (
	"switch-admin/internal/service"
)

// Handler 配置模块处理器
type Handler struct {
	service *service.ConfigService
}

// New 创建配置模块处理器
func New() *Handler {
	return &Handler{
		service: service.GetConfigService(),
	}
}
