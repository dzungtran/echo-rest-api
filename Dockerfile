# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies in a single layer
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
ARG APP_PATH
RUN go build -ldflags="-s -w" -buildvcs=false -o appbin $APP_PATH

# Runtime stage
FROM alpine

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --chown=65532:65532 --from=builder /app/appbin .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

ENTRYPOINT ["/app/appbin"]
