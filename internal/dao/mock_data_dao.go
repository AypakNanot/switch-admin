package dao

import (
	"database/sql"
	"time"

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
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
	数据库 := getDB()
	defer 数据库.Close()

	rows, err := 数据库.Query("SELECT id, port_name, admin_status, link_status, speed, duplex, description, rx_bytes, tx_bytes, rx_packets, tx_packets, rx_errors, tx_errors, updated_at FROM mock_port ORDER BY port_name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []model.MockPort
	for rows.Next() {
		var port model.MockPort
		var updatedAt string
		err := rows.Scan(&port.ID, &port.PortName, &port.AdminStatus, &port.LinkStatus, &port.Speed, &port.Duplex, &port.Description,
			&port.RxBytes, &port.TxBytes, &port.RxPackets, &port.TxPackets, &port.RxErrors, &port.TxErrors, &updatedAt)
		if err != nil {
			return nil, err
		}
		port.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		ports = append(ports, port)
	}

	return ports, nil
}

// GetMockPortByName 根据名称获取模拟端口
func (d *MockDataDAO) GetMockPortByName(portName string) (*model.MockPort, error) {
	数据库 := getDB()
	defer 数据库.Close()

	var port model.MockPort
	var updatedAt string
	err := 数据库.QueryRow("SELECT id, port_name, admin_status, link_status, speed, duplex, description, rx_bytes, tx_bytes, rx_packets, tx_packets, rx_errors, tx_errors, updated_at FROM mock_port WHERE port_name = ?", portName).Scan(
		&port.ID, &port.PortName, &port.AdminStatus, &port.LinkStatus, &port.Speed, &port.Duplex, &port.Description,
		&port.RxBytes, &port.TxBytes, &port.RxPackets, &port.TxPackets, &port.RxErrors, &port.TxErrors, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	port.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return &port, nil
}

// UpdateMockPortAdminStatus 更新端口管理状态
func (d *MockDataDAO) UpdateMockPortAdminStatus(portName string, enabled bool) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err := 数据库.Exec("UPDATE mock_port SET admin_status = ?, updated_at = ? WHERE port_name = ?",
		enabledInt, now, portName)
	return err
}

// UpdateMockPortLinkStatus 更新端口链路状态
func (d *MockDataDAO) UpdateMockPortLinkStatus(portName string, up bool) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	upInt := 0
	if up {
		upInt = 1
	}

	_, err := 数据库.Exec("UPDATE mock_port SET link_status = ?, updated_at = ? WHERE port_name = ?",
		upInt, now, portName)
	return err
}

// ClearMockPortStats 清零单个端口统计
func (d *MockDataDAO) ClearMockPortStats(portName string) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := 数据库.Exec("UPDATE mock_port SET rx_bytes = 0, tx_bytes = 0, rx_packets = 0, tx_packets = 0, rx_errors = 0, tx_errors = 0, updated_at = ? WHERE port_name = ?",
		now, portName)
	return err
}

// ClearAllMockPortsStats 清零所有端口统计
func (d *MockDataDAO) ClearAllMockPortsStats() error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := 数据库.Exec("UPDATE mock_port SET rx_bytes = 0, tx_bytes = 0, rx_packets = 0, tx_packets = 0, rx_errors = 0, tx_errors = 0, updated_at = ?",
		now)
	return err
}

// GetMockSystemInfo 获取模拟系统信息
func (d *MockDataDAO) GetMockSystemInfo() (*model.MockSystemInfo, error) {
	数据库 := getDB()
	defer 数据库.Close()

	var info model.MockSystemInfo
	var bootTime string
	err := 数据库.QueryRow("SELECT id, model, serial_number, mac_address, software_version, hardware_version, uptime_seconds, boot_time FROM mock_system_info LIMIT 1").Scan(
		&info.ID, &info.Model, &info.SerialNumber, &info.MACAddress, &info.SoftwareVersion, &info.HardwareVersion,
		&info.UptimeSeconds, &bootTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	info.BootTime, _ = time.Parse("2006-01-02 15:04:05", bootTime)
	return &info, nil
}

// UpdateMockSystemInfoUptime 更新系统运行时间
func (d *MockDataDAO) UpdateMockSystemInfoUptime(uptimeSeconds int64) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := 数据库.Exec("UPDATE mock_system_info SET uptime_seconds = ?, updated_at = ? WHERE id = ?",
		uptimeSeconds, now, 1)
	return err
}
