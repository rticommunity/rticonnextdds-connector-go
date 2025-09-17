#!/bin/bash

# RTI Connext DDS Connector Library Downloader
# Downloads the latest (or specified version) of RTI Connext libraries
# from GitHub releases and extracts them to the correct location.

set -e

# Configuration
REPO_OWNER="rticommunity"
REPO_NAME="rticonnextdds-connector"
BASE_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}"
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download"
LIB_DIR="rticonnextdds-connector"
TEMP_DIR=$(mktemp -d)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cleanup function
cleanup() {
    if [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}
trap cleanup EXIT

# Print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}$1${NC}"
}

# Show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Downloads RTI Connext DDS Connector libraries from GitHub releases.

OPTIONS:
    -v, --version VERSION   Download specific version (e.g., v1.3.1)
    -l, --list              List available versions
    -c, --current           Show current installed version
    -f, --force             Force download even if version exists
    -h, --help              Show this help message

EXAMPLES:
    $0                      # Download latest version
    $0 -v v1.3.1           # Download specific version
    $0 -l                  # List available versions
    $0 -c                  # Show current version

EOF
}

# Detect platform and architecture
detect_platform() {
    local os=$(uname -s)
    local arch=$(uname -m)
    
    case "$os" in
        "Linux")
            case "$arch" in
                "x86_64") echo "linux-x64" ;;
                "aarch64"|"arm64") echo "linux-arm64" ;;
                *) print_error "Unsupported Linux architecture: $arch"; exit 1 ;;
            esac
            ;;
        "Darwin")
            case "$arch" in
                "x86_64") echo "osx-x64" ;;
                "arm64") echo "osx-arm64" ;;
                *) print_error "Unsupported macOS architecture: $arch"; exit 1 ;;
            esac
            ;;
        "MINGW"*|"CYGWIN"*|"MSYS"*)
            echo "win-x64"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Get available versions from GitHub API
get_versions() {
    print_status "Fetching available versions..."
    curl -s "$BASE_URL/releases" | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/' | head -10
}

# Get latest version
get_latest_version() {
    curl -s "$BASE_URL/releases/latest" | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

# Show current installed version
show_current_version() {
    if [ -d "$LIB_DIR" ]; then
        print_header "Current Installation:"
        
        # Try to find version info in libraries
        local platform=$(detect_platform)
        local lib_path="$LIB_DIR/lib/$platform"
        
        if [ -d "$lib_path" ]; then
            echo "  Platform: $platform"
            echo "  Library path: $lib_path"
            
            # Check if we can extract version from library files
            case "$platform" in
                linux-*)
                    if [ -f "$lib_path/librtiddsconnector.so" ]; then
                        local version_info=$(strings "$lib_path/librtiddsconnector.so" | grep -i "BUILD\|VERSION" | head -3 2>/dev/null || echo "Version info not found")
                        echo "  Version info: $version_info"
                    fi
                    ;;
                osx-*)
                    if [ -f "$lib_path/librtiddsconnector.dylib" ]; then
                        local version_info=$(strings "$lib_path/librtiddsconnector.dylib" | grep -i "BUILD\|VERSION" | head -3 2>/dev/null || echo "Version info not found")
                        echo "  Version info: $version_info"
                    fi
                    ;;
                win-*)
                    if [ -f "$lib_path/rtiddsconnector.dll" ]; then
                        echo "  Version info: Available (use findstr BUILD on Windows)"
                    fi
                    ;;
            esac
            
            # Show library files
            echo "  Libraries:"
            ls -la "$lib_path" | grep -E '\.(so|dylib|dll)$' | awk '{print "    " $9 " (" $5 " bytes)"}'
        else
            print_warning "No libraries found for platform: $platform"
        fi
    else
        print_warning "No RTI Connector libraries installed"
    fi
}

# Download and extract libraries
download_version() {
    local version="$1"
    local force="$2"
    local platform=$(detect_platform)
    
    print_header "RTI Connext DDS Connector Library Download"
    echo "Version: $version"
    echo "Platform: $platform"
    echo ""
    
    # Check if version already exists
    if [ -d "$LIB_DIR" ] && [ "$force" != "true" ]; then
        print_warning "Libraries already exist. Use --force to overwrite or --current to see current version."
        echo "Current installation:"
        show_current_version
        read -p "Continue with download? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Download cancelled."
            exit 0
        fi
    fi
    
    # Construct download URL - newer releases contain connectorlibs-<version>.zip
    local archive_name="connectorlibs-${version#v}.zip"
    local download_url="$DOWNLOAD_URL/$version/$archive_name"
    
    print_status "Downloading $archive_name..."
    print_status "URL: $download_url"
    
    # Download to temp directory
    cd "$TEMP_DIR"
    if ! curl -L -f -o "$archive_name" "$download_url"; then
        print_error "Failed to download $archive_name"
        print_error "Please check that version $version exists at:"
        print_error "https://github.com/$REPO_OWNER/$REPO_NAME/releases"
        exit 1
    fi
    
    print_status "Download completed. Extracting..."
    
    # Extract archive
    if command -v unzip >/dev/null 2>&1; then
        unzip -q "$archive_name"
    else
        print_error "unzip command not found. Please install unzip."
        exit 1
    fi
    
    # Find the extracted directory (could be connectorlibs-* or RTI-connector-*)
    local extracted_dir=$(find . -name "connectorlibs-*" -o -name "RTI-connector-*" -type d | head -1)
    if [ -z "$extracted_dir" ]; then
        print_error "Could not find extracted directory"
        exit 1
    fi
    
    print_status "Moving libraries to project directory..."
    
    # Go back to project root
    cd - > /dev/null
    
    # Remove existing libraries if they exist
    if [ -d "$LIB_DIR" ]; then
        print_status "Removing existing libraries..."
        rm -rf "$LIB_DIR"
    fi
    
    # Move extracted content
    mv "$TEMP_DIR/$extracted_dir" "$LIB_DIR"
    
    print_status "Libraries successfully installed!"
    print_status "Platform-specific libraries are available in: $LIB_DIR/lib/$platform"
    
    # Show what was installed
    echo ""
    show_current_version
    
    # Show library path setup instructions
    echo ""
    print_header "Setup Instructions:"
    case "$platform" in
        linux-*)
            echo "Add to your environment:"
            echo "export LD_LIBRARY_PATH=\$(pwd)/$LIB_DIR/lib/$platform:\$LD_LIBRARY_PATH"
            ;;
        osx-*)
            echo "Add to your environment:"
            echo "export DYLD_LIBRARY_PATH=\$(pwd)/$LIB_DIR/lib/$platform:\$DYLD_LIBRARY_PATH"
            ;;
        win-*)
            echo "Add to your PATH:"
            echo "set PATH=%CD%\\$LIB_DIR\\lib\\$platform;%PATH%"
            ;;
    esac
    
    print_status "You can now run: make test-local"
}

# Main script logic
main() {
    local version=""
    local force="false"
    local list_versions="false"
    local show_current="false"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
                shift 2
                ;;
            -f|--force)
                force="true"
                shift
                ;;
            -l|--list)
                list_versions="true"
                shift
                ;;
            -c|--current)
                show_current="true"
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Check if required tools are available
    for tool in curl unzip; do
        if ! command -v $tool >/dev/null 2>&1; then
            print_error "$tool is required but not installed."
            exit 1
        fi
    done
    
    # Handle different modes
    if [ "$list_versions" = "true" ]; then
        print_header "Available Versions:"
        get_versions
        exit 0
    fi
    
    if [ "$show_current" = "true" ]; then
        show_current_version
        exit 0
    fi
    
    # Get version to download
    if [ -z "$version" ]; then
        print_status "No version specified, fetching latest..."
        version=$(get_latest_version)
        if [ -z "$version" ]; then
            print_error "Could not determine latest version"
            exit 1
        fi
        print_status "Latest version: $version"
    fi
    
    # Download the version
    download_version "$version" "$force"
}

# Run main function
main "$@"