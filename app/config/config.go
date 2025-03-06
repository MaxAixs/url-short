package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Env        string     `mapstructure:"env" yaml:"env"`
	HTTPServer HTTPServer `mapstructure:"http_server" yaml:"http_server"`
	DBConfig   DBConfig   `mapstructure:"db" yaml:"db"`
}

type HTTPServer struct {
	Port              string        `mapstructure:"port" yaml:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"reader_header_timeout" yaml:"reader_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"db_name" yaml:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
}

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("internal/config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w ", err)
	}

	err = setConfigEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	return &cfg, nil
}

func setConfigEnv(cfg *Config) error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error read from env file %w", err)
	}

	cfg.DBConfig.Password = viper.GetString("DB_PASSWORD")
	if cfg.DBConfig.Password == "" {
		return fmt.Errorf("cant get DB password from env %w", err)
	}

	return nil
}
