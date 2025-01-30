package config

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST" default:"localhost"`
	Port     string `mapstructure:"REDIS_PORT" default:"6379"`
	Password string `mapstructure:"REDIS_PASSWORD" default:""`
	DB       int    `mapstructure:"REDIS_DB" default:"0"`
}
