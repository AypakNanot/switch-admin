package network

import (
	"context"

	"switch-admin/internal/model"
)

// NetworkProvider 网络模块接口
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
