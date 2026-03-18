package cli

import (
	"os/exec"
)

// ConfigProvider CLI 模式的 Config Provider
// 通过执行系统命令和 SSH 调用交换机 API 实现配置功能
type ConfigProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewConfigProvider 创建 CLI Config Provider
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}
