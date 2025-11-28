package config

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl-mode"`
	User     string // из ENV
	Password string // из Env
}
