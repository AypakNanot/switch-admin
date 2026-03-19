package config

import (
	"context"

	"switch-admin/internal/model"
)

// GetVLANConfig 获取 VLAN 配置 (Mock)
func (p *MockProvider) GetVLANConfig(ctx context.Context) (*model.VLANConfigListResponse, error) {
	return &model.VLANConfigListResponse{}, nil
}

// CreateVLAN 创建 VLAN (Mock)
func (p *MockProvider) CreateVLAN(ctx context.Context, req model.VLANCreateRequest) error {
	return nil
}

// UpdateVLAN 更新 VLAN 配置 (Mock)
func (p *MockProvider) UpdateVLAN(ctx context.Context, vlanID string, req model.VLANUpdateRequest) error {
	return nil
}

// DeleteVLAN 删除 VLAN (Mock)
func (p *MockProvider) DeleteVLAN(ctx context.Context, vlanID string) error {
	return nil
}

// GetSTPConfig 获取 STP 配置 (Mock)
func (p *MockProvider) GetSTPConfig(ctx context.Context) (*model.STPConfig, error) {
	return &model.STPConfig{}, nil
}

// GetSTPStatus 获取 STP 状态 (Mock)
func (p *MockProvider) GetSTPStatus(ctx context.Context) (*model.STPStatus, error) {
	return &model.STPStatus{}, nil
}

// UpdateSTPConfig 更新 STP 配置 (Mock)
func (p *MockProvider) UpdateSTPConfig(ctx context.Context, req model.STPConfigRequest) error {
	return nil
}

// GetPoEConfig 获取 PoE 配置 (Mock)
func (p *MockProvider) GetPoEConfig(ctx context.Context) (*model.PoEConfig, error) {
	return &model.PoEConfig{}, nil
}

// UpdatePoEConfig 更新 PoE 配置 (Mock)
func (p *MockProvider) UpdatePoEConfig(ctx context.Context, portID string, req model.PoEPortRequest) error {
	return nil
}

// GetStormControl 获取风暴控制配置 (Mock)
func (p *MockProvider) GetStormControl(ctx context.Context) (*model.StormControlConfig, error) {
	return &model.StormControlConfig{
		Enabled:   true,
		Mode:      "kbps",
		StormType: "broadcast",
		MaxRate:   10000,
		Interval:  1,
		Action:    "drop",
		Ports: []model.StormControlPort{
			{PortID: "GE1/0/1", Enabled: true, StormType: "broadcast", MaxRate: 5000, CurrentRate: 1200, Status: "normal"},
			{PortID: "GE1/0/2", Enabled: true, StormType: "multicast", MaxRate: 3000, CurrentRate: 2800, Status: "warning"},
			{PortID: "GE1/0/3", Enabled: false, StormType: "unknown-unicast", MaxRate: 2000, CurrentRate: 0, Status: "disabled"},
			{PortID: "GE1/0/4", Enabled: true, StormType: "broadcast", MaxRate: 5000, CurrentRate: 150, Status: "normal"},
		},
	}, nil
}

// UpdateStormControlGlobal 更新全局风暴控制配置 (Mock)
func (p *MockProvider) UpdateStormControlGlobal(ctx context.Context, req model.StormControlRequest) error {
	return nil
}

// UpdateStormControlPort 更新端口风暴控制配置 (Mock)
func (p *MockProvider) UpdateStormControlPort(ctx context.Context, portID string, req model.StormControlPortRequest) error {
	return nil
}

// GetFlowControl 获取流控配置 (Mock)
func (p *MockProvider) GetFlowControl(ctx context.Context) (*model.FlowControlConfig, error) {
	return &model.FlowControlConfig{
		Enabled:      true,
		Mode:         "auto",
		Backpressure: false,
		PauseType:    "symmetric",
		Ports: []model.FlowControlPort{
			{PortID: "GE1/0/1", Enabled: true, Status: "up", Negotiation: "Full/On", PauseDirection: "both"},
			{PortID: "GE1/0/2", Enabled: false, Status: "down", Negotiation: "-", PauseDirection: "none"},
			{PortID: "GE1/0/3", Enabled: true, Status: "up", Negotiation: "Full/Off", PauseDirection: "none"},
			{PortID: "GE1/0/4", Enabled: true, Status: "up", Negotiation: "Half/Backpressure", PauseDirection: "backpressure"},
		},
	}, nil
}

// UpdateFlowControlGlobal 更新全局流控配置 (Mock)
func (p *MockProvider) UpdateFlowControlGlobal(ctx context.Context, req model.FlowControlRequest) error {
	return nil
}

// UpdateFlowControlPort 更新端口流控配置 (Mock)
func (p *MockProvider) UpdateFlowControlPort(ctx context.Context, portID string, req model.FlowControlPortRequest) error {
	return nil
}

// GetPortIsolation 获取端口隔离配置 (Mock)
func (p *MockProvider) GetPortIsolation(ctx context.Context) (*model.PortIsolationConfig, error) {
	return &model.PortIsolationConfig{
		Enabled: true,
		IsolationGroups: []model.PortIsolationGroup{
			{GroupID: 1, Name: "隔离组 1", Ports: []string{"GE1/0/1", "GE1/0/2"}, IsolationMode: "l2"},
			{GroupID: 2, Name: "隔离组 2", Ports: []string{"GE1/0/3", "GE1/0/4"}, IsolationMode: "all"},
		},
	}, nil
}

// UpdatePortIsolation 更新端口隔离配置 (Mock)
func (p *MockProvider) UpdatePortIsolation(ctx context.Context, req model.PortIsolationRequest) error {
	return nil
}

// DeletePortIsolation 删除端口隔离配置 (Mock)
func (p *MockProvider) DeletePortIsolation(ctx context.Context, groupID int) error {
	return nil
}

// GetPortMirror 获取端口镜像配置 (Mock)
func (p *MockProvider) GetPortMirror(ctx context.Context) (*model.PortMonitorConfig, error) {
	return &model.PortMonitorConfig{
		Sessions: []model.PortMirrorSession{
			{SessionID: 1, Name: "镜像会话 1", MonitorPort: "GE1/0/24", SourcePorts: []string{"GE1/0/1", "GE1/0/2"}, Direction: "both", Enabled: true},
			{SessionID: 2, Name: "镜像会话 2", MonitorPort: "GE1/0/23", SourcePorts: []string{"GE1/0/3"}, Direction: "ingress", Enabled: false},
		},
	}, nil
}

// UpdatePortMirror 更新端口镜像配置 (Mock)
func (p *MockProvider) UpdatePortMirror(ctx context.Context, req model.PortMirrorRequest) error {
	return nil
}

// DeletePortMirror 删除端口镜像配置 (Mock)
func (p *MockProvider) DeletePortMirror(ctx context.Context, sessionID int) error {
	return nil
}

// GetMacTable 获取 MAC 地址表 (Mock)
func (p *MockProvider) GetMacTable(ctx context.Context) (*model.MacTableListResponse, error) {
	return &model.MacTableListResponse{
		Entries: []model.MacTableEntry{
			{VLANID: 1, MACAddress: "00:1A:2B:3C:4D:01", PortID: "GE1/0/1", Type: "dynamic", AgingTime: 300},
			{VLANID: 1, MACAddress: "00:1A:2B:3C:4D:02", PortID: "GE1/0/2", Type: "dynamic", AgingTime: 300},
			{VLANID: 10, MACAddress: "00:1A:2B:3C:4D:03", PortID: "GE1/0/3", Type: "static", AgingTime: 0},
			{VLANID: 20, MACAddress: "00:1A:2B:3C:4D:04", PortID: "GE1/0/4", Type: "dynamic", AgingTime: 300},
		},
		Total: 4,
	}, nil
}

// ClearDynamicMacEntries 清除动态 MAC 表项 (Mock)
func (p *MockProvider) ClearDynamicMacEntries(ctx context.Context) error {
	return nil
}

// GetERPSConfig 获取 ERPS 配置 (Mock)
func (p *MockProvider) GetERPSConfig(ctx context.Context) (*model.ERPSConfig, error) {
	return &model.ERPSConfig{
		Enabled:        true,
		RingID:         1,
		ControlVLAN:    4000,
		DataVLANs:      []int{10, 20, 30},
		Role:           "auto",
		WTR:            5,
		RingStatus:     "normal",
		ActiveTopology: "clockwise",
	}, nil
}

// UpdateERPSConfig 更新 ERPS 配置 (Mock)
func (p *MockProvider) UpdateERPSConfig(ctx context.Context, req model.ERPSRequest) error {
	return nil
}

// GetMulticastConfig 获取组播配置 (Mock)
func (p *MockProvider) GetMulticastConfig(ctx context.Context) (*model.MulticastConfig, error) {
	return &model.MulticastConfig{
		Enabled:       true,
		Mode:          "igmp-snoop",
		RouterPorts:   []string{"GE1/0/24"},
		HostPorts:     []string{"GE1/0/1", "GE1/0/2", "GE1/0/3"},
		FastLeave:     true,
	}, nil
}

// UpdateMulticastConfig 更新组播配置 (Mock)
func (p *MockProvider) UpdateMulticastConfig(ctx context.Context, req model.MulticastRequest) error {
	return nil
}

// GetResource 获取资源使用情况 (Mock)
func (p *MockProvider) GetResource(ctx context.Context) (*model.ResourceUsage, error) {
	return &model.ResourceUsage{
		CPUUsage:    15,
		MemoryUsage: 42,
		Temperature: 45,
		FanStatus:   "normal",
		PowerStatus: "normal",
		Uptime:      "15d 8h 32m",
		FlashUsage:  35,
		DRAMSize:    512,
		FlashSize:   256,
	}, nil
}

// GetStackConfig 获取堆叠配置 (Mock)
func (p *MockProvider) GetStackConfig(ctx context.Context) (*model.StackConfig, error) {
	return &model.StackConfig{
		Enabled:     true,
		MasterID:    1,
		MemberCount: 3,
		Members: []model.StackMember{
			{MemberID: 1, MACAddress: "00:1A:2B:3C:4D:01", Priority: 10, Role: "master", Status: "active"},
			{MemberID: 2, MACAddress: "00:1A:2B:3C:4D:02", Priority: 5, Role: "slave", Status: "active"},
			{MemberID: 3, MACAddress: "00:1A:2B:3C:4D:03", Priority: 5, Role: "slave", Status: "active"},
		},
		Topology: "ring",
	}, nil
}

// UpdateStackMember 更新堆叠成员配置 (Mock)
func (p *MockProvider) UpdateStackMember(ctx context.Context, req model.StackRequest) error {
	return nil
}
