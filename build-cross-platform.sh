#!/bin/bash

# Cross-platform build script for KrankyBear Tetris
# This script builds for Windows, macOS, and Linux from a single machine
# Note: Cross-compiling Fyne apps requires proper CGO setup

set -e

echo "KrankyBear Tetris - Cross-Platform Build Script"
echo "================================================"
echo ""

# Create bin directory if it doesn't exist
if [ ! -d "bin" ]; then
    mkdir -p bin
fi

# Cleanup previous binaries
rm -f bin/tetris-*

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Update dependencies
echo "Updating dependencies..."
go get fyne.io/fyne/v2@latest || true
go mod tidy
go mod vendor

echo ""
echo "Building for multiple platforms..."
echo ""

# Build for Linux (amd64)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/tetris-linux-amd64
if [ $? -eq 0 ]; then
    echo "✓ Linux (amd64) build successful"
else
    echo "✗ Linux (amd64) build failed"
fi

# Build for Linux (arm64)
echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/tetris-linux-arm64
if [ $? -eq 0 ]; then
    echo "✓ Linux (arm64) build successful"
else
    echo "✗ Linux (arm64) build failed (may require native build)"
fi

# Build for macOS (amd64) - only works on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Building for macOS (amd64)..."
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/tetris-macos-amd64
    if [ $? -eq 0 ]; then
        echo "✓ macOS (amd64) build successful"
    else
        echo "✗ macOS (amd64) build failed"
    fi
    
    echo "Building for macOS (arm64)..."
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/tetris-macos-arm64
    if [ $? -eq 0 ]; then
        echo "✓ macOS (arm64) build successful"
    else
        echo "✗ macOS (arm64) build failed"
    fi
else
    echo "⚠ Skipping macOS builds (requires macOS host)"
fi

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
# Check for MinGW cross-compiler
if command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then
    export CC=x86_64-w64-mingw32-gcc
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w -H windowsgui" -trimpath -o bin/tetris-windows-amd64.exe
    if [ $? -eq 0 ]; then
        echo "✓ Windows (amd64) build successful"
    else
        echo "✗ Windows (amd64) build failed"
    fi
    unset CC
else
    echo "⚠ Skipping Windows build (requires x86_64-w64-mingw32-gcc)"
    echo "  Install with: sudo apt-get install gcc-mingw-w64-x86-64 (Linux)"
    echo "  or: brew install mingw-w64 (macOS)"
fi

echo ""
echo "================================================"
echo "Build complete! Binaries are in the bin/ directory:"
ls -lh bin/tetris-* 2>/dev/null || echo "No binaries were built"

echo ""
echo "Note:"
echo "- Linux builds: May require native compilation for best results"
echo "- macOS builds: Must be built on macOS"
echo "- Windows builds: Require MinGW cross-compiler or native Windows build"

