package dao

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	"switch-admin/internal/model"
)

var (
	ErrConfigNotFound = errors.New("配置不存在")
	dbPath            = "data/admin.db"
)

// getDB 获取数据库连接
func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil
	}
	return db
}

// ConfigDAO 配置数据访问对象
type ConfigDAO struct{}

// NewConfigDAO 创建配置 DAO
func NewConfigDAO() *ConfigDAO {
	return &ConfigDAO{}
}

// GetRunMode 获取运行模式
func (d *ConfigDAO) GetRunMode() (string, error) {
数据库 := getDB()
	defer 数据库.Close()

	var configValue string
	err := 数据库.QueryRow("SELECT config_value FROM sys_config WHERE config_key = ?", "run_mode").Scan(&configValue)
	if err != nil {
		return "mock", nil // 默认模式
	}

	return configValue, nil
}

// SetRunMode 设置运行模式
func (d *ConfigDAO) SetRunMode(mode string) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	// 检查配置是否存在
	var exists int
	err := 数据库.QueryRow("SELECT COUNT(*) FROM sys_config WHERE config_key = ?", "run_mode").Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		// 更新
		_, err = 数据库.Exec("UPDATE sys_config SET config_value = ?, updated_at = ? WHERE config_key = ?",
			mode, now, "run_mode")
		return err
	} else {
		// 插入
		_, err = 数据库.Exec("INSERT INTO sys_config (config_key, config_value, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			"run_mode", mode, "运行模式：mock=离线测试，switch=交换机", now, now)
		return err
	}
}

// GetConfig 获取配置项
func (d *ConfigDAO) GetConfig(key string) (string, error) {
	数据库 := getDB()
	defer 数据库.Close()

	var configValue string
	err := 数据库.QueryRow("SELECT config_value FROM sys_config WHERE config_key = ?", key).Scan(&configValue)
	if err != nil {
		return "", err
	}

	return configValue, nil
}

// SetConfig 设置配置项
func (d *ConfigDAO) SetConfig(key, value, description string) error {
	数据库 := getDB()
	defer 数据库.Close()

	exists, err := d.configExists(key)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if exists {
		_, err = 数据库.Exec("UPDATE sys_config SET config_value = ?, updated_at = ? WHERE config_key = ?",
			value, now, key)
		return err
	} else {
		_, err = 数据库.Exec("INSERT INTO sys_config (config_key, config_value, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			key, value, description, now, now)
		return err
	}
}

// configExists 检查配置是否存在
func (d *ConfigDAO) configExists(key string) (bool, error) {
	数据库 := getDB()
	defer 数据库.Close()

	var count int
	err := 数据库.QueryRow("SELECT COUNT(*) FROM sys_config WHERE config_key = ?", key).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAdapterConfigs 获取指定功能的所有适配器配置
func (d *ConfigDAO) GetAdapterConfigs(functionName string) ([]model.AdapterConfig, error) {
	数据库 := getDB()
	defer 数据库.Close()

	rows, err := 数据库.Query("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config WHERE function_name = ? ORDER BY priority DESC", functionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.AdapterConfig
	for rows.Next() {
		var config model.AdapterConfig
		err := rows.Scan(&config.ID, &config.FunctionName, &config.AdapterType, &config.Priority, &config.Enabled, &config.Config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// GetEnabledAdapterConfig 获取指定功能的已启用适配器配置（优先级最高）
func (d *ConfigDAO) GetEnabledAdapterConfig(functionName string) (*model.AdapterConfig, error) {
	数据库 := getDB()
	defer 数据库.Close()

	var config model.AdapterConfig
	err := 数据库.QueryRow("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config WHERE function_name = ? AND enabled = 1 ORDER BY priority DESC LIMIT 1", functionName).Scan(
		&config.ID, &config.FunctionName, &config.AdapterType, &config.Priority, &config.Enabled, &config.Config)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	return &config, nil
}

// UpdateAdapterConfig 更新适配器配置
func (d *ConfigDAO) UpdateAdapterConfig(id int64, enabled bool, priority int, config string) error {
	数据库 := getDB()
	defer 数据库.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err := 数据库.Exec("UPDATE adapter_config SET enabled = ?, priority = ?, config = ?, updated_at = ? WHERE id = ?",
		enabledInt, priority, config, now, id)
	return err
}

// GetAllAdapterConfigs 获取所有适配器配置
func (d *ConfigDAO) GetAllAdapterConfigs() ([]model.AdapterConfig, error) {
	数据库 := getDB()
	defer 数据库.Close()

	rows, err := 数据库.Query("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config ORDER BY function_name ASC, priority DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.AdapterConfig
	for rows.Next() {
		var config model.AdapterConfig
		err := rows.Scan(&config.ID, &config.FunctionName, &config.AdapterType, &config.Priority, &config.Enabled, &config.Config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
