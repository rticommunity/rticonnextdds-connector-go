rticonnextdds-connector-go
=======
[![Coverage](https://codecov.io/gh/rticommunity/rticonnextdds-connector-go/branch/master/graph/badge.svg)](https://codecov.io/gh/rticommunity/rticonnextdds-connector-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rticommunity/rticonnextdds-connector-go)](https://goreportcard.com/report/github.com/rticommunity/rticonnextdds-connector-go)
[![Build and Test](https://github.com/rticommunity/rticonnextdds-connector-go/actions/workflows/build.yml/badge.svg)](https://github.com/rticommunity/rticonnextdds-connector-go/actions/workflows/build.yml)


### RTI Connector for Connext DDS
*RTI Connector* for Connext DDS is a quick and easy way to access the power and
functionality of [RTI Connext DDS](http://www.rti.com/products/index.html).
It is based on [XML-Based Application Creation](https://community.rti.com/static/documentation/connext-dds/6.0.0/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf) and Dynamic Data.

*Connector* was created by the RTI Research Group to quickly and easily develop demos
and proofs of concept. It can be useful for anybody that needs
a quick way to develop an application communicating over the Connext DDS Databus.
Thanks to the binding with multiple programming languages, you can integrate
with many other available technologies.

The *Connector* library is provided in binary form for [select architectures](https://github.com/rticommunity/rticonnextdds-connector/tree/master/lib). Language bindings and examples are provided in source format.

Go *Connector* leverages [cgo](https://golang.org/cmd/cgo) to call its C library;
this detail is hidden in a Go wrapper. 

### Examples
#### Simple Reader
```golang
package main
import (
	rti "github.com/rticommunity/rticonnextdds-connector-go"
	"log"
)

func main() {
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", "./ShapeExample.xml")
	if err != nil {
		log.Panic(err)
	}
	defer connector.Delete()
	input, err := connector.GetInput("MySubscriber::MySquareReader")
	if err != nil {
		log.Panic(err)
	}

	for {
		connector.Wait(-1)
		input.Take()
		numOfSamples, _ := input.Samples.GetLength()
		for j := 0; j < numOfSamples; j++ {
			valid, _ := input.Infos.IsValid(j)
			if valid {
				json, err := input.Samples.GetJSON(j)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(string(json))
				}
			}
		}
	}
}
```

#### Simple Writer
```golang
package main
import (
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
	"time"
)

func main() {
	// Create a connector defined in the XML configuration
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", "./ShapeExample.xml")
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer connector.Delete()

	// Get an output from the connector
	output, err := connector.GetOutput("MyPublisher::MySquareWriter")
	if err != nil {
		log.Panic(err)
	}

	// Set values to the instance and write the instance
	for i := 0; i < 10; i++ {
		output.Instance.SetInt("x", i)
		output.Instance.SetInt("y", i*2)
		output.Instance.SetInt("shapesize", 30)
		output.Instance.SetString("color", "BLUE")
		output.Write()
		log.Println("Writing...")
		time.Sleep(time.Second * 1)
	}
}
```

#### XML Configurations
```xml
<dds>
  <!-- Qos Library -->
  <qos_library name="QosLibrary">
    <qos_profile name="DefaultProfile" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true">
      <participant_qos>
        <transport_builtin>
          <mask>UDPV4 | SHMEM</mask>
        </transport_builtin>
      </participant_qos>
    </qos_profile>
  </qos_library>
  <!-- types -->
  <types>
    <struct name="ShapeType" extensibility="extensible">
      <member name="color" stringMaxLength="128" id="0" type="string" key="true"/>
      <member name="x" id="1" type="long"/>
      <member name="y" id="2" type="long"/>
      <member name="shapesize" id="3" type="long"/>
    </struct>
    <enum name="ShapeFillKind" extensibility="extensible">
      <enumerator name="SOLID_FILL" value="0"/>
      <enumerator name="TRANSPARENT_FILL" value="1"/>
      <enumerator name="HORIZONTAL_HATCH_FILL" value="2"/>
      <enumerator name="VERTICAL_HATCH_FILL" value="3"/>
    </enum>
    <struct name="ShapeTypeExtended" baseType="ShapeType" extensibility="extensible">
      <member name="fillKind" id="4" type="nonBasic" nonBasicTypeName="ShapeFillKind"/>
      <member name="angle" id="5" type="float"/>
    </struct>
  </types>
  <!-- Domain Library -->
  <domain_library name="MyDomainLibrary">
    <domain name="MyDomain" domain_id="0">
      <register_type name="ShapeType" type_ref="ShapeType"/>
      <topic name="Square" register_type_ref="ShapeType"/>
    </domain>
  </domain_library>
  <!-- Participant library -->
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
Please see [examples](examples/README.md) for usage details.

### Getting started
#### Using Go Modules
Be sure you have golang installed (we tested with golang v1.17). 

Import:

```golang
import "github.com/rticommunity/rticonnextdds-connector-go"
```

Build:
```bash
$ go build reader.go
```

A dependency to the latest stable version of rticonnextdds-connector-go should be automatically added to your `go.mod` file.

Run:

To run your application, you need to add the Connector C library to your library path.
```bash
$ export LD_LIBRARY_PATH=$GOPATH//go/pkg/mod/github.com/rticommunity/rticonnextdds-connector-go\@{version}-{YYYYMMDDHHmm}-{commit_id}/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH
$ ./simple_reader
```

### Platform support
Go *Connector* builds its library for few [select architectures](https://github.com/rticommunity/rticonnextdds-connector/tree/master/lib). If you need another architecture, please contact your RTI account manager or sales@rti.com.

If you want to check the version of the libraries you can run the following command:

``` bash
strings ./rticonnextdds-connector/lib/linux-x64/librtiddsconnector.so | grep BUILD
```

### Threading model
The *Connector* Native API does not yet implement any mechanism for thread safety. Originally, the *Connector* native code was built to work with *RTI Prototyper* and Lua. That was a single-threaded loop. RTI then introduced support for JavaScript, Python, and Go. For now, you are responsible for protecting calls to *Connector*. Thread safety
may be implemented in the future.

### Support
*Connector* is an experimental RTI product. If you have questions, please use the [RTI Community Forum](https://community.rti.com/forums/technical-questions). If you would like to report a bug or have a feature request, please create an [issue](https://github.com/rticommunity/rticonnextdds-connector-go/issues).

### Documentation
The best way to get started with *Connector* is to look at the
examples; you will see that it is very easy to use.

### Contributing
Contributions to the code, examples, documentation are really appreciated. Please follow the steps below for contributions. 

1. [Sign the CLA](CONTRIBUTING.md).
1. Create a fork and make your changes.
1. Run tests and linters (make test lint).
1. Push your branch.
1. Open a new [pull request](https://github.com/rticommunity/rticonnextdds-connector-go/compare).
