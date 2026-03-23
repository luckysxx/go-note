package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	AppEnv       string             `mapstructure:"app_env"`
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Redis        RedisConfig        `mapstructure:"redis"`
	UserPlatform UserPlatformConfig `mapstructure:"user_platform"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver      string `mapstructure:"driver"`
	Source      string `mapstructure:"source"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// UserPlatformConfig user-platform gRPC 连接配置
type UserPlatformConfig struct {
	Addr string `mapstructure:"addr"`
}

// LoadConfig 从 Viper 加载配置
func LoadConfig() *Config {
	// 优先加载 .env 文件（生产环境一般不用此文件）
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // 支持项目根目录启动
	viper.AddConfigPath(".")

	// 允许环境变量覆盖配置 (比如 export DATABASE_SOURCE=xyz)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config.yaml found, relying entirely on ENV variables: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
	return &cfg
}
