package mode

import (
	"errors"
	"sync"

	"switch-admin/internal/service/provider"
	"switch-admin/internal/service/provider/diagnostic"
	"switch-admin/internal/service/provider/maintenance"
	"switch-admin/internal/service/provider/network"
	"switch-admin/internal/service/provider/config"
)

// RunMode 运行模式
type RunMode string

const (
	// ModeMock 离线测试模式 - 使用数据库模拟数据
	ModeMock RunMode = "mock"
	// ModeSwitch 交换机模式 - 使用真实交换机硬件
	ModeSwitch RunMode = "switch"
)

// IsValid 检查模式是否有效
func (m RunMode) IsValid() bool {
	return m == ModeMock || m == ModeSwitch
}

// String 返回模式字符串
func (m RunMode) String() string {
	return string(m)
}

// Description 返回模式描述
func (m RunMode) Description() string {
	switch m {
	case ModeMock:
		return "离线测试模式"
	case ModeSwitch:
		return "交换机模式"
	default:
		return "未知模式"
	}
}

// 错误定义
var (
	ErrInvalidMode    = errors.New("无效的运行模式")
	ErrModeSwitchBusy = errors.New("模式切换中，请稍后重试")
)

// ModeResolver 模式解析器
// 负责管理当前运行模式，并提供模式切换功能
type ModeResolver struct {
	mu sync.RWMutex

	currentMode RunMode

	// 配置持久化接口
	configDAO ConfigDAO

	// Provider 缓存
	mockDiagnosticProvider      *diagnostic.MockProvider
	cliDiagnosticProvider       *diagnostic.CLIProvider
	mockMaintenanceProvider     *maintenance.MockProvider
	cliMaintenanceProvider      *maintenance.CLIProvider
	mockNetworkProvider         *network.MockProvider
	cliNetworkProvider          *network.CLIProvider
	mockConfigProvider          *config.MockProvider
	cliConfigProvider           *config.CLIProvider
}

// ConfigDAO 配置数据访问接口（用于解耦）
type ConfigDAO interface {
	GetRunMode() (string, error)
	SetRunMode(mode string) error
}

// ModeResolverConfig 模式解析器配置
type ModeResolverConfig struct {
	InitialMode RunMode
	ConfigDAO   ConfigDAO
}

// NewModeResolver 创建模式解析器
func NewModeResolver(cfg ModeResolverConfig) *ModeResolver {
	mode := cfg.InitialMode
	if mode == "" {
		mode = ModeMock // 默认为离线测试模式
	}

	return &ModeResolver{
		currentMode:            mode,
		configDAO:              cfg.ConfigDAO,
		mockDiagnosticProvider:      diagnostic.NewMockProvider(),
		cliDiagnosticProvider:       diagnostic.NewCLIProvider(),
		mockMaintenanceProvider:     maintenance.NewMockProvider(),
		cliMaintenanceProvider:      maintenance.NewCLIProvider(),
		mockNetworkProvider:         network.NewMockProvider(),
		cliNetworkProvider:          network.NewCLIProvider(),
		mockConfigProvider:          config.NewMockProvider(),
		cliConfigProvider:           config.NewCLIProvider(),
	}
}

// GetCurrentMode 获取当前运行模式
func (r *ModeResolver) GetCurrentMode() RunMode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currentMode
}

// GetCurrentModeString 获取当前运行模式字符串
func (r *ModeResolver) GetCurrentModeString() string {
	return r.GetCurrentMode().String()
}

// GetModeDescription 获取当前模式描述
func (r *ModeResolver) GetModeDescription() string {
	return r.GetCurrentMode().Description()
}

// SwitchMode 切换运行模式
// 此操作会持久化到数据库
func (r *ModeResolver) SwitchMode(newMode RunMode) error {
	if !newMode.IsValid() {
		return ErrInvalidMode
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	oldMode := r.currentMode

	// 如果模式相同，无需切换
	if oldMode == newMode {
		return nil
	}

	// 持久化到数据库
	if r.configDAO != nil {
		if err := r.configDAO.SetRunMode(newMode.String()); err != nil {
			return err
		}
	}

	// 更新内存中的模式
	r.currentMode = newMode

	// TODO: 触发模式切换事件（清理旧资源、初始化新组件）

	return nil
}

// IsMockMode 是否为离线测试模式
func (r *ModeResolver) IsMockMode() bool {
	return r.GetCurrentMode() == ModeMock
}

// IsSwitchMode 是否为交换机模式
func (r *ModeResolver) IsSwitchMode() bool {
	return r.GetCurrentMode() == ModeSwitch
}

// LoadFromDatabase 从数据库加载模式
func (r *ModeResolver) LoadFromDatabase() error {
	if r.configDAO == nil {
		return nil
	}

	modeStr, err := r.configDAO.GetRunMode()
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	mode := RunMode(modeStr)
	if mode.IsValid() {
		r.currentMode = mode
	}

	return nil
}

// GetDiagnosticProvider 根据当前模式返回对应的 Diagnostic Provider
func (r *ModeResolver) GetDiagnosticProvider() provider.DiagnosticProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var diagnosticProvider provider.DiagnosticProvider
	if r.currentMode == ModeSwitch {
		diagnosticProvider = r.cliDiagnosticProvider
	} else {
		diagnosticProvider = r.mockDiagnosticProvider
	}

	return diagnosticProvider
}

// GetMaintenanceProvider 根据当前模式返回对应的 Maintenance Provider
func (r *ModeResolver) GetMaintenanceProvider() provider.MaintenanceProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.currentMode == ModeSwitch {
		return r.cliMaintenanceProvider
	}
	return r.mockMaintenanceProvider
}

// GetNetworkProvider 根据当前模式返回对应的 Network Provider
func (r *ModeResolver) GetNetworkProvider() provider.NetworkProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.currentMode == ModeSwitch {
		return r.cliNetworkProvider
	}
	return r.mockNetworkProvider
}

// GetConfigProvider 根据当前模式返回对应的 Config Provider
func (r *ModeResolver) GetConfigProvider() provider.ConfigProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.currentMode == ModeSwitch {
		return r.cliConfigProvider
	}
	return r.mockConfigProvider
}
