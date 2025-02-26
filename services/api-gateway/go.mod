module github.com/FreibergVlad/url-shortener/api-gateway

go 1.23.2

require github.com/FreibergVlad/url-shortener/proto v0.0.0

require (
	github.com/FreibergVlad/url-shortener/shared/go v0.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1
	github.com/rs/zerolog v1.33.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.5-20250219170025-d39267d9df8f.1 // indirect
	github.com/caarlos0/env/v11 v11.3.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.mongodb.org/mongo-driver/v2 v2.0.1 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250224174004-546df14abb99 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250224174004-546df14abb99 // indirect
)

replace github.com/FreibergVlad/url-shortener/shared/go => ../../shared/go

replace github.com/FreibergVlad/url-shortener/proto => ../../proto
