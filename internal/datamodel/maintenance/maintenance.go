package maintenance

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// 页面内容函数导出 - 保持与 main.go 中的调用兼容

// GetRebootSaveContent 重启/保存配置页面
func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
	return getRebootSaveContent(ctx)
}

// GetUsersContent 用户管理页面
func GetUsersContent(ctx *context.Context) (types.Panel, error) {
	return getUsersContent(ctx)
}

// GetMaintenanceSystemConfigContent 系统配置页面
func GetMaintenanceSystemConfigContent(ctx *context.Context) (types.Panel, error) {
	return getSystemConfigContent(ctx)
}

// GetLoadConfigContent 加载配置页面
func GetLoadConfigContent(ctx *context.Context) (types.Panel, error) {
	return getLoadConfigContent(ctx)
}

// GetFilesContent 文件管理页面
func GetFilesContent(ctx *context.Context) (types.Panel, error) {
	return getFilesContent(ctx)
}

// GetLogsContent 日志管理页面
func GetLogsContent(ctx *context.Context) (types.Panel, error) {
	return getLogsContent(ctx)
}

// GetSNMPContent SNMP 配置页面
func GetSNMPContent(ctx *context.Context) (types.Panel, error) {
	return getSNMPContent(ctx)
}

// GetSNMPTrapContent SNMP Trap 配置页面
func GetSNMPTrapContent(ctx *context.Context) (types.Panel, error) {
	return getSNMPTrapContent(ctx)
}

// GetWormProtectionContent 蠕虫攻击防护页面
func GetWormProtectionContent(ctx *context.Context) (types.Panel, error) {
	return getWormProtectionContent(ctx)
}

// GetDDoSProtectionContent DDoS 攻击防护页面
func GetDDoSProtectionContent(ctx *context.Context) (types.Panel, error) {
	return getDDoSProtectionContent(ctx)
}

// GetARPProtectionContent ARP 攻击防护页面
func GetARPProtectionContent(ctx *context.Context) (types.Panel, error) {
	return getARPProtectionContent(ctx)
}

// GetSessionsContent 当前会话页面
func GetSessionsContent(ctx *context.Context) (types.Panel, error) {
	return getSessionsContent(ctx)
}
