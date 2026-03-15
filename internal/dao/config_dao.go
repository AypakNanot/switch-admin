package dao

import (
	"errors"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/db"
	"switch-admin/internal/model"
)

var (
	ErrConfigNotFound = errors.New("配置不存在")
)

// ConfigDAO 配置数据访问对象
type ConfigDAO struct{}

// NewConfigDAO 创建配置 DAO
func NewConfigDAO() *ConfigDAO {
	return &ConfigDAO{}
}

// GetRunMode 获取运行模式
func (d *ConfigDAO) GetRunMode() (string, error) {
	res, err := db.Table("sys_config").
		Where("config_key", "=", "run_mode").
		First()
	if err != nil {
		return "mock", nil // 默认模式
	}

	if res == nil {
		return "mock", nil
	}

	if configValue, ok := res["config_value"].(string); ok {
		return configValue, nil
	}

	return "mock", nil
}

// SetRunMode 设置运行模式
func (d *ConfigDAO) SetRunMode(mode string) error {
	// 检查配置是否存在
	exists, err := d.configExists("run_mode")
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if exists {
		// 更新
		_, err = db.Table("sys_config").
			Where("config_key", "=", "run_mode").
			Update(map[string]interface{}{
				"config_value": mode,
				"updated_at":   now,
			})
		return err
	} else {
		// 插入
		_, err = db.Table("sys_config").
			Insert(map[string]interface{}{
				"config_key":   "run_mode",
				"config_value": mode,
				"description":  "运行模式：mock=离线测试，switch=交换机",
				"created_at":   now,
				"updated_at":   now,
			})
		return err
	}
}

// GetConfig 获取配置项
func (d *ConfigDAO) GetConfig(key string) (string, error) {
	res, err := db.Table("sys_config").
		Where("config_key", "=", key).
		First()
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", ErrConfigNotFound
	}

	if configValue, ok := res["config_value"].(string); ok {
		return configValue, nil
	}

	return "", ErrConfigNotFound
}

// SetConfig 设置配置项
func (d *ConfigDAO) SetConfig(key, value, description string) error {
	exists, err := d.configExists(key)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if exists {
		_, err = db.Table("sys_config").
			Where("config_key", "=", key).
			Update(map[string]interface{}{
				"config_value": value,
				"updated_at":   now,
			})
		return err
	} else {
		_, err = db.Table("sys_config").
			Insert(map[string]interface{}{
				"config_key":   key,
				"config_value": value,
				"description":  description,
				"created_at":   now,
				"updated_at":   now,
			})
		return err
	}
}

// configExists 检查配置是否存在
func (d *ConfigDAO) configExists(key string) (bool, error) {
	res, err := db.Table("sys_config").
		Where("config_key", "=", key).
		First()
	if err != nil {
		return false, err
	}
	return res != nil, nil
}

// GetAdapterConfigs 获取指定功能的所有适配器配置
func (d *ConfigDAO) GetAdapterConfigs(functionName string) ([]model.AdapterConfig, error) {
	rows, err := db.Table("adapter_config").
		Where("function_name", "=", functionName).
		OrderBy("priority", "DESC").
		All()
	if err != nil {
		return nil, err
	}

	var configs []model.AdapterConfig
	for _, row := range rows {
		config := model.AdapterConfig{
			ID:           int64(row["id"].(int)),
			FunctionName: row["function_name"].(string),
			AdapterType:  row["adapter_type"].(string),
			Priority:     row["priority"].(int),
			Enabled:      row["enabled"].(int) == 1,
			Config:       row["config"].(string),
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// GetEnabledAdapterConfig 获取指定功能的已启用适配器配置（优先级最高）
func (d *ConfigDAO) GetEnabledAdapterConfig(functionName string) (*model.AdapterConfig, error) {
	res, err := db.Table("adapter_config").
		Where("function_name", "=", functionName).
		Where("enabled", "=", 1).
		OrderBy("priority", "DESC").
		First()
	if err != nil {
		return nil, ErrConfigNotFound
	}

	if res == nil {
		return nil, ErrConfigNotFound
	}

	return &model.AdapterConfig{
		ID:           int64(res["id"].(int)),
		FunctionName: res["function_name"].(string),
		AdapterType:  res["adapter_type"].(string),
		Priority:     res["priority"].(int),
		Enabled:      res["enabled"].(int) == 1,
		Config:       res["config"].(string),
	}, nil
}

// UpdateAdapterConfig 更新适配器配置
func (d *ConfigDAO) UpdateAdapterConfig(id int64, enabled bool, priority int, config string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	_, err := db.Table("adapter_config").
		Where("id", "=", id).
		Update(map[string]interface{}{
			"enabled":    enabledInt,
			"priority":   priority,
			"config":     config,
			"updated_at": now,
		})
	return err
}
