# syntax=docker/dockerfile:1

FROM golang:1.23.2-alpine AS build

WORKDIR /build/services/api-gateway

COPY shared/go ../../shared/go
COPY proto ../../proto
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./
RUN go mod download

COPY services/api-gateway ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api-gateway ./cmd/api-gateway/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS api-gateway

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/api-gateway/api-gateway .

ENTRYPOINT ["/home/nonroot/api-gateway"]