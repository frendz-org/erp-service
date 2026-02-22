# ---- Build stage ----
# TARGETPLATFORM and TARGETARCH are set automatically by Docker BuildKit
# based on the machine or --platform flag.
#   Mac M1/M2/M3 (Apple Silicon) → TARGETARCH=arm64
#   Linux x86 VPS / Intel Mac     → TARGETARCH=amd64
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETARCH

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /app/bin/erp-service ./cmd/http

# ---- Runtime stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/bin/erp-service /app/erp-service
COPY --from=builder /app/migration /app/migration
COPY --from=builder /app/doc/openapi /app/doc/openapi

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/erp-service"]
