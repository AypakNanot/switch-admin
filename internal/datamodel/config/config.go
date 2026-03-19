package config

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// 页面内容函数导出 - 保持与 main.go 中的调用兼容

// GetPortsContent 端口状态页面
func GetPortsContent(ctx *context.Context) (types.Panel, error) {
	return getPortsContent(ctx)
}

// GetLinkAggregationContent 链路聚合页面
func GetLinkAggregationContent(ctx *context.Context) (types.Panel, error) {
	return getLinkAggregationContent(ctx)
}

// GetStormControlContent 风暴控制页面
func GetStormControlContent(ctx *context.Context) (types.Panel, error) {
	return getStormControlContent(ctx)
}

// GetFlowControlContent 流量控制页面
func GetFlowControlContent(ctx *context.Context) (types.Panel, error) {
	return getFlowControlContent(ctx)
}

// GetPortIsolationContent 端口隔离页面
func GetPortIsolationContent(ctx *context.Context) (types.Panel, error) {
	return getPortIsolationContent(ctx)
}

// GetPortMonitorContent 端口监控页面
func GetPortMonitorContent(ctx *context.Context) (types.Panel, error) {
	return getPortMonitorContent(ctx)
}

// GetVLANContent VLAN 配置页面
func GetVLANContent(ctx *context.Context) (types.Panel, error) {
	return getVLANContent(ctx)
}

// GetMacTableContent MAC 地址表页面
func GetMacTableContent(ctx *context.Context) (types.Panel, error) {
	return getMacTableContent(ctx)
}

// GetSTPContent STP 配置页面
func GetSTPContent(ctx *context.Context) (types.Panel, error) {
	return getSTPContent(ctx)
}

// GetERPSContent ERPS 配置页面
func GetERPSContent(ctx *context.Context) (types.Panel, error) {
	return getERPSContent(ctx)
}

// GetPoEContent PoE 配置页面
func GetPoEContent(ctx *context.Context) (types.Panel, error) {
	return getPoEContent(ctx)
}

// GetPortMirrorContent 端口镜像页面
func GetPortMirrorContent(ctx *context.Context) (types.Panel, error) {
	return getPortMirrorContent(ctx)
}

// GetMulticastContent 组播配置页面
func GetMulticastContent(ctx *context.Context) (types.Panel, error) {
	return getMulticastContent(ctx)
}

// GetResourceContent 资源管理页面
func GetResourceContent(ctx *context.Context) (types.Panel, error) {
	return getResourceContent(ctx)
}

// GetStackContent 堆叠配置页面
func GetStackContent(ctx *context.Context) (types.Panel, error) {
	return getStackContent(ctx)
}
