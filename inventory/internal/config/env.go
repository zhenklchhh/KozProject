package config

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var cfg *appConfig

type appConfig struct {
	grpcConfig  GrpcConfig
	httpConfig  HttpConfig
	mongoConfig MongoConfig
	loggerConfig LoggerConfig
}

func (c *appConfig) GRPC() GrpcConfig   { return c.grpcConfig }
func (c *appConfig) HTTP() HttpConfig   { return c.httpConfig }
func (c *appConfig) Mongo() MongoConfig { return c.mongoConfig }
func (c *appConfig) Logger() LoggerConfig { return c.loggerConfig }

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
	PingTimeout       time.Duration `env:"PING_TIMEOUT,required"`
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
	mongoConfig, err := newMongoConfig()
	if err != nil {
		return err
	}
	loggerConfig, err := newLoggerConfig()
	if err != nil {
		return err
	}
	cfg = &appConfig{
		grpcConfig:  grpcConfig,
		httpConfig:  httpConfig,
		mongoConfig: mongoConfig,
		loggerConfig: loggerConfig,
	}
	return nil
}

func AppConfig() Config {
	return cfg
}
