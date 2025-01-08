module github.com/FreibergVlad/url-shortener/domain-service

go 1.23.2

replace github.com/FreibergVlad/url-shortener/shared/go => ../../shared/go

replace github.com/FreibergVlad/url-shortener/proto => ../../proto

require (
	github.com/FreibergVlad/url-shortener/proto v0.0.0
	github.com/FreibergVlad/url-shortener/shared/go v0.0.0
	github.com/brianvoe/gofakeit/v7 v7.1.2
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.9.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240920164238-5a7b106cbb87.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/bufbuild/protovalidate-go v0.7.2 // indirect
	github.com/caarlos0/env/v11 v11.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/cel-go v0.21.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.mongodb.org/mongo-driver/v2 v2.0.0-beta2 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
