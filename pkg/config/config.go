package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type App struct {
	Server   Server   `mapstructure:"server"`
	JWT      JWT      `mapstructure:"jwt"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}

type Server struct {
	Addr string `mapstructure:"addr"`
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
	// Secret string `mapstructure:"secret"`
}

type JWT struct {
	ActiveKID             string        `mapstructure:"active_key_id"`
	KeyDirPath            string        `mapstructure:"key_dir_path"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`

	MaxOpenConns    int32         `mapstructure:"max_open_conns"`
	MinOpenConns    int32         `mapstructure:"min_open_conns"`
	MaxIdleConns    int32         `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func New(path string) (App, error) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		// Fallback default jika path kosong
		viper.SetConfigFile("./config.yaml")
	}

	var cfg App
	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
