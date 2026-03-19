package network

import (
	"os/exec"
)

// CLIProvider CLI 模式的 Network Provider
type CLIProvider struct {
	execFunc func(command string, args ...string) ([]byte, error)
}

// NewCLIProvider 创建 CLI Network Provider
func NewCLIProvider() *CLIProvider {
	return &CLIProvider{
		execFunc: func(command string, args ...string) ([]byte, error) {
			cmd := exec.Command(command, args...)
			return cmd.CombinedOutput()
		},
	}
}

// MockProvider Mock 模式的 Network Provider
type MockProvider struct{}

// NewMockProvider 创建 Mock Network Provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}
