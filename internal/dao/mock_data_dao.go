package dao

import (
	"time"

	"github.com/GoAdminGroup/go-admin/modules/db"
	"switch-admin/internal/model"
)

// MockDataDAO 模拟数据访问对象
type MockDataDAO struct{}

// NewMockDataDAO 创建模拟数据 DAO
func NewMockDataDAO() *MockDataDAO {
	return &MockDataDAO{}
}

// GetAllMockPorts 获取所有模拟端口
func (d *MockDataDAO) GetAllMockPorts() ([]model.MockPort, error) {
	rows, err := db.Table("mock_port").
		OrderBy("port_name", "ASC").
		All()
	if err != nil {
		return nil, err
	}

	var ports []model.MockPort
	for _, row := range rows {
		port := d.rowToMockPort(row)
		ports = append(ports, port)
	}

	return ports, nil
}

// GetMockPortByName 根据名称获取模拟端口
func (d *MockDataDAO) GetMockPortByName(portName string) (*model.MockPort, error) {
	res, err := db.Table("mock_port").
		Where("port_name", "=", portName).
		First()
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrConfigNotFound
	}

	port := d.rowToMockPort(res)
	return &port, nil
}

// UpdateMockPortAdminStatus 更新端口管理状态
func (d *MockDataDAO) UpdateMockPortAdminStatus(portName string, enabled bool) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err := db.Table("mock_port").
		Where("port_name", "=", portName).
		Update(map[string]interface{}{
			"admin_status": enabledInt,
			"updated_at":   now,
		})
	return err
}

// UpdateMockPortLinkStatus 更新端口链路状态
func (d *MockDataDAO) UpdateMockPortLinkStatus(portName string, up bool) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	upInt := 0
	if up {
		upInt = 1
	}

	_, err := db.Table("mock_port").
		Where("port_name", "=", portName).
		Update(map[string]interface{}{
			"link_status": upInt,
			"updated_at":  now,
		})
	return err
}

// ClearMockPortStats 清零单个端口统计
func (d *MockDataDAO) ClearMockPortStats(portName string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Table("mock_port").
		Where("port_name", "=", portName).
		Update(map[string]interface{}{
			"rx_bytes":   0,
			"tx_bytes":   0,
			"rx_packets": 0,
			"tx_packets": 0,
			"rx_errors":  0,
			"tx_errors":  0,
			"updated_at": now,
		})
	return err
}

// ClearAllMockPortsStats 清零所有端口统计
func (d *MockDataDAO) ClearAllMockPortsStats() error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Table("mock_port").
		Update(map[string]interface{}{
			"rx_bytes":   0,
			"tx_bytes":   0,
			"rx_packets": 0,
			"tx_packets": 0,
			"rx_errors":  0,
			"tx_errors":  0,
			"updated_at": now,
		})
	return err
}

// GetMockSystemInfo 获取模拟系统信息
func (d *MockDataDAO) GetMockSystemInfo() (*model.MockSystemInfo, error) {
	res, err := db.Table("mock_system_info").
		Take(1).
		First()
	if err != nil {
		return nil, ErrConfigNotFound
	}

	if res == nil {
		return nil, ErrConfigNotFound
	}

	return d.rowToMockSystemInfo(res), nil
}

// UpdateMockSystemInfoUptime 更新系统运行时间
func (d *MockDataDAO) UpdateMockSystemInfoUptime(uptimeSeconds int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Table("mock_system_info").
		Where("id", "=", 1).
		Update(map[string]interface{}{
			"uptime_seconds": uptimeSeconds,
			"updated_at":     now,
		})
	return err
}

// rowToMockPort 将数据库行转换为 MockPort
func (d *MockDataDAO) rowToMockPort(row map[string]interface{}) model.MockPort {
	port := model.MockPort{
		ID:          int64(row["id"].(int)),
		PortName:    getString(row["port_name"]),
		AdminStatus: row["admin_status"].(int) == 1,
		LinkStatus:  row["link_status"].(int) == 1,
		Speed:       getString(row["speed"]),
		Duplex:      getString(row["duplex"]),
		Description: getString(row["description"]),
		RxBytes:     getUint64(row["rx_bytes"]),
		TxBytes:     getUint64(row["tx_bytes"]),
		RxPackets:   getUint64(row["rx_packets"]),
		TxPackets:   getUint64(row["tx_packets"]),
		RxErrors:    getUint64(row["rx_errors"]),
		TxErrors:    getUint64(row["tx_errors"]),
	}

	if updatedAt, ok := row["updated_at"].(string); ok {
		port.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	}

	return port
}

// rowToMockSystemInfo 将数据库行转换为 MockSystemInfo
func (d *MockDataDAO) rowToMockSystemInfo(row map[string]interface{}) *model.MockSystemInfo {
	info := &model.MockSystemInfo{
		ID:              int64(row["id"].(int)),
		Model:           getString(row["model"]),
		SerialNumber:    getString(row["serial_number"]),
		MACAddress:      getString(row["mac_address"]),
		SoftwareVersion: getString(row["software_version"]),
		HardwareVersion: getString(row["hardware_version"]),
		UptimeSeconds:   int64(row["uptime_seconds"].(int)),
	}

	if bootTime, ok := row["boot_time"].(string); ok {
		info.BootTime, _ = time.Parse("2006-01-02 15:04:05", bootTime)
	}

	return info
}

// 辅助函数
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getInt(v interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func getInt64(v interface{}) int64 {
	if i, ok := v.(int64); ok {
		return i
	}
	if i, ok := v.(int); ok {
		return int64(i)
	}
	return 0
}

func getUint64(v interface{}) uint64 {
	if u, ok := v.(uint64); ok {
		return u
	}
	if i, ok := v.(int64); ok {
		return uint64(i)
	}
	if i, ok := v.(int); ok {
		return uint64(i)
	}
	return 0
}
