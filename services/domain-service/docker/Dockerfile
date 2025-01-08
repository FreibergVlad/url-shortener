# syntax=docker/dockerfile:1

FROM golang:1.23.2-alpine AS build

WORKDIR /build/services/domain-service

COPY shared/go ../../shared/go
COPY proto ../../proto
COPY services/domain-service/go.mod services/domain-service/go.sum ./
RUN go mod download

COPY services/domain-service ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o domain-service ./cmd/domain-service/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS domain-service

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/domain-service/domain-service .

ENTRYPOINT ["/home/nonroot/domain-service"]
