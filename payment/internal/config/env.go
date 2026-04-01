package config

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	grpcConfig GrpcConfig
	httpConfig HttpConfig
}

func (c *appConfig) GRPC() GrpcConfig {
	return c.grpcConfig
}

func (c *appConfig) HTTP() HttpConfig {
	return c.httpConfig
}

type grpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

func newGRPCConfig() (GrpcConfig, error) {
	var cfg grpcEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *grpcEnvConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type httpEnvConfig struct {
	Host              string        `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port              string        `env:"HTTP_PORT,required"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT,required"`
}

func newHttpConfig() (HttpConfig, error) {
	var cfg httpEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *httpEnvConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *httpEnvConfig) GetReadHeaderTimeout() time.Duration {
	return c.ReadHeaderTimeout
}

func Load(path string) (Config, error) {
	_ = godotenv.Load(path)

	grpcConfig, err := newGRPCConfig()
	if err != nil {
		return nil, err
	}

	httpConfig, err := newHttpConfig()
	if err != nil {
		return nil, err
	}

	return &appConfig{
		grpcConfig: grpcConfig,
		httpConfig: httpConfig,
	}, nil
}
