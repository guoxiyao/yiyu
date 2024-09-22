package config

import (
	"errors"
	"os"
)

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	//Database *DatabaseConfig
}

// LoadConfig 加载配置
func LoadConfig() (*DatabaseConfig, error) {
	dbConfig := &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	// 检查必要配置是否已设置
	if dbConfig.Host == "" || dbConfig.User == "" || dbConfig.Name == "" {
		return nil, ErrConfigMissing
	}

	return dbConfig, nil
}

// ErrConfigMissing 配置缺失错误
var ErrConfigMissing = errors.New("one or more database configuration values are missing")
