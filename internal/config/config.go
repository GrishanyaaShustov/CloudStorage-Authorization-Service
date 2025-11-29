package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env      string         `mapstructure:"env"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
}

func MustLoad() *Config {
	viper.SetConfigName("prod")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// Подгружаем .env
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()

	// Ошибка чтения файла -> PANIC
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config: " + err.Error())
	}

	var cfg Config

	// Ошибка при маршалинге -> PANIC
	if err := viper.Unmarshal(&cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	// Подтягиваем креды из ENV
	cfg.Database.Password = viper.GetString("DB_PASSWORD")
	cfg.Database.User = viper.GetString("DB_USER")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")

	if cfg.Database.Password == "" || cfg.Database.User == "" {
		panic("DATABASE credentials are missing (DB_USER / DB_PASSWORD not set)")
	}

	if cfg.Redis.Password == "" {
		panic("Redis credentials are missing password")
	}

	return &cfg
}
