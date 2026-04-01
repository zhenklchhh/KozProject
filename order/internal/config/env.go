package config

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	httpConfig            HttpConfig
	ctxConfig             ContextConfig
	inventoryClientConfig ClientConfig
	paymentClientConfig   ClientConfig
	migrationsConfig      MigrationsConfig
	postgresConfig        PostgresConfig
}

func (c *appConfig) HTTP() HttpConfig {
	return c.httpConfig
}

func (c *appConfig) Context() ContextConfig {
	return c.ctxConfig
}

func (c *appConfig) InventoryClient() ClientConfig {
	return c.inventoryClientConfig
}

func (c *appConfig) PaymentClient() ClientConfig {
	return c.paymentClientConfig
}

func (c *appConfig) Migrations() MigrationsConfig {
	return c.migrationsConfig
}

func (c *appConfig) Postgres() PostgresConfig {
	return c.postgresConfig
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

type contextEnvConfig struct {
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
}

func newContextConfig() (ContextConfig, error) {
	var cfg contextEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *contextEnvConfig) GetShutdownTimeout() time.Duration {
	return c.ShutdownTimeout
}

type inventoryClientEnvConfig struct {
	URL string `env:"INVENTORY_CLIENT_URL, required"`
}

func newInventoryClientConfig() (ClientConfig, error) {
	var cfg inventoryClientEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *inventoryClientEnvConfig) URI() string {
	return c.URL
}

type paymentClientEnvConfig struct {
	URL string `env:"PAYMENT_CLIENT_URL, required"`
}

func newPaymentClientConfig() (ClientConfig, error) {
	var cfg paymentClientEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *paymentClientEnvConfig) URI() string {
	return c.URL
}

type migrationEnvConfig struct {
	Directory string `env:"MIGRATIONS_DIR, required"`
}

func newMigrationsConfig() (MigrationsConfig, error) {
	var cfg migrationEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *migrationEnvConfig) Dir() string {
	return c.Directory
}

type postgresEnvConfig struct {
	Username string `env:"POSTGRES_USER, required"`
	Password string `env:"POSTGRES_PASSWORD, required"`
	DB       string `env:"POSTGRES_DB, required"`
	URL      string `env:"DB_URI, required"`
}

func newPostgresConfig() (PostgresConfig, error) {
	var cfg postgresEnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *postgresEnvConfig) Database() string {
	return c.DB
}

func (c *postgresEnvConfig) URI() string {
	return c.URL
}

func Load(path string) (Config, error) {
	_ = godotenv.Load(path)
	httpConfig, err := newHttpConfig()
	if err != nil {
		return nil, err
	}

	contextConfig, err := newContextConfig()
	if err != nil {
		return nil, err
	}

	invClientConfig, err := newInventoryClientConfig()
	if err != nil {
		return nil, err
	}

	paymentClientConfig, err := newPaymentClientConfig()
	if err != nil {
		return nil, err
	}

	migrationConfig, err := newMigrationsConfig()
	if err != nil {
		return nil, err
	}

	postgresConfig, err := newPostgresConfig()
	if err != nil {
		return nil, err
	}
	return &appConfig{
		httpConfig:            httpConfig,
		ctxConfig:             contextConfig,
		inventoryClientConfig: invClientConfig,
		paymentClientConfig:   paymentClientConfig,
		migrationsConfig:      migrationConfig,
		postgresConfig:        postgresConfig,
	}, nil
}
