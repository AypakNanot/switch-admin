package network

import (
	"switch-admin/internal/service"
)

// Handler 网络模块主处理器
type Handler struct {
	service *service.NetworkService
}

// New 创建网络模块处理器
func New() *Handler {
	return &Handler{
		service: service.GetNetworkService(),
	}
}
