package maintenance

import (
	"os/exec"
)

// CLIProvider CLI 模式的 Maintenance Provider
type CLIProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewCLIProvider 创建 CLI Maintenance Provider
func NewCLIProvider() *CLIProvider {
	return &CLIProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}

// MockProvider Mock 模式的 Maintenance Provider
type MockProvider struct{}

// NewMockProvider 创建 Mock Maintenance Provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}
