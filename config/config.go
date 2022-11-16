package config

import (
    "fmt"
    "github.com/evrone/go-clean-template/pkg/logger/prepare"

    "github.com/ilyakaznacheev/cleanenv"
)

type (
    // Config -.
    Config struct {
        App   App            `yaml:"app"`
        HTTP  HTTP           `yaml:"http"`
        Log   prepare.Config `yaml:"logger"`
        PG    PG             `yaml:"postgres"`
        RMQ   RMQ            `yaml:"rabbitmq"`
        GRPC  GRPC           `yaml:"grpc"`
        Redis Redis          `yaml:"redis"`
    }

    // App -.
    App struct {
        Name         string `env-required:"true" yaml:"name"    env:"APP_NAME"`
        Version      string `env-required:"true" yaml:"version" env:"APP_VERSION"`
        ResourcesDir string `env-required:"true" yaml:"resources_dir" env:"RESOURCES_DIR"`
    }

    // HTTP -.
    HTTP struct {
        Port       string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
        MetricPort string `env-required:"true" yaml:"metric_port" env:"HTTP_METRIC_PORT"`
    }

    // PG -.
    PG struct {
        PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
        URL     string `env-required:"true"                 env:"PG_URL"`
    }

    // Redis -.
    Redis struct {
        URL string `env-required:"true"                 env:"REDIS_URL"`
    }

    // GRPC -.
    GRPC struct {
        URL string `env-required:"true"                 env:"GRPC_URL"`
    }

    // RMQ -.
    RMQ struct {
        ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
        ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
        URL            string `env-required:"true"                            env:"RMQ_URL"`
    }
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
    cfg := &Config{}

    err := cleanenv.ReadConfig("./config/config.yml", cfg)
    if err != nil {
        return nil, fmt.Errorf("config error: %w", err)
    }

    err = cleanenv.ReadEnv(cfg)
    if err != nil {
        return nil, err
    }

    return cfg, nil
}
