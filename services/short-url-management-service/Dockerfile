# syntax=docker/dockerfile:1

FROM golang:1.23.2-alpine AS build

WORKDIR /build/services/short-url-management-service

COPY shared/go ../../shared/go
COPY proto ../../proto
COPY services/short-url-management-service/go.mod services/short-url-management-service/go.sum ./
RUN go mod download

COPY services/short-url-management-service ./
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o short-url-management-service ./cmd/short-url-management-service/main.go && \
    go build -ldflags="-w -s" -o short-url-generator-service ./cmd/short-url-generator-service/main.go && \
    go build -ldflags="-w -s" -o short-url-service-migrations ./cmd/migrations/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS migrations

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/short-url-management-service/short-url-service-migrations .
COPY --from=build --chown=root:root --chmod=755 /build/services/short-url-management-service/internal/db/migrations db/migrations/

ENTRYPOINT ["/home/nonroot/short-url-service-migrations"]

FROM gcr.io/distroless/static-debian12:nonroot AS short-url-management-service

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/short-url-management-service/short-url-management-service .

ENTRYPOINT ["/home/nonroot/short-url-management-service"]

FROM gcr.io/distroless/static-debian12:nonroot AS short-url-generator-service

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/short-url-management-service/short-url-generator-service .

ENTRYPOINT ["/home/nonroot/short-url-generator-service"]