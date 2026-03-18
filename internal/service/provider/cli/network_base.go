package cli

import (
	"os/exec"
)

// NetworkProvider CLI 模式的 Network Provider
// 通过执行系统命令和 SSH 调用交换机 API 实现网络功能
type NetworkProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewNetworkProvider 创建 CLI Network Provider
func NewNetworkProvider() *NetworkProvider {
	return &NetworkProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}
