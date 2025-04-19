package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port int    `env:"PORT" envDefault:"8080"`
	DBHost string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort int    `env:"DB_PORT" envDefault:"3306"`
	DBUser string `env:"DB_USER" envDefault:"todo_user"`
	DBPass string `env:"DB_PASS" envDefault:"p@ssw0rd"`
	DBName string `env:"DB_NAME" envDefault:"todo_app"`
	DBLoc  string `env:"DB_LOC"  envDefault:"Asia/Tokyo"`
}

func New() (*Config, string, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, "", err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=%s",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBLoc)
	return cfg, dsn, nil
}