package network

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetVLANContent VLAN 管理页面
func GetVLANContent(ctx *context.Context) (types.Panel, error) {
	return getVLANContent(ctx)
}

// GetPortContent 端口管理页面
func GetPortContent(ctx *context.Context) (types.Panel, error) {
	return getPortContent(ctx)
}

// GetLAGContent 链路聚合管理页面
func GetLAGContent(ctx *context.Context) (types.Panel, error) {
	return getLAGContent(ctx)
}

// GetSTPContent STP 管理页面
func GetSTPContent(ctx *context.Context) (types.Panel, error) {
	return getSTPContent(ctx)
}

// GetACLContent ACL 管理页面
func GetACLContent(ctx *context.Context) (types.Panel, error) {
	return getACLContent(ctx)
}
