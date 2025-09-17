# Go Get User Guide

This guide explains how to use RTI Connector when you install it via `go get`.

## The Challenge

When users run `go get github.com/rticommunity/rticonnextdds-connector-go`, they get the Go package but not the RTI Connector C libraries that are required for the CGO bindings to work. This creates a bootstrapping problem.

## The Solution

We provide multiple approaches to handle this scenario:

### 1. Go-based Library Downloader (Recommended)

A Go program that users can run directly from the installed package:

```bash
# Download latest libraries
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest

# Download specific version
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -version v1.3.1

# Check current installation
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -current

# List available versions
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -list
```

**Benefits:**
- ✅ Works immediately after `go get`
- ✅ No additional dependencies
- ✅ Cross-platform compatible
- ✅ Automatic platform detection
- ✅ Version management

### 2. Repository Cloning (For Development)

Users who want to contribute or access development tools:

```bash
git clone https://github.com/rticommunity/rticonnextdds-connector-go.git
cd rticonnextdds-connector-go
make download-libs
```

## Complete Workflow for Go Get Users

### Step 1: Install Package
```bash
go mod init my-rti-project
go get github.com/rticommunity/rticonnextdds-connector-go
```

### Step 2: Download Libraries
```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest
```

### Step 3: Set Library Path
The downloader will show platform-specific instructions:

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

### Step 4: Create Your Application
See the [go-get-example](../examples/go-get-example/) for a complete working example.

## Architecture Details

### Library Downloader Implementation

The downloader (`cmd/download-libs/main.go`) provides:

1. **GitHub API Integration**: Fetches latest release information
2. **Platform Detection**: Automatically determines OS/architecture
3. **Archive Handling**: Downloads and extracts ZIP files
4. **Directory Management**: Organizes libraries in expected structure
5. **Status Reporting**: Shows current installation and version info

### Key Features

- **No Dependencies**: Uses only Go standard library
- **Security**: Validates download URLs and prevents path traversal
- **Error Handling**: Comprehensive error messages and troubleshooting
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Version Management**: Can download specific versions or latest

### Library Organization

After download, libraries are organized as:

```
your-project/
├── rticonnextdds-connector/
│   ├── lib/
│   │   ├── linux-x64/          # Intel/AMD Linux
│   │   ├── linux-arm64/        # ARM Linux
│   │   ├── osx-x64/           # Intel macOS
│   │   ├── osx-arm64/         # Apple Silicon macOS
│   │   └── win-x64/           # Windows x64
│   ├── include/               # Header files
│   └── examples/              # C/Lua examples
├── main.go
└── go.mod
```

## Best Practices

### 1. Library Path Management

Create a setup script for your project:

```bash
#!/bin/bash
# setup.sh
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH  # macOS
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH    # Linux
```

### 2. Version Pinning

For reproducible builds, pin to specific versions:

```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -version v1.3.1
```

### 3. CI/CD Integration

In CI environments:

```bash
# Download libraries non-interactively
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -force

# Set library path
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH

# Run tests
go test ./...
```

### 4. Docker Integration

For Docker builds:

```dockerfile
FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Download RTI libraries
RUN go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest

COPY . .
RUN go build -o myapp ./cmd/myapp

# Set library path
ENV LD_LIBRARY_PATH=/app/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH

CMD ["./myapp"]
```

## Troubleshooting

### Common Issues

1. **Libraries not found at runtime**
   - Ensure library path is set correctly
   - Verify libraries were downloaded to expected location

2. **Wrong architecture libraries**
   - The downloader auto-detects platform
   - Check with: `go run .../download-libs@latest -current`

3. **Network/download issues**
   - Check internet connectivity
   - Verify GitHub access (no corporate firewall blocking)

4. **Version not found**
   - List available versions: `go run .../download-libs@latest -list`
   - Check RTI releases: https://github.com/rticommunity/rticonnextdds-connector/releases

### Debug Commands

```bash
# Check current status
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -current

# Verify library files
ls -la rticonnextdds-connector/lib/*/

# Test library loading
go run your-program.go
```

## Migration from Other Installation Methods

### From Git Clone

If you previously cloned the repository but want to use `go get`:

1. Remove the cloned directory from GOPATH
2. Use `go get` in your project
3. Download libraries with the Go downloader
4. Update import paths if needed

### From Manual Library Management

If you manually managed libraries:

1. Remove old library directories
2. Use the Go downloader for consistent organization
3. Update library paths to use new structure

## Future Enhancements

Potential improvements being considered:

1. **Automatic Library Management**: Libraries downloaded automatically on first import
2. **Version Synchronization**: Automatic matching of Go package and library versions
3. **Local Caching**: Shared library cache across projects
4. **Dependency Detection**: Automatic library path configuration