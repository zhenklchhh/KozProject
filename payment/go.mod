module github.com/zhenklchhh/KozProject/payment

go 1.26.1

require (
	buf.build/go/protovalidate v1.1.3
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	github.com/joho/godotenv v1.5.1
	github.com/zhenklchhh/KozProject/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.2
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/brianvoe/gofakeit/v7 v7.14.1 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/zhenklchhh/KozProject/shared => ../shared
