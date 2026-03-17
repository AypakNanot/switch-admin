package dao

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"switch-admin/internal/model"
)

var (
	ErrConfigNotFound = errors.New("配置不存在")
	dbPath            = "data/admin.db"
)

// openDB 打开数据库连接（每次操作创建新连接，避免并发问题）
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	// 自动初始化 sys_config 表
	if err := initConfigTable(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// initConfigTable 初始化 sys_config 表（如果不存在）
func initConfigTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sys_config (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			config_key VARCHAR(64) NOT NULL UNIQUE,
			config_value TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// ConfigDAO 配置数据访问对象
type ConfigDAO struct {
	mu sync.Mutex // 数据库访问互斥锁
}

// NewConfigDAO 创建配置 DAO
func NewConfigDAO() *ConfigDAO {
	return &ConfigDAO{}
}

// GetRunMode 获取运行模式
func (d *ConfigDAO) GetRunMode() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return "mock", nil
	}
	defer db.Close()

	var configValue string
	err = db.QueryRow("SELECT config_value FROM sys_config WHERE config_key = ?", "run_mode").Scan(&configValue)
	if err != nil {
		return "mock", nil // 默认模式
	}

	return configValue, nil
}

// SetRunMode 设置运行模式
func (d *ConfigDAO) SetRunMode(mode string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	// 检查配置是否存在
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM sys_config WHERE config_key = ?", "run_mode").Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		// 更新
		_, err = db.Exec("UPDATE sys_config SET config_value = ?, updated_at = ? WHERE config_key = ?",
			mode, now, "run_mode")
		return err
	} else {
		// 插入
		_, err = db.Exec("INSERT INTO sys_config (config_key, config_value, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			"run_mode", mode, "运行模式：mock=离线测试，switch=交换机", now, now)
		return err
	}
}

// GetConfig 获取配置项
func (d *ConfigDAO) GetConfig(key string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var configValue string
	err = db.QueryRow("SELECT config_value FROM sys_config WHERE config_key = ?", key).Scan(&configValue)
	if err != nil {
		return "", err
	}

	return configValue, nil
}

// SetConfig 设置配置项
func (d *ConfigDAO) SetConfig(key, value, description string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	exists, err := d.configExists(key)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if exists {
		_, err = db.Exec("UPDATE sys_config SET config_value = ?, updated_at = ? WHERE config_key = ?",
			value, now, key)
		return err
	} else {
		_, err = db.Exec("INSERT INTO sys_config (config_key, config_value, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			key, value, description, now, now)
		return err
	}
}

// configExists 检查配置是否存在
func (d *ConfigDAO) configExists(key string) (bool, error) {
	db, err := openDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sys_config WHERE config_key = ?", key).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAdapterConfigs 获取指定功能的所有适配器配置
func (d *ConfigDAO) GetAdapterConfigs(functionName string) ([]model.AdapterConfig, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config WHERE function_name = ? ORDER BY priority DESC", functionName)
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
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var config model.AdapterConfig
	err = db.QueryRow("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config WHERE function_name = ? AND enabled = 1 ORDER BY priority DESC LIMIT 1", functionName).Scan(
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
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err = db.Exec("UPDATE adapter_config SET enabled = ?, priority = ?, config = ?, updated_at = ? WHERE id = ?",
		enabledInt, priority, config, now, id)
	return err
}

// GetAllAdapterConfigs 获取所有适配器配置
func (d *ConfigDAO) GetAllAdapterConfigs() ([]model.AdapterConfig, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, function_name, adapter_type, priority, enabled, config FROM adapter_config ORDER BY function_name ASC, priority DESC")
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
