.PHONY: build clean install test release help

APP_NAME := resolve-transcoder
BUILD_DIR := build

# Get the latest Git tag as the version, fallback to "v0.0.0" if no tags exist
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Go build flags (inject version from Git tag)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(APP_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Build for all major platforms
build-all: clean
	@echo "Building for all platforms version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)

	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .

	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .

	@echo "Cross-platform builds complete!"
	@ls -la $(BUILD_DIR)/

# Install to system PATH (Unix-like systems)
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	go clean

# Create release packages
release: build-all
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/release

	# Windows packages
	zip -j $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe README.md LICENSE
	zip -j $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-windows-arm64.zip $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe README.md LICENSE

	# macOS packages
	tar -czf $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(APP_NAME)-darwin-amd64 -C .. README.md LICENSE
	tar -czf $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(APP_NAME)-darwin-arm64 -C .. README.md LICENSE

	# Linux packages
	tar -czf $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(APP_NAME)-linux-amd64 -C .. README.md LICENSE
	tar -czf $(BUILD_DIR)/release/$(APP_NAME)-$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR) $(APP_NAME)-linux-arm64 -C .. README.md LICENSE

	@echo "Release packages created in $(BUILD_DIR)/release/"
	@ls -la $(BUILD_DIR)/release/

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform (version from latest Git tag)"
	@echo "  build-all  - Build for all major platforms"
	@echo "  install    - Install to system PATH (requires sudo)"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  release    - Create release packages"
	@echo "  help       - Show this help message"

