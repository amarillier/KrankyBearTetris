.PHONY: all build test clean install run help

.PHONY: macos
macos: build-darwin
.PHONY: mac
macos: build-darwin

# Binary name
BINARY_NAME=tetris

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
LDFLAGS=-w -s

# Default target
all: test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v

# Build with WebView support (better fallback UI)
build-webview:
	@echo "Building $(BINARY_NAME) with WebView support..."
	@echo "Note: Linux requires webkit2gtk (see BUILD_WEBVIEW_LINUX.md)"
	$(GOGET) github.com/webview/webview_go || true
	$(GOMOD) tidy
	CGO_ENABLED=1 $(GOBUILD) -tags webview -o $(BINARY_NAME) -v

# Build for all desktop platforms
build-all: build-linux build-darwin build-windows

# Build for iOS (requires Xcode and iOS SDK, only works on macOS)
build-ios:
	@echo "Building for iOS..."
	@if [[ "$(shell uname)" != "Darwin" ]]; then \
		echo "Error: iOS builds require macOS and Xcode"; \
		false; \
	fi
	@if ! command -v fyne >/dev/null 2>&1; then \
		echo "Installing fyne command-line tool..."; \
		go install fyne.io/fyne/v2/cmd/fyne@latest; \
	fi
	@echo "Note: iOS builds require Xcode, iOS SDK, and proper code signing"
	@echo "Building iOS app bundle..."
	fyne package -os ios -appID com.github.amarillier.KrankyBearTetris -name "KrankyBear Tetris"
	@echo "iOS build complete. Output: KrankyBearTetris.app"

# Build for Android (requires Android SDK/NDK)
build-android:
	@echo "Building for Android..."
	@if ! command -v fyne >/dev/null 2>&1; then \
		echo "Installing fyne command-line tool..."; \
		go install fyne.io/fyne/v2/cmd/fyne@latest; \
	fi
	@echo "Note: Android builds require Android SDK and NDK"
	@echo "Set ANDROID_NDK_HOME environment variable if needed"
	@echo "Building Android APK..."
	fyne package -os android -appID com.github.amarillier.KrankyBearTetris -name "KrankyBear Tetris"
	@echo "Android build complete. Output: KrankyBearTetris.apk"

# Build for all platforms including mobile
build-all-platforms: build-all build-ios build-android

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@echo "Note: Cross-compiling Fyne apps from macOS to Linux requires fyne-cross or native Linux build"
	@echo "For direct compilation on Linux, use: ./compile-linux.sh"
	@echo ""
	@echo "Attempting cross-compile with fyne-cross (requires Docker)..."
	@if command -v fyne-cross >/dev/null 2>&1; then \
		fyne-cross linux -arch=amd64,arm64 -output $(BINARY_NAME)-linux; \
	else \
		echo "fyne-cross not found. Install with: go install github.com/fyne-io/fyne-cross@latest"; \
		echo "Or build directly on Linux using: ./compile-linux.sh"; \
		false; \
	fi

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) -ldflags="-w -s" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -ldflags="-w -s" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64
	# set executable icon
	./setIcon.sh Resources/Images/KrankyBearVikingHelmet.png $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64
	./setIcon.sh Resources/Images/KrankyBearVikingHelmet.png $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@echo "Note: Requires mingw-w64 (brew install mingw-w64 on macOS)"
	@echo "Note: Console window enabled so flags (-version, -help) work. Use Start-Process -WindowStyle Hidden to hide."
	@if command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) -ldflags="-w -s" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe -v; \
		./setIcon.sh Resources/Images/KrankyBearVikingHelmet.png $(BUILD_DIR)/$(BINARY_NAME)-windows.exe; \
	else \
		echo "mingw-w64 not found. Install with: brew install mingw-w64"; \
		echo "Or use: ./compile-windows.ps1 on Windows"; \
		false; \
	fi

build-windows-debug:
	@echo "Building Windows DEBUG version (with console output)..."
	@mkdir -p $(BUILD_DIR)
	@echo "Note: This version shows console output for troubleshooting"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-win-debug.exe -v
	@echo "Debug build created: $(BUILD_DIR)/$(BINARY_NAME)-win-debug.exe"

build-windows-webview:
	@echo "==========================================================================="
	@echo "ERROR: WebView cannot be cross-compiled from macOS to Windows"
	@echo "==========================================================================="
	@echo ""
	@echo "The WebView build requires Windows SDK headers (shlobj.h, etc.) that"
	@echo "are not available in the mingw-w64 cross-compilation toolchain."
	@echo ""
	@echo "To build with WebView support:"
	@echo ""
	@echo "  1. Copy source files to your Windows machine"
	@echo "  2. Install Go on Windows: https://go.dev/dl/"
	@echo "  3. Run these commands on Windows:"
	@echo "       go get github.com/webview/webview_go"
	@echo "       go mod tidy"
	@echo "       go build -tags webview -o template-webview.exe"
	@echo ""
	@echo "Or use the build script: build-webview-windows.bat"
	@echo ""
	@echo "See: WINDOWS_VM_BUILD_INSTRUCTIONS.md for details"
	@echo "==========================================================================="
	@false

build-windows-webview-debug:
	@echo "==========================================================================="
	@echo "ERROR: WebView cannot be cross-compiled from macOS to Windows"
	@echo "==========================================================================="
	@echo ""
	@echo "WebView requires building directly on Windows with the Windows SDK."
	@echo ""
	@echo "To build WebView debug version on Windows:"
	@echo "       go build -tags webview -o template-webview-debug.exe"
	@echo ""
	@echo "See: WINDOWS_VM_BUILD_INSTRUCTIONS.md for complete instructions"
	@echo "==========================================================================="
	@false

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f latestcheck.json

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) /usr/local/bin/

# Run the application with default settings
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Check if GUI is available
check-gui: build
	@echo "Checking GUI availability..."
	./$(BINARY_NAME) -check-gui

# Show application help
app-help: build
	@echo "Showing application help..."
	./$(BINARY_NAME) -help

# Show version
version: build
	@echo "Showing version information..."
	./$(BINARY_NAME) -version

# Show help
help:
	@echo "KrankyBear Tetris - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the application for current platform"
	@echo "  make build-webview  - Build with WebView support (better fallback UI)"
	@echo "  make build-all      - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make build-linux    - Build for Linux"
	@echo "  make build-darwin   - Build for macOS (Intel and ARM)"
	@echo "  make build-windows  - Build for Windows"
	@echo "  make build-windows-debug - Build Windows version with console output (for troubleshooting)"
	@echo "  make build-windows-webview - Build Windows with WebView support (better UI than MessageBox)"
	@echo "  make build-windows-webview-debug - Build Windows WebView with console output"
	@echo "  make build-ios      - Build for iOS (requires macOS, Xcode, iOS SDK)"
	@echo "  make build-android  - Build for Android (requires Android SDK/NDK)"
	@echo "  make build-all-platforms - Build for all platforms including mobile"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make bench          - Run benchmarks"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Install/update dependencies"
	@echo "  make install        - Install binary to /usr/local/bin"
	@echo "  make run            - Build and run the application"
	@echo "  make check-gui      - Check if GUI is available"
	@echo "  make app-help       - Show application help (flags and options)"
	@echo "  make version        - Show application version"
	@echo "  make help           - Show this help message"
	@echo ""
