package config

import "time"

type Config interface {
	HTTP() HttpConfig
	Context() ContextConfig
	InventoryClient() ClientConfig
	PaymentClient() ClientConfig
	Migrations() MigrationsConfig
	Postgres() PostgresConfig
	Logger() LoggerConfig
}

type LoggerConfig interface{
	Level() string
	AsJson() bool
}

type HttpConfig interface {
	Address() string
	GetReadHeaderTimeout() time.Duration
}

type ContextConfig interface {
	GetShutdownTimeout() time.Duration
}

type ClientConfig interface {
	URI() string
}

type MigrationsConfig interface {
	Dir() string
}

type PostgresConfig interface {
	Database() string
	URI() string
}
