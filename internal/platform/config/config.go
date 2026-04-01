package config

import (
	"github.com/luckysxx/common/conf"
)

// Config 全局配置
type Config struct {
	AppEnv      string              `mapstructure:"app_env"`
	Server      conf.ServerConfig   `mapstructure:"server"`
	GRPCServer  GRPCServerConfig    `mapstructure:"grpc_server"`
	IDGenerator IDGeneratorConfig   `mapstructure:"id_generator"`
	Database    conf.DatabaseConfig `mapstructure:"database"`
	Redis       conf.RedisConfig    `mapstructure:"redis"`
	OTel        conf.OTelConfig     `mapstructure:"otel"`
	Metrics     MetricsConfig       `mapstructure:"metrics"`
}

type GRPCServerConfig struct {
	Port string `mapstructure:"port"`
}

type IDGeneratorConfig struct {
	Addr string `mapstructure:"addr"`
}

type MetricsConfig struct {
	Port string `mapstructure:"port"`
}

// LoadConfig 从 Viper 加载配置
func LoadConfig() *Config {
	var cfg Config
	conf.Load(&cfg)
	return &cfg
}
