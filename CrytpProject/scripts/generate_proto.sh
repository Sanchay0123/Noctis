#!/bin/sh
set -e

# Run inside Docker to avoid host dependencies and ensure reproducibility
# Generator version is controlled via the pinned base image (golang:1.21-alpine) 
# and the explicitly pinned protoc-gen-go version (v1.33.0).

docker run --rm -v "$(pwd):/app" -w /app golang:1.21-alpine sh -c '
  echo "Installing protoc..."
  apk add --no-cache protoc &&
  echo "Installing protoc-gen-go v1.33.0..." &&
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 &&
  export PATH="$PATH:$(go env GOPATH)/bin" &&
  echo "Generating Protobuf files..." &&
  protoc --go_out=. --go_opt=paths=source_relative internal/protocol/meshchat.proto
'

echo "Protobuf generation complete."
