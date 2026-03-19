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

	// 其他配置功能（待扩展）
	// GetStormControl(ctx context.Context) (*model.StormControlConfig, error)
	// GetFlowControl(ctx context.Context) (*model.FlowControlConfig, error)
	// GetPortIsolation(ctx context.Context) (*model.PortIsolationConfig, error)
	// GetPortMonitor(ctx context.Context) (*model.PortMonitorConfig, error)
	// GetMacTable(ctx context.Context) (*model.MacTableListResponse, error)
	// GetERPSConfig(ctx context.Context) (*model.ERPSConfig, error)
	// GetPortMirror(ctx context.Context) (*model.PortMirrorConfig, error)
	// GetMulticastConfig(ctx context.Context) (*model.MulticastConfig, error)
	// GetResource(ctx context.Context) (*model.ResourceUsage, error)
	// GetStackConfig(ctx context.Context) (*model.StackConfig, error)
}
