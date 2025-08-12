#!/bin/bash
# Build binary for current platform
CGO_ENABLED=0 go build -ldflags="-s -w" -o palettesmith ./cmd/palettesmith
