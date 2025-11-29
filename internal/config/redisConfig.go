package config

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string // from ENV
	DB       int    `mapstructure:"db"`
}
