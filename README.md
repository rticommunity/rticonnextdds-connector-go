# RTI Connector for Connext DDS - Go

[![Coverage](https://img.shields.io/badge/coverage-90.4%25-brightgreen)](https://github.com/rticommunity/rticonnextdds-connector-go/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rticommunity/rticonnextdds-connector-go)](https://goreportcard.com/report/github.com/rticommunity/rticonnextdds-connector-go)
[![Build and Test](https://github.com/rticommunity/rticonnextdds-connector-go/actions/workflows/build.yml/badge.svg)](https://github.com/rticommunity/rticonnextdds-connector-go/actions/workflows/build.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rticommunity/rticonnextdds-connector-go)](https://github.com/rticommunity/rticonnextdds-connector-go/blob/master/go.mod)

> A lightweight, easy-to-use Go binding for RTI Connext DDS that enables rapid development of distributed applications.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage Examples](#usage-examples)
- [Platform Support](#platform-support)
- [Development](#development)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Support](#support)

## Overview

**RTI Connector** for Connext DDS is a lightweight, easy-to-use API that provides access to the power and functionality of [RTI Connext DDS](http://www.rti.com/products/index.html). Built on [XML-Based Application Creation](https://community.rti.com/static/documentation/connext-dds/6.0.0/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf) and Dynamic Data, it enables rapid prototyping and development of distributed applications.

Originally created by the RTI Research Group for demos and proof-of-concepts, RTI Connector is perfect for developers who need to quickly integrate DDS communication into their Go applications without the complexity of the full RTI Connext DDS Professional SDK.

## Key Features

✨ **Simple API** - Easy-to-use Go interface that hides DDS complexity  
🚀 **Rapid Development** - Get up and running with DDS in minutes  
🔄 **Dynamic Data** - No need to generate code from IDL files  
📄 **XML Configuration** - Define data types and QoS policies in XML  
🌐 **Cross-Platform** - Supports Linux x64, macOS, and Windows x64  

## Quick Start

### Prerequisites

- Go 1.21 or later
- Supported platform (Linux x64, macOS, Windows x64)

### Installation

1. **Get the Go package:**
```bash
go get github.com/rticommunity/rticonnextdds-connector-go
```

2. **Download RTI Connector libraries:**
```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest
```

3. **Set library path:**
The download tool will show you the exact command for your platform. For example:
```bash
# macOS
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH

# Linux  
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH

# Windows (PowerShell)
$env:PATH = "$(pwd)\rticonnextdds-connector\lib\win-x64;$env:PATH"
```

> **💡 New to RTI Connector Go?** Try the [go-get-example](examples/go-get-example/) first - it provides a complete walkthrough of this installation process with a simple working example.

### Your First DDS Application

1. **Create an XML configuration file** (`ShapeExample.xml`):

```xml
<dds>
  <qos_library name="QosLibrary">
    <qos_profile name="DefaultProfile" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true"/>
  </qos_library>
  
  <types>
    <struct name="ShapeType">
      <member name="color" type="string" key="true"/>
      <member name="x" type="long"/>
      <member name="y" type="long"/>
      <member name="shapesize" type="long"/>
    </struct>
  </types>
  
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
      <subscriber name="MySubscriber">
        <data_reader name="MySquareReader" topic_ref="Square"/>
      </subscriber>
    </domain_participant>
  </domain_participant_library>
</dds>
```

2. **Create a publisher** (`publisher.go`):

```go
package main

import (
    "log"
    "time"
    
    rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
    // Create connector
    connector, err := rti.NewConnector("MyParticipantLibrary::Zero", "./ShapeExample.xml")
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

3. **Create a subscriber** (`subscriber.go`):

```go
package main

import (
    "log"
    
    rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
    // Create connector
    connector, err := rti.NewConnector("MyParticipantLibrary::Zero", "./ShapeExample.xml")
    if err != nil {
        log.Fatal(err)
    }
    defer connector.Delete()

    // Get input (reader)
    input, err := connector.GetInput("MySubscriber::MySquareReader")
    if err != nil {
        log.Fatal(err)
    }

    // Read data
    log.Println("Waiting for data...")
    for {
        connector.Wait(-1) // Wait indefinitely for data
        input.Take()
        
        numSamples, _ := input.Samples.GetLength()
        for i := 0; i < numSamples; i++ {
            if valid, _ := input.Infos.IsValid(i); valid {
                color, _ := input.Samples.GetString(i, "color")
                x, _ := input.Samples.GetInt(i, "x")
                y, _ := input.Samples.GetInt(i, "y")
                shapesize, _ := input.Samples.GetInt(i, "shapesize")
                
                log.Printf("Received: color=%s, x=%d, y=%d, size=%d", color, x, y, shapesize)
            }
        }
    }
}
```

4. **Run your application:**

Make sure you've completed the installation steps above (including library download and path setup), then:

```bash
# Terminal 1 - Start subscriber
go run subscriber.go

# Terminal 2 - Start publisher  
go run publisher.go
```

You should see the subscriber receiving data published by the publisher!

## Library Management

The installation above uses our automated library download tool. For advanced scenarios, see our comprehensive guides:

- **[Library Management Documentation](docs/LIBRARY_MANAGEMENT.md)** - Complete guide to library installation options
- **[Go Get Users Guide](docs/GO_GET_USERS.md)** - Specific help for `go get` workflow
- **[Examples](examples/)** - Ready-to-run code samples

### Library Download Tool Options

```bash
# Download specific version
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -version v1.3.1

# Check what's currently installed
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -current

# List available versions
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -list
```

### Development Setup

For contributors or users who want to build from source:

```bash
git clone https://github.com/rticommunity/rticonnextdds-connector-go.git
cd rticonnextdds-connector-go
make download-libs
make test-local
```

See our [contribution guide](CONTRIBUTING.md) for more details.
- ✅ Shows version information and setup instructions
- ✅ Validates downloads and provides helpful error messages

## Usage Examples

Explore our comprehensive examples to learn different patterns and use cases:

| Example | Description | Key Features |
|---------|-------------|--------------|
| [Simple](examples/simple/) | Basic publisher/subscriber | Getting started, basic data flow |
| [Go Get Example](examples/go-get-example/) | For `go get` users | Inline XML, library download workflow |
| [Shapes Demo](examples/array/) | Array data handling | Complex data types, arrays |
| [JSON Integration](examples/go_struct/) | Go struct mapping | JSON serialization, struct binding |
| [Request-Reply](examples/request_reply/) | RPC pattern | Synchronous communication |
| [Security](examples/security/) | Secure communication | Authentication, encryption |
| [Multiple Files](examples/module/) | Modular configuration | XML organization, reusability |

**📁 [Browse all examples →](examples/README.md)**

## Platform Support

RTI Connector supports the following platforms:

| Platform | Architecture | Status |
|----------|-------------|---------|
| **Linux** | x86_64 | ✅ Supported |
| **macOS** | x86_64 (Intel/Apple Silicon) | ✅ Supported |
| **Windows** | x86_64 | ✅ Supported |

> 📝 **Note**: If you need support for additional architectures, please contact your RTI account manager or [sales@rti.com](mailto:sales@rti.com).

### Version Information

To check the version of your installed RTI libraries:

```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest -current
```

## Development

### Building and Testing

For contributors and developers working with the source code:

```bash
# Clone and setup
git clone https://github.com/rticommunity/rticonnextdds-connector-go.git
cd rticonnextdds-connector-go
make download-libs

# Run all tests with coverage
make test-local

# Run comprehensive test suite with quality checks
./test_improvements.sh

# Build all packages
go build ./...

# Run static analysis
go vet ./...

# Format code
go fmt ./...
```

### Threading Model

> ⚠️ **Important**: The Connector Native API does not implement thread safety. You are responsible for protecting calls to Connector in multi-threaded applications.

The native code was originally designed for single-threaded environments (RTI Prototyper and Lua). While we've added support for Go, Python, and JavaScript, thread safety remains the developer's responsibility.

## Documentation

- 📖 **[Examples](examples/README.md)** - Comprehensive examples and tutorials  
- 🧪 **[Testing Guide](TESTING.md)** - Development and testing guidelines
- 📚 **[Library Management](docs/LIBRARY_MANAGEMENT.md)** - Managing RTI Connector libraries
- 🚀 **[Go Get Users Guide](docs/GO_GET_USERS.md)** - Complete guide for `go get` users
- 🤝 **[Contributing](CONTRIBUTING.md)** - How to contribute to the project

### Additional Resources

- [RTI Connext DDS Documentation](https://community.rti.com/documentation)
- [XML-Based Application Creation Guide](https://community.rti.com/static/documentation/connext-dds/6.0.0/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf)
- [RTI Community Forum](https://community.rti.com/forums/technical-questions)

## Contributing

We welcome contributions! Here's how to get started:

1. **📝 [Sign the CLA](CONTRIBUTING.md)** - Required for all contributions
2. **🍴 Fork and clone** the repository
3. **🔧 Make your changes** with tests
4. **✅ Run quality checks**: `make test-local` or `./test_improvements.sh`
5. **📤 Submit a pull request**

All contributions are automatically tested for quality, including:
- Build verification across platforms
- Code linting and formatting  
- Comprehensive test suite validation
- Coverage analysis

### Development Environment

```bash
# Clone the repository
git clone https://github.com/rticommunity/rticonnextdds-connector-go.git
cd rticonnextdds-connector-go

# Run tests to verify setup
./test_improvements.sh

# Start developing!
```

## Support

### Getting Help

- 💬 **[RTI Community Forum](https://community.rti.com/forums/technical-questions)** - Technical questions and discussions
- 🐛 **[GitHub Issues](https://github.com/rticommunity/rticonnextdds-connector-go/issues)** - Bug reports and feature requests
- 📧 **[Contact RTI](mailto:sales@rti.com)** - Commercial support and licensing

---

**Ready to get started?** Check out our [Quick Start](#quick-start) guide or explore the [examples](examples/README.md)!
