#!/bin/sh
# *******************************************************************************
# Copyright (c) 2025 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Apache License Version 2.0 which is available at
# https://www.apache.org/licenses/LICENSE-2.0
#
# SPDX-License-Identifier: Apache-2.0
# *******************************************************************************
#
# scorex installer script
# Usage: curl -sSL https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/install.sh | sh
#

set -e

REPO="eclipse-score/score_scrample"
INSTALL_DIR="${SCOREX_INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s)
    ARCH=$(uname -m)

    case "$OS" in
        Linux)
            case "$ARCH" in
                x86_64) PLATFORM="linux-x86_64" ;;
                *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
            esac
            ;;
        Darwin)
            case "$ARCH" in
                arm64) PLATFORM="macos-arm64" ;;
                x86_64) PLATFORM="macos-x86_64" ;;
                *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
            esac
            ;;
        *)
            echo "Unsupported OS: $OS"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Could not fetch latest version"
        exit 1
    fi
    echo "Latest version: $VERSION"
}

# Download and install
install_scorex() {
    echo "Detecting platform..."
    detect_platform
    echo "Platform: $PLATFORM"

    echo "Fetching latest version..."
    get_latest_version

    DOWNLOAD_URL="https://github.com/$REPO/releases/download/v$VERSION/scorex-$VERSION-$PLATFORM.tar.gz"
    TEMP_DIR=$(mktemp -d)

    echo "Downloading scorex from $DOWNLOAD_URL..."
    curl -sL "$DOWNLOAD_URL" | tar -xz -C "$TEMP_DIR"

    echo "Installing to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"

    # Find the binary in temp dir
    BINARY=$(find "$TEMP_DIR" -name "scorex-*" -type f | head -n 1)
    if [ -z "$BINARY" ]; then
        echo "Error: Binary not found in download"
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    mv "$BINARY" "$INSTALL_DIR/scorex"
    chmod +x "$INSTALL_DIR/scorex"

    rm -rf "$TEMP_DIR"

    echo ""
    echo "✓ scorex $VERSION installed successfully!"
    echo ""
    echo "Make sure $INSTALL_DIR is in your PATH:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo ""
    echo "Run 'scorex --help' to get started."
}

install_scorex
