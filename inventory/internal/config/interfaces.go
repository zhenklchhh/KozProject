package config

import "time"

type Config interface {
	GRPC() GrpcConfig
	Mongo() MongoConfig
	HTTP() HttpConfig
	Logger() LoggerConfig
}

type MongoConfig interface {
	Database() string
	URI() string
}

type GrpcConfig interface {
	Address() string
}

type HttpConfig interface {
	Address() string
	GetReadHeaderTimeout() time.Duration
	GetPingTimeout() time.Duration
	StaticDir() string
	GetSwaggerFile() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}