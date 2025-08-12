#!/bin/bash
# Install locally for testing
go build -o palettesmith ./cmd/palettesmith
sudo mv palettesmith /usr/local/bin/
