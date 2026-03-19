package diagnostic

// CLIProvider CLI 模式的 Diagnostic Provider
type CLIProvider struct{}

// MockProvider Mock 模式的 Diagnostic Provider
type MockProvider struct{}

// NewCLIProvider 创建 CLI Diagnostic Provider
func NewCLIProvider() *CLIProvider {
	return &CLIProvider{}
}

// NewMockProvider 创建 Mock Diagnostic Provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}
