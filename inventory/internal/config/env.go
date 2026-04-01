package config

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	grpcConfig  GrpcConfig
	httpConfig  HttpConfig
	mongoConfig MongoConfig
}

func (c *appConfig) GRPC() GrpcConfig   { return c.grpcConfig }
func (c *appConfig) HTTP() HttpConfig   { return c.httpConfig }
func (c *appConfig) Mongo() MongoConfig { return c.mongoConfig }

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
	PingTimeout       time.Duration `env:"PING_TIMEOUT, required"`
	StaticDirectory   string        `env:"HTTP_STATIC_DIR,required"`
	SwaggerFile       string        `env:"HTTP_SWAGGER_FILE,required"`
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

func (c *httpEnvConfig) GetPingTimeout() time.Duration {
	return c.PingTimeout
}

func (c *httpEnvConfig) GetSwaggerFile() string {
	return c.SwaggerFile
}

func (c *httpEnvConfig) StaticDir() string {
	return c.StaticDirectory
}

type mongoEnvConfig struct {
	Username string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	DB       string `env:"MONGO_INITDB_ROOT_DATABASE,required"`
	URL      string `env:"MONGO_DB_URI,required"`
}

func newMongoConfig() (MongoConfig, error) {
	var cfg mongoEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *mongoEnvConfig) URI() string {
	return c.URL
}

func (c *mongoEnvConfig) Database() string {
	return c.DB
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
	mongoConfig, err := newMongoConfig()
	if err != nil {
		return nil, err
	}
	return &appConfig{
		grpcConfig:  grpcConfig,
		httpConfig:  httpConfig,
		mongoConfig: mongoConfig,
	}, nil
}
