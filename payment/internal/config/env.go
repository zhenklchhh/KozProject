package config

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	cfg *appConfig
)

type appConfig struct {
	grpcConfig GrpcConfig
	httpConfig HttpConfig
	loggerConfig LoggerConfig
}

func (c *appConfig) GRPC() GrpcConfig {
	return c.grpcConfig
}

func (c *appConfig) HTTP() HttpConfig {
	return c.httpConfig
}

func (c *appConfig) Logger() LoggerConfig {
	return c.loggerConfig
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

type loggerEnvConfig struct {
	LoggerLevel string `env:"LOGGER_LEVEL,required" envDefault:""`
	LogsAsJson bool `env:"LOGGER_AS_JSON,required"`
}

func newLoggerConfig() (LoggerConfig, error) {
	var cfg loggerEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *loggerEnvConfig) Level() string {
	return c.LoggerLevel
}

func (c *loggerEnvConfig) AsJson() bool {
	return c.LogsAsJson
}

func Load(path string) error {
	_ = godotenv.Load(path)

	grpcConfig, err := newGRPCConfig()
	if err != nil {
		return err
	}

	httpConfig, err := newHttpConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := newLoggerConfig()
	if err != nil {
		return err
	}

	cfg = &appConfig{
		grpcConfig: grpcConfig,
		httpConfig: httpConfig,
		loggerConfig: loggerConfig,
	}
	return nil
}

func AppConfig() Config {
	return cfg
}
