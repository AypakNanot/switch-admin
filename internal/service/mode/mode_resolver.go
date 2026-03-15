package mode

import (
	"errors"
	"sync"
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
func NewModeResolver(config ModeResolverConfig) *ModeResolver {
	mode := config.InitialMode
	if mode == "" {
		mode = ModeMock // 默认为离线测试模式
	}

	return &ModeResolver{
		currentMode: mode,
		configDAO:   config.ConfigDAO,
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
