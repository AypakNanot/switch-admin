package cli

import (
	"os/exec"
)

// MaintenanceProvider CLI 模式的 Maintenance Provider
// 通过执行系统命令和 SSH 调用交换机 API 实现维护功能
type MaintenanceProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewMaintenanceProvider 创建 CLI Maintenance Provider
func NewMaintenanceProvider() *MaintenanceProvider {
	return &MaintenanceProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}
