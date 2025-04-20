#!/bin/bash

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the root of the project (parent of scripts/)
cd "$SCRIPT_DIR/.." || exit 1

# Run the Go server
go run ./cmd/server
