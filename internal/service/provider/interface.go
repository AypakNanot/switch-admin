package provider

import (
	"context"

	"switch-admin/internal/model"
)

// PingProvider Ping 诊断接口
// 所有 Ping 操作必须通过此接口，实现模式隔离
type PingProvider interface {
	// ExecutePing 执行 Ping 诊断
	// - ctx: 上下文，支持取消操作
	// - req: Ping 请求参数（目标 IP、VRF、Count、Timeout 等）
	// - 返回：Ping 任务响应（包含统计信息）和错误
	ExecutePing(ctx context.Context, req model.PingRequest) (*model.PingTaskResponse, error)
}

// TracerouteProvider Traceroute 诊断接口
type TracerouteProvider interface {
	// ExecuteTraceroute 执行 Traceroute 诊断
	ExecuteTraceroute(ctx context.Context, req model.TracerouteRequest) (*model.TracerouteResponse, error)
}

// CableTestProvider 电缆检测接口
type CableTestProvider interface {
	// ExecuteCableTest 执行电缆检测
	ExecuteCableTest(ctx context.Context, req model.CableTestRequest) (*model.CableTestResult, error)
}

// MaintenanceProvider 维护模块接口
// 所有维护操作必须通过此接口，实现模式隔离
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

// NetworkProvider 网络模块接口
// 所有网络操作必须通过此接口，实现模式隔离
type NetworkProvider interface {
	// VLAN 管理
	GetVLANList(ctx context.Context) (*model.VLANListResponse, error)
	CreateVLAN(ctx context.Context, req model.VLANRequest) error
	UpdateVLAN(ctx context.Context, id int, req model.VLANRequest) error
	DeleteVLAN(ctx context.Context, id int) error
	DeleteVLANs(ctx context.Context, ids []int) error
	AddVLANPort(ctx context.Context, vlanID int, port string, mode string) error
	RemoveVLANPort(ctx context.Context, vlanID int, port string) error

	// 端口管理
	GetPortList(ctx context.Context) (*model.NetworkPortListResponse, error)
	GetPortDetail(ctx context.Context, portName string) (*model.PortDetail, error)
	UpdatePort(ctx context.Context, portName string, req model.PortUpdateRequest) error
	ResetPort(ctx context.Context, portName string) error
	RestartPort(ctx context.Context, portName string) error

	// 链路聚合管理
	GetLAGList(ctx context.Context) (*model.LAGListResponse, error)
	CreateLAG(ctx context.Context, req model.LAGRequest) error
	UpdateLAG(ctx context.Context, id int, req model.LAGRequest) error
	DeleteLAG(ctx context.Context, id int) error
	AddLAGPort(ctx context.Context, lagID int, port string) error
	RemoveLAGPort(ctx context.Context, lagID int, port string) error

	// STP 管理
	GetSTPConfig(ctx context.Context) (*model.STPConfig, error)
	UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error
	GetSTPStatus(ctx context.Context) (*model.STPStatus, error)

	// ACL 管理
	GetACLList(ctx context.Context) (*model.ACLListResponse, error)
	CreateACL(ctx context.Context, req model.ACLRequest) error
	UpdateACL(ctx context.Context, id int, req model.ACLRequest) error
	DeleteACL(ctx context.Context, id int) error
	GetACLRules(ctx context.Context, aclID int) (*model.ACLRuleListResponse, error)
	AddACLRule(ctx context.Context, aclID int, req model.ACLRuleRequest) error
	UpdateACLRule(ctx context.Context, aclID int, ruleID int, req model.ACLRuleRequest) error
	DeleteACLRule(ctx context.Context, aclID int, ruleID int) error
}

// ConfigProvider 配置模块接口
// 所有配置操作必须通过此接口，实现模式隔离
type ConfigProvider interface {
	// 端口配置
	GetPortList(ctx context.Context) (*model.PortConfigListResponse, error)
	GetPortDetail(ctx context.Context, portID string) (*model.PortConfig, error)
	UpdatePort(ctx context.Context, portID string, req model.PortConfigRequest) error

	// 链路聚合配置
	GetLinkAggregationList(ctx context.Context) (*model.LinkAggregationListResponse, error)
	CreateLinkAggregation(ctx context.Context, req model.LinkAggregationRequest) error
	UpdateLinkAggregation(ctx context.Context, id int, req model.LinkAggregationRequest) error
	DeleteLinkAggregation(ctx context.Context, id int) error
}
