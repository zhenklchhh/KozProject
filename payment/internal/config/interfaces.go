package config

import "time"

type Config interface {
	GRPC() GrpcConfig
	HTTP() HttpConfig
	Logger() LoggerConfig
}

type GrpcConfig interface {
	Address() string
}

type HttpConfig interface {
	Address() string
	GetReadHeaderTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
