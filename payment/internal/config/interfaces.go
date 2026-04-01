package config

import "time"

type Config interface {
	GRPC() GrpcConfig
	HTTP() HttpConfig
}

type GrpcConfig interface {
	Address() string
}

type HttpConfig interface {
	Address() string
	GetReadHeaderTimeout() time.Duration
}
