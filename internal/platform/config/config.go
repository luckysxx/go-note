package config

import (
	"github.com/luckysxx/common/conf"
)

// Config 全局配置
type Config struct {
	AppEnv   string              `mapstructure:"app_env"`
	Server   conf.ServerConfig   `mapstructure:"server"`
	Database conf.DatabaseConfig `mapstructure:"database"`
	Redis    conf.RedisConfig    `mapstructure:"redis"`
	OTel     conf.OTelConfig     `mapstructure:"otel"`
}

// LoadConfig 从 Viper 加载配置
func LoadConfig() *Config {
	var cfg Config
	conf.Load(&cfg)
	return &cfg
}
