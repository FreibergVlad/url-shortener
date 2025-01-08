module github.com/FreibergVlad/url-shortener/api-gateway

go 1.23.2

require github.com/FreibergVlad/url-shortener/proto v0.0.0

require (
	github.com/FreibergVlad/url-shortener/shared/go v0.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0
	github.com/rs/zerolog v1.33.0
	google.golang.org/grpc v1.67.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240920164238-5a7b106cbb87.2 // indirect
	github.com/caarlos0/env/v11 v11.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	go.mongodb.org/mongo-driver/v2 v2.0.0-beta2 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

replace github.com/FreibergVlad/url-shortener/shared/go => ../../shared/go

replace github.com/FreibergVlad/url-shortener/proto => ../../proto
