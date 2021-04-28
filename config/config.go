package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	env      string
	HTTP     HTTPServerConfig
	DBConfig *DBConfig
	JWT      *JWTConfig
}

func (c *Config) GetEnv() string {
	return c.env
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		env: viper.GetString("APP_ENV"),
		HTTP: HTTPServerConfig{
			host: viper.GetString("HTTP_SERVER"),
			port: viper.GetInt("HTTP_PORT"),
		},
		DBConfig: &DBConfig{
			port:     viper.GetInt("DB_PORT"),
			server:   viper.GetString("DB_SERVER"),
			user:     viper.GetString("DB_USER"),
			password: viper.GetString("DB_PASSWORD"),
			database: viper.GetString("DB_DATABASE"),
		},
		JWT: &JWTConfig{
			secret: viper.GetString("JWT_SECRET"),
		},
	}, nil
}
