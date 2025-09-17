# RTI Connector Library Management

This document explains how to manage RTI Connext DDS Connector libraries in the Go project.

## Overview

The RTI Connector Go binding requires native C libraries from RTI Connext DDS. This project provides automated tools to download and manage these libraries from the official RTI GitHub releases.

## Quick Start

### Download Latest Libraries

```bash
# Using Make (recommended)
make download-libs

# Or directly using Go command
go run ./cmd/download-libs
```

### Check Current Installation

```bash
make check-libs
```

### List Available Versions

```bash
make list-lib-versions
```

## Library Download Tool

The Go-based download tool (`cmd/download-libs`) provides comprehensive library management:

### Basic Usage

```bash
# Download latest version
go run ./cmd/download-libs

# Download specific version  
go run ./cmd/download-libs -version v1.3.1

# Force download (overwrite existing)
go run ./cmd/download-libs -force

# Show current installation
go run ./cmd/download-libs -current

# List available versions
go run ./cmd/download-libs -list

# Use make commands (recommended)
make download-libs          # Download latest
make download-libs-latest   # Force download latest
make check-libs            # Show current installation
make list-lib-versions     # List available versions
```

### Platform Support

The tool automatically detects your platform and downloads the appropriate libraries:

- **Linux x64**: `linux-x64` libraries
- **Linux ARM64**: `linux-arm64` libraries  
- **macOS Intel**: `osx-x64` libraries
- **macOS Apple Silicon**: `osx-arm64` libraries
- **Windows x64**: `win-x64` libraries

### Library Path Setup

After downloading, the tool shows the appropriate environment setup:

**Linux:**
```bash
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH
```

**macOS:**
```bash
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH
```

**Windows:**
```cmd
set PATH=%CD%\rticonnextdds-connector\lib\win-x64;%PATH%
```

## Makefile Integration

The following Make targets are available:

| Target | Description |
|--------|-------------|
| `make download-libs` | Download latest libraries (interactive) |
| `make download-libs-latest` | Force download latest libraries |
| `make check-libs` | Show current installation info |
| `make list-lib-versions` | List available versions |

## Library Sources

Libraries are downloaded from the official RTI repository:
- **Repository**: https://github.com/rticommunity/rticonnextdds-connector
- **Releases**: https://github.com/rticommunity/rticonnextdds-connector/releases

## Version Management

### Current Version Detection

The tool can detect the currently installed version by examining the library files:

```bash
go run ./cmd/download-libs -current
# or
make check-libs
```

This shows:
- Platform and architecture
- Library path
- Version information extracted from binaries
- List of installed library files

### Available Versions

Check what versions are available for download:

```bash
go run ./cmd/download-libs -list
# or
make list-lib-versions
```

Recent versions include:
- v1.3.1 (latest)
- v1.3.0
- v1.2.3
- v1.2.2
- v1.2.0

## Directory Structure

After downloading, libraries are organized as:

```
rticonnextdds-connector/
├── lib/
│   ├── linux-x64/
│   │   ├── libnddsc.so
│   │   ├── libnddscore.so
│   │   └── librtiddsconnector.so
│   ├── linux-arm64/
│   ├── osx-x64/
│   ├── osx-arm64/
│   │   ├── libnddsc.dylib
│   │   ├── libnddscore.dylib
│   │   └── librtiddsconnector.dylib
│   └── win-x64/
├── include/
└── examples/
```

## Troubleshooting

### Common Issues

1. **Libraries not found**: Ensure library path is set correctly
2. **Go tool issues**: Ensure Go is properly installed and accessible
3. **Network issues**: Check internet connection and GitHub access
4. **Version not found**: Verify version exists in releases

### Debug Commands

```bash
# Check if libraries are in path
echo $LD_LIBRARY_PATH    # Linux
echo $DYLD_LIBRARY_PATH  # macOS

# Verify library files
ls -la rticonnextdds-connector/lib/$(uname -s | tr '[:upper:]' '[:lower:]')-*

# Test library loading
make test-local
```

## Integration with CI/CD

For automated builds, you can use:

```bash
# In CI scripts
go run ./cmd/download-libs -force  # Force download latest
make test-local                   # Run tests
```

The tool is designed to work in both interactive and automated environments.

## Contributing

When contributing to the project:

1. Always test with multiple library versions
2. Ensure the download tool works on all supported platforms
3. Update documentation if adding new library management features
4. Test both Go module and manual installation paths