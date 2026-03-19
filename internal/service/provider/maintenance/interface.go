package maintenance

import (
	"context"

	"switch-admin/internal/model"
)

// MaintenanceProvider 维护模块接口
type MaintenanceProvider interface {
	// 系统配置
	GetSystemConfig(ctx context.Context) (*model.SystemConfig, error)
	UpdateNetworkConfig(ctx context.Context, req model.NetworkConfigRequest) error
	UpdateTemperatureConfig(ctx context.Context, req model.TemperatureConfigRequest) error
	UpdateDeviceInfo(ctx context.Context, req model.DeviceInfoRequest) error
	UpdateDateTime(ctx context.Context, req model.DateTimeRequest) error

	// 重启/保存
	SaveConfig(ctx context.Context) error
	RebootSwitch(ctx context.Context, delay int) error
	FactoryReset(ctx context.Context) error

	// 用户管理
	GetUsers(ctx context.Context) (*model.UserListResponse, error)
	CreateUser(ctx context.Context, req model.UserRequest) error
	DeleteUser(ctx context.Context, username string) error
	DeleteUsers(ctx context.Context, usernames []string) error

	// 会话管理
	GetSessions(ctx context.Context) (*model.SessionListResponse, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteSessions(ctx context.Context, sessionIDs []string) error

	// 日志管理
	GetLogs(ctx context.Context) (*model.LogListResponse, error)
	ClearLogs(ctx context.Context, levels []string) error

	// 文件管理
	GetFiles(ctx context.Context, path string) (*model.FileListResponse, error)
	UploadFile(ctx context.Context, req model.FileUploadRequest) error
	DeleteFile(ctx context.Context, path string) error
	DeleteFiles(ctx context.Context, paths []string) error
	DownloadFile(ctx context.Context, path string) ([]byte, string, error)

	// SNMP 配置
	GetSNMPConfig(ctx context.Context) (*model.SNMPConfig, error)
	UpdateSNMPConfig(ctx context.Context, req model.SNMPConfigRequest) error
	GetTrapHosts(ctx context.Context) ([]model.TrapHost, error)
	AddTrapHost(ctx context.Context, req model.TrapHostRequest) error
	DeleteTrapHost(ctx context.Context, host string) error
	TestTrap(ctx context.Context, host string) error
	GetSNMPCommunities(ctx context.Context) ([]model.SNMPCommunity, error)
	AddCommunity(ctx context.Context, name, access, description string) error
	DeleteCommunity(ctx context.Context, name string) error

	// 蠕虫防护
	GetWormRules(ctx context.Context) (*model.WormRuleList, error)
	AddWormRule(ctx context.Context, req model.WormRuleRequest) error
	UpdateWormRule(ctx context.Context, id string, req model.WormRuleRequest) error
	DeleteWormRule(ctx context.Context, id string) error
	DeleteWormRules(ctx context.Context, ids []string) error
	ClearWormStats(ctx context.Context) error

	// DDoS 防护
	GetDDoSConfig(ctx context.Context) (*model.DDoSConfig, error)
	UpdateDDoSConfig(ctx context.Context, req model.DDoSConfigRequest) error

	// ARP 防护
	GetARPConfig(ctx context.Context) (*model.ARPConfig, error)
	UpdateARPConfig(ctx context.Context, req model.ARPConfigRequest) error

	// 加载配置
	GetConfigFiles(ctx context.Context) (*model.LoadConfigListResponse, error)
	LoadConfig(ctx context.Context, configFile string) error
}
