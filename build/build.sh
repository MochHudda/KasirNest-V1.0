#!/bin/bash

# KasirNest Build Script
# This script builds the KasirNest application with optimization and obfuscation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="kasirnest"
VERSION="1.0.0"
BUILD_DIR="build/output"
DIST_DIR="build/dist"

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.buildTime=$(date -u +%Y%m%d.%H%M%S)"
CGO_ENABLED=0

echo -e "${BLUE}=== KasirNest Build Script ===${NC}"
echo -e "${BLUE}Version: ${VERSION}${NC}"
echo -e "${BLUE}Build Time: $(date)${NC}"
echo ""

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}')
echo "Go version: $GO_VERSION"

# Check if garble is installed
if ! command -v garble &> /dev/null; then
    echo -e "${YELLOW}Installing garble for obfuscation...${NC}"
    go install mvdan.cc/garble@latest
fi

# Check if fyne is installed
if ! command -v fyne &> /dev/null; then
    echo -e "${YELLOW}Installing fyne build tools...${NC}"
    go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# Create build directories
echo -e "${YELLOW}Creating build directories...${NC}"
mkdir -p $BUILD_DIR
mkdir -p $DIST_DIR

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf $BUILD_DIR/*
rm -rf $DIST_DIR/*

# Download dependencies
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod download
go mod tidy

# Run tests (optional)
if [ "$1" != "--no-test" ]; then
    echo -e "${YELLOW}Running tests...${NC}"
    go test ./... -v
fi

# Build function
build_binary() {
    local platform=$1
    local arch=$2
    local extension=$3
    local output="${BUILD_DIR}/${APP_NAME}_${platform}_${arch}${extension}"
    
    echo -e "${YELLOW}Building for ${platform}/${arch}...${NC}"
    
    export GOOS=$platform
    export GOARCH=$arch
    export CGO_ENABLED=$CGO_ENABLED
    
    # Build with garble for obfuscation
    if command -v garble &> /dev/null && [ "$OBFUSCATE" = "true" ]; then
        echo -e "${BLUE}Building with obfuscation...${NC}"
        garble -literals -seed=random build -ldflags="$LDFLAGS" -o "$output" .
    else
        echo -e "${BLUE}Building without obfuscation...${NC}"
        go build -ldflags="$LDFLAGS" -o "$output" .
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built: $output${NC}"
        
        # Get file size
        if [ "$platform" = "windows" ]; then
            SIZE=$(stat -c%s "$output" 2>/dev/null || echo "unknown")
        else
            SIZE=$(stat -f%z "$output" 2>/dev/null || stat -c%s "$output" 2>/dev/null || echo "unknown")
        fi
        echo -e "${BLUE}  Size: $SIZE bytes${NC}"
    else
        echo -e "${RED}✗ Failed to build for ${platform}/${arch}${NC}"
        exit 1
    fi
}

# Parse command line arguments
OBFUSCATE="true"
BUILD_ALL="false"

while [[ $# -gt 0 ]]; do
    case $1 in
        --no-obfuscate)
            OBFUSCATE="false"
            shift
            ;;
        --all-platforms)
            BUILD_ALL="true"
            shift
            ;;
        --no-test)
            # Already handled above
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--no-obfuscate] [--all-platforms] [--no-test]"
            exit 1
            ;;
    esac
done

# Build for different platforms
echo -e "${YELLOW}Starting builds...${NC}"

if [ "$BUILD_ALL" = "true" ]; then
    # Build for all platforms
    build_binary "windows" "amd64" ".exe"
    build_binary "windows" "386" ".exe"
    build_binary "linux" "amd64" ""
    build_binary "linux" "386" ""
    build_binary "darwin" "amd64" ""
    build_binary "darwin" "arm64" ""
else
    # Build for current platform only
    CURRENT_OS=$(go env GOOS)
    CURRENT_ARCH=$(go env GOARCH)
    
    if [ "$CURRENT_OS" = "windows" ]; then
        build_binary "$CURRENT_OS" "$CURRENT_ARCH" ".exe"
    else
        build_binary "$CURRENT_OS" "$CURRENT_ARCH" ""
    fi
fi

# Create distribution packages
echo -e "${YELLOW}Creating distribution packages...${NC}"

for binary in $BUILD_DIR/*; do
    if [ -f "$binary" ]; then
        filename=$(basename "$binary")
        platform_arch=$(echo "$filename" | sed "s/${APP_NAME}_//")
        
        echo -e "${YELLOW}Packaging $filename...${NC}"
        
        # Create package directory
        package_dir="$DIST_DIR/${APP_NAME}_${platform_arch}"
        mkdir -p "$package_dir"
        
        # Copy binary
        cp "$binary" "$package_dir/"
        
        # Copy configuration template
        cp "config/app.ini.example" "$package_dir/"
        
        # Copy assets
        cp -r "assets" "$package_dir/" 2>/dev/null || true
        
        # Copy documentation
        cp "README.md" "$package_dir/" 2>/dev/null || true
        cp "FIREBASE_SETUP.md" "$package_dir/" 2>/dev/null || true
        
        # Create archive
        cd "$DIST_DIR"
        if [[ $filename == *"windows"* ]]; then
            # Create ZIP for Windows
            zip -r "${APP_NAME}_${platform_arch}.zip" "${APP_NAME}_${platform_arch}/"
        else
            # Create TAR.GZ for Unix-like systems
            tar -czf "${APP_NAME}_${platform_arch}.tar.gz" "${APP_NAME}_${platform_arch}/"
        fi
        cd - > /dev/null
        
        echo -e "${GREEN}✓ Package created: ${APP_NAME}_${platform_arch}${NC}"
    fi
done

# Build summary
echo ""
echo -e "${GREEN}=== Build Summary ===${NC}"
echo -e "${GREEN}Built binaries:${NC}"
ls -la $BUILD_DIR/

echo ""
echo -e "${GREEN}Distribution packages:${NC}"
ls -la $DIST_DIR/*.{zip,tar.gz} 2>/dev/null || echo "No packages found"

echo ""
echo -e "${GREEN}✓ Build completed successfully!${NC}"

# Optional: Run the binary for current platform
if [ "$1" = "--run" ]; then
    echo -e "${YELLOW}Starting application...${NC}"
    CURRENT_OS=$(go env GOOS)
    CURRENT_ARCH=$(go env GOARCH)
    
    if [ "$CURRENT_OS" = "windows" ]; then
        $BUILD_DIR/${APP_NAME}_${CURRENT_OS}_${CURRENT_ARCH}.exe
    else
        $BUILD_DIR/${APP_NAME}_${CURRENT_OS}_${CURRENT_ARCH}
    fi
fi