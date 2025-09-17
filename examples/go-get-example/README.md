# Getting Started with Go Get

This example shows the easiest way to get started with RTI Connector Go using `go get`.

## Quick Start

1. **Create a new project and initialize Go module:**
```bash
mkdir my-rti-demo
cd my-rti-demo
go mod init my-rti-demo
```

2. **Get the package:**
```bash
go get github.com/rticommunity/rticonnextdds-connector-go
```

3. **Download RTI libraries:**
```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest
```

4. **Set library path (for runtime):**
```bash
# macOS (Apple Silicon/ARM64)
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH

# macOS (Intel/x86_64)
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-x64:$DYLD_LIBRARY_PATH

# Linux  
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH

# Windows (PowerShell)
$env:PATH = "$(pwd)\rticonnextdds-connector\lib\win-x64;$env:PATH"
```

> **💡 macOS Users**: Use `osx-arm64` for Apple Silicon Macs (M1/M2/M3) and `osx-x64` for Intel Macs. You can check your architecture with `uname -m` (arm64 = Apple Silicon, x86_64 = Intel).

> **✨ New**: CGO compilation is now automatic! No manual CGO flags needed.

5. **Run the example:**
```bash
go run publisher.go
```

You should see:
```
Creating RTI Connector...
✅ RTI Connector created successfully!
✅ Successfully published test message!
RTI Connector Go is working with libraries downloaded via go get workflow!
```

## The Code

The `publisher.go` demonstrates the minimal working example:

- **Inline XML configuration** using `str://"<xml>..."` syntax (no external files needed)
- **Simple connector creation** and basic publishing
- **Clean resource management** with `defer connector.Delete()`

This example proves your RTI Connector Go setup is working correctly.

## Next Steps

Once you've verified this works:

1. **Explore other examples** in the parent directory for more advanced features
2. **Check the main README** for comprehensive documentation
3. **Read the library management docs** at `docs/GO_GET_USERS.md` for troubleshooting

## Need Help?

- **Library issues?** See `docs/LIBRARY_MANAGEMENT.md`
- **Go get specific problems?** See `docs/GO_GET_USERS.md`  
- **General usage?** Check the main project README

This getting started guide focuses on the essential workflow - once it works, you're ready for any RTI Connector Go project!
  
  <domain_library name="MyDomainLibrary">
    <domain name="MyDomain" domain_id="0">
      <register_type name="ShapeType" type_ref="ShapeType"/>
      <topic name="Square" register_type_ref="ShapeType"/>
    </domain>
  </domain_library>
  
  <domain_participant_library name="MyParticipantLibrary">
    <domain_participant name="Zero" domain_ref="MyDomainLibrary::MyDomain">
      <publisher name="MyPublisher">
        <data_writer name="MySquareWriter" topic_ref="Square"/>
      </publisher>
    </domain_participant>
  </domain_participant_library>
</dds>"`

    // Create connector from XML string
    connector, err := rti.NewConnector("MyParticipantLibrary::Zero", xmlConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer connector.Delete()

    // Get output (writer)
    output, err := connector.GetOutput("MyPublisher::MySquareWriter")
    if err != nil {
        log.Fatal(err)
    }

    // Publish data
    for i := 0; i < 10; i++ {
        output.Instance.SetString("color", "BLUE")
        output.Instance.SetInt("x", i*10)
        output.Instance.SetInt("y", i*20)
        output.Instance.SetInt("shapesize", 30)
        
        output.Write()
        log.Printf("Published sample %d", i)
        time.Sleep(time.Second)
    }
}
```

## Run the Example

```bash
go run publisher.go
```

## Library Management Commands

```bash
# Check current installation
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -current

# Update to latest libraries
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -force

# Download specific version
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -version v1.3.1

# List available versions
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -list
```

## Directory Structure After Setup

```
your-project/
├── publisher.go
├── rticonnextdds-connector/
│   ├── lib/
│   │   ├── linux-x64/          # Linux libraries
│   │   ├── linux-arm64/        # Linux ARM libraries  
│   │   ├── osx-x64/           # macOS Intel libraries
│   │   ├── osx-arm64/         # macOS Apple Silicon libraries
│   │   └── win-x64/           # Windows libraries
│   ├── include/
│   └── examples/
└── go.mod
```

The library downloader automatically detects your platform and downloads the appropriate libraries.