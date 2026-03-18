package mock

import (
	"math/rand"
	"time"
)

// MaintenanceProvider Mock 模式的 Maintenance Provider
// 用于离线测试模式，生成模拟数据
type MaintenanceProvider struct{}

// NewMaintenanceProvider 创建 Mock Maintenance Provider
func NewMaintenanceProvider() *MaintenanceProvider {
	rand.Seed(time.Now().UnixNano())
	return &MaintenanceProvider{}
}
