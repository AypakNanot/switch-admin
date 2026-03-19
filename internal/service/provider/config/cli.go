package config

import (
	"os/exec"
)

// CLIProvider CLI 模式的 Config Provider
// 通过执行系统命令和 SSH 调用交换机 API 实现配置功能
type CLIProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewCLIProvider 创建 CLI Config Provider
func NewCLIProvider() *CLIProvider {
	return &CLIProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}
