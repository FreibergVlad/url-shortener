# syntax=docker/dockerfile:1

FROM golang:1.23.2-alpine AS build

WORKDIR /build/services/auth-service

COPY shared/go ../../shared/go
COPY proto ../../proto
COPY services/auth-service/go.mod services/auth-service/go.sum ./
RUN go mod download

COPY services/auth-service ./
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o auth-service ./cmd/auth-service/main.go && \
    go build -ldflags="-w -s" -o auth-service-migrations ./cmd/migrations/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS auth-service

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/auth-service/auth-service .

ENTRYPOINT ["/home/nonroot/auth-service"]

FROM gcr.io/distroless/static-debian12:nonroot AS migrations

USER nonroot:nonroot

WORKDIR /home/nonroot

COPY --from=build --chown=root:root --chmod=755 /build/services/auth-service/internal/db/migrations /db/migrations/
COPY --from=build --chown=root:root --chmod=755 /build/services/auth-service/auth-service-migrations .

ENTRYPOINT ["/home/nonroot/auth-service-migrations"]