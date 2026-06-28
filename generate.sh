#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

# Generate ogen client for hosted API
echo "Generating hosted API client..."
go run github.com/ogen-go/ogen/cmd/ogen@v1.22.0 \
    -target ./internal/ogenhosted \
    -package ogenhosted \
    -clean \
    ./openapi/hosted.json

echo "Running go mod tidy..."
go mod tidy

echo "Done!"
