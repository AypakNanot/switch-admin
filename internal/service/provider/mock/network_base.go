package mock

import (
	"math/rand"
	"time"
)

// NetworkProvider Mock 模式的 Network Provider
// 用于离线测试模式，生成模拟数据
type NetworkProvider struct{}

// NewNetworkProvider 创建 Mock Network Provider
func NewNetworkProvider() *NetworkProvider {
	rand.Seed(time.Now().UnixNano())
	return &NetworkProvider{}
}
