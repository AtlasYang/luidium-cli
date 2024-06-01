#!/bin/bash

BINARY_NAME="luidium"

CURRENT_DIR="$(cd "$(dirname "$0")" && pwd)"

case "$(uname -s)" in
    Linux*)     BINARY_PATH="$CURRENT_DIR/luidium-linux-amd64";;
    Darwin*)    BINARY_PATH="$CURRENT_DIR/luidium-darwin-amd64";;
    *)          echo "Unsupported OS"; exit 1;;
esac

SYSTEM_BIN_DIR="/usr/local/bin"

sudo cp "$BINARY_PATH" "$SYSTEM_BIN_DIR/$BINARY_NAME"

echo "Luidium CLI installed successfully"