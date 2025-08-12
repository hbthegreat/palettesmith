#!/bin/bash
# Development mode - hot reload with air or simple rebuild
# Install: go install github.com/cosmtrek/air@latest
air || while true; do go run ./cmd/palettesmith "$@"; done
