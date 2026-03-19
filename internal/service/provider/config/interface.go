package config

import (
	"context"

	"switch-admin/internal/model"
)

// ConfigProvider 配置模块接口
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

	// VLAN 配置
	GetVLANConfig(ctx context.Context) (*model.VLANConfigListResponse, error)
	CreateVLAN(ctx context.Context, req model.VLANCreateRequest) error
	UpdateVLAN(ctx context.Context, vlanID string, req model.VLANUpdateRequest) error
	DeleteVLAN(ctx context.Context, vlanID string) error

	// STP 配置
	GetSTPConfig(ctx context.Context) (*model.STPConfig, error)
	GetSTPStatus(ctx context.Context) (*model.STPStatus, error)
	UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error

	// PoE 配置
	GetPoEConfig(ctx context.Context) (*model.PoEConfig, error)
	UpdatePoEConfig(ctx context.Context, portID string, req model.PoEPortRequest) error

	// 风暴控制配置
	GetStormControl(ctx context.Context) (*model.StormControlConfig, error)
	UpdateStormControlGlobal(ctx context.Context, req model.StormControlRequest) error
	UpdateStormControlPort(ctx context.Context, portID string, req model.StormControlPortRequest) error

	// 流控配置
	GetFlowControl(ctx context.Context) (*model.FlowControlConfig, error)
	UpdateFlowControlGlobal(ctx context.Context, req model.FlowControlRequest) error
	UpdateFlowControlPort(ctx context.Context, portID string, req model.FlowControlPortRequest) error

	// 端口隔离配置
	GetPortIsolation(ctx context.Context) (*model.PortIsolationConfig, error)
	UpdatePortIsolation(ctx context.Context, req model.PortIsolationRequest) error
	DeletePortIsolation(ctx context.Context, groupID int) error

	// 端口镜像配置
	GetPortMirror(ctx context.Context) (*model.PortMonitorConfig, error)
	UpdatePortMirror(ctx context.Context, req model.PortMirrorRequest) error
	DeletePortMirror(ctx context.Context, sessionID int) error

	// MAC 地址表
	GetMacTable(ctx context.Context) (*model.MacTableListResponse, error)
	ClearDynamicMacEntries(ctx context.Context) error

	// ERPS 配置
	GetERPSConfig(ctx context.Context) (*model.ERPSConfig, error)
	UpdateERPSConfig(ctx context.Context, req model.ERPSRequest) error

	// 组播配置
	GetMulticastConfig(ctx context.Context) (*model.MulticastConfig, error)
	UpdateMulticastConfig(ctx context.Context, req model.MulticastRequest) error

	// 资源使用情况
	GetResource(ctx context.Context) (*model.ResourceUsage, error)

	// 堆叠配置
	GetStackConfig(ctx context.Context) (*model.StackConfig, error)
	UpdateStackMember(ctx context.Context, req model.StackRequest) error
}
