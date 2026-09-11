FROM golang:1.21-alpine AS builder

WORKDIR /app

# Dependency download
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Source code copy
COPY . .

# Build step
RUN go build -o /meshchat ./cmd/meshchat

FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /meshchat /meshchat

# No secrets should be baked in. Config provided via environment or mounted volume.
ENTRYPOINT ["/meshchat"]
