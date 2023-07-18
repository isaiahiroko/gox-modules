# Get the version from the command line argument (default: v1.0.0)
VERSION ?= 1.0.0

# Set the name of your Go application
APP_NAME = makr

# Set the directories
SRC_DIR = ./
OUT_DIR = ./dist

# Set the platforms and corresponding output filenames
PLATFORMS = linux darwin windows

# Set the output filenames for each platform
OUT_LINUX = $(OUT_DIR)/$(VERSION)/$(APP_NAME)-linux
OUT_MACOS = $(OUT_DIR)/$(VERSION)/$(APP_NAME)-macos
OUT_WINDOWS = $(OUT_DIR)/$(VERSION)/$(APP_NAME)-windows.exe

# Default target
default: build

# Build the Go application for the current platform
build:
	@echo "Building $(APP_NAME) v$(VERSION) for the current platform..."
	go build -o $(OUT_DIR)/$(VERSION)/$(APP_NAME) $(SRC_DIR)

# Build the Go application for all supported platforms
build-all: clean
	@echo "Building $(APP_NAME) v$(VERSION) for all platforms..."
	mkdir -p $(OUT_DIR)/$(VERSION)
	for platform in $(PLATFORMS); do \
		out_file=$(OUT_DIR)/$(VERSION)/$(APP_NAME)-$(VERSION)-$$platform; \
		GOOS=$$platform GOARCH=amd64 go build -o $$out_file $(SRC_DIR); \
		if [ $$platform = "windows" ]; then \
			mv $$out_file $$out_file.exe; \
		fi; \
	done

# Clean the compiled binaries
clean:
	@echo "Cleaning..."
	rm -rf $(OUT_DIR)/$(VERSION)

.PHONY: build build-all clean
