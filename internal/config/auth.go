package config

type AuthConfig struct {
	SecretKey  string `mapstructure:"secret_key"`
	Realm      string `mapstructure:"realm" default:"api"`
	Timeout    int    `mapstructure:"timeout_hours" default:"24"`
	MaxRefresh int    `mapstructure:"max_refresh_hours" default:"24"`
}
