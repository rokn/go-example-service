package config

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST" default:"localhost"`
	Port     string `mapstructure:"DB_PORT" default:"5432"`
	User     string `mapstructure:"DB_USER" default:"postgres"`
	Password string `mapstructure:"DB_PASSWORD" default:"postgres"`
	Name     string `mapstructure:"DB_NAME" default:"postgres"`
}
