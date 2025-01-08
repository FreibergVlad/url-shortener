# syntax=docker/dockerfile:1

FROM golang:1.23.2-alpine AS build

WORKDIR /build/services/short-url-redirection-service

COPY shared/go ../../shared/go
COPY proto ../../proto
COPY services/short-url-redirection-service/go.mod services/short-url-redirection-service/go.sum ./
RUN go mod download

COPY services/short-url-redirection-service ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o short-url-redirection-service ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS server

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/short-url-redirection-service/short-url-redirection-service .

ENTRYPOINT ["/home/nonroot/short-url-redirection-service"]
