module github.com/FreibergVlad/url-shortener/shared/go

go 1.23.2

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.5-20250219170025-d39267d9df8f.1
	github.com/FreibergVlad/url-shortener/proto v0.0.0-00010101000000-000000000000
	github.com/bufbuild/protovalidate-go v0.9.2
	github.com/caarlos0/env/v11 v11.3.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.7.1
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.10.0
	go.mongodb.org/mongo-driver/v2 v2.0.1
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250224174004-546df14abb99
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	cel.dev/expr v0.21.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/cel-go v0.24.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/exp v0.0.0-20250218142911-aa4b98e5adaa // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250224174004-546df14abb99 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/FreibergVlad/url-shortener/proto => ../../proto
