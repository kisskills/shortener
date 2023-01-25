package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"strings"
)

type Config struct {
	cfg *koanf.Koanf
}

func New(configPath string) (*Config, error) {
	cfg := koanf.New(".")

	err := cfg.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		return nil, err
	}

	//В docker compose создаем новую переменную в env, нужную для выбора хранилища
	err = cfg.Load(env.Provider("STORAGE_", "=", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "STORAGE_"))
	}), nil)
	if err != nil {
		return nil, err
	}

	return &Config{
		cfg: cfg,
	}, nil
}

func (c *Config) DatabaseMode() string {
	return c.cfg.String("mode")
}

func (c *Config) PostgresDSN() string {
	return c.cfg.String("storage.postgres")
}

func (c *Config) GrpcPort() int {
	return c.cfg.Int("service.grpc_port")
}

func (c *Config) HttpPort() int {
	return c.cfg.Int("service.http_port")
}
