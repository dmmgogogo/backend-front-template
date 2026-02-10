package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type SystemConfig struct {
	ID          int64  `orm:"column(id);pk;auto" json:"id"`
	ConfigKey   string `orm:"column(config_key);size(100);unique" json:"config_key"`
	ConfigValue string `orm:"column(config_value);size(255)" json:"config_value"`
	ConfigDesc  string `orm:"column(config_desc);size(255)" json:"config_desc"`
	CreatedTime int64  `orm:"column(created_time)" json:"created_time"`
	UpdatedTime int64  `orm:"column(updated_time)" json:"updated_time"`
}

func (sc *SystemConfig) TableName() string {
	return "app_system_config"
}

func init() {
	orm.RegisterModel(new(SystemConfig))
}

// GetConfigByKey 根据配置键获取配置
func GetConfigByKey(key string) (string, error) {
	db := orm.NewOrm()
	config := &SystemConfig{}
	err := db.QueryTable("app_system_config").
		Filter("config_key", key).
		One(config)
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// GetAllConfigs 获取所有配置（返回map）
func GetAllConfigs() (map[string]string, error) {
	db := orm.NewOrm()
	var configs []SystemConfig
	_, err := db.QueryTable("app_system_config").All(&configs)
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.ConfigKey] = config.ConfigValue
	}
	return configMap, nil
}

// GetConfigList 获取配置列表（支持分页）
func GetConfigList(page, pageSize int) (configs []SystemConfig, total int64, err error) {
	db := orm.NewOrm()
	qs := db.QueryTable("app_system_config")

	// 统计总数
	total, _ = qs.Count()

	// 分页查询
	offset := (page - 1) * pageSize
	_, err = qs.OrderBy("-id").Limit(pageSize, offset).All(&configs)
	return
}

// CreateConfig 创建新配置
func CreateConfig(key, value, desc string) error {
	db := orm.NewOrm()
	now := time.Now().Unix()

	config := &SystemConfig{
		ConfigKey:   key,
		ConfigValue: value,
		ConfigDesc:  desc,
		CreatedTime: now,
		UpdatedTime: now,
	}

	_, err := db.Insert(config)
	return err
}

// UpdateConfigByID 更新配置（按 ID）
func UpdateConfigByID(id int64, value, desc string) error {
	db := orm.NewOrm()
	config := &SystemConfig{ID: id}
	err := db.Read(config)
	if err != nil {
		return err
	}

	config.ConfigValue = value
	if desc != "" {
		config.ConfigDesc = desc
	}
	config.UpdatedTime = time.Now().Unix()

	_, err = db.Update(config, "config_value", "config_desc", "updated_time")
	return err
}

// DeleteConfigByID 删除配置（按 ID）
func DeleteConfigByID(id int64) error {
	db := orm.NewOrm()
	config := &SystemConfig{ID: id}

	// 先检查是否存在
	err := db.Read(config)
	if err != nil {
		return err
	}

	_, err = db.Delete(config)
	return err
}

// GetSystemConfigValue 获取配置值并转换为float64
func GetSystemConfigValue(key string) (float64, error) {
	valueStr, err := GetConfigByKey(key)
	if err != nil {
		return 0, err
	}

	// Trim 空格，防止解析失败
	valueStr = strings.TrimSpace(valueStr)

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
