RTI Connext Go Connector Examples
========

### Installation and Platform support
Check [here](https://github.com/rticommunity/rticonnextdds-connector#getting-started-with-python) and [here](https://github.com/rticommunity/rticonnextdds-connector#platform-support).
If you still have trouble write on the [RTI Community Forum](https://community.rti.com/forums/technical-questions)

### Available examples
In this directory you can find the following examples

 * **simple**: shows simple examples of how to write samples and how to read/take

### Building and running examples
``` bash
$ go build reader.go
```
Currently, Go Connector links to DDS library dynamically so the path for DDS library needs to be added to a shared library path (e.g. LD_LIBRARY_PATH for Linux). "ARCH" needs to be replaced with your architecture(e.g. x64Linux2.6gcc4.45 for 64-bit Linux)
``` bash
$ export LD_LIBRARY_PATH=$GOPATH/src/github.com/kyoungho/rticonnextdds-connector/lib/ARCH:$LD_LIBRARY_PATH
$ ./reader
```

### API Overview
#### Import the Connector library
If you want to use the RTI Connext Go Connector, you have to import the package.

```go
import "github.com/kyoungho/rticonnextdds-connector"
```

#### Instantiate a new connector
To create a new connector you have to pass a location of an XML configuration file and a configuration name in XML. For more information on
the XML format check the [XML App Creation guide](https://community.rti.com/rti-doc/510/ndds.5.1.0/doc/pdf/RTI_CoreLibrariesAndUtilities_XML_AppCreation_GettingStarted.pdf) or take a look at the [ShapeExample.xml](ShapeExample.xml) file included in this directory.  

```go
connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
```
#### Delete a connector
To destroy all the DDS entities created by a connector, you should call the ```Delete``` function.

```go
connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
...
...
connector.Delete()
```

#### Write a data sample
To write a data sample, you have to get a reference to the output port:

```go
output, err := connector.GetOutput("MyPublisher::MySquareWriter")
```

then you have to set the fields in a sample instance:

```go
output.Instance.SetInt("x", i)
output.Instance.SetInt("y", i*2)
output.Instance.SetInt("shapesize", 30)
output.Instance.SetString("color", "BLUE")
```

and then you can write:

```go
output.Write();
```

#### Setting the fields in a sample instance:
The content of an instance can be set using a dictionary that matches the original type, or field by field:

* **Using a Go type object via JSON encoding**:

```go
// Define a Go type for shape data
// Add an annotation (e.g. json:"color") that indicates an corresponding field in a DDS type
// The Set function uses the built-in encoding function for JSON
type Shape struct {
	Color     string `json:"color"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Shapesize int    `json:"shapesize"`
}

var shape Shape
shape.Y = 2
output.Instance.Set(&shape)
```

 * **Field by field**:

```go
output.Instance.SetInt("y", 2);
```

Nested fields can be accessed with the dot notation: `"x.y.z"`, and array or sequences with square brakets: `"x.y[1].z"`. For more info on how to access
fields, check Section 6.4 'Data Access API' of the
[RTI Prototyper Getting Started Guide](https://community.rti.com/rti-doc/510/ndds.5.1.0/doc/pdf/RTI_CoreLibrariesAndUtilities_Prototyper_GettingStarted.pdf)


#### reading/taking data
To read/take samples first you have to get a reference to the input port:

```go
input, err := connector.GetInput("MySubscriber::MySquareReader")
```

then you can call the `Read()` or `Take()` API:

```go
input.Read();
```

 or

```go
input.Take();
```

The read or take operation can return multiple samples. So we have to iterate on an array:

```go
input.Take()
numOfSamples := input.Samples.GetLength()
for j := 0; j < numOfSamples; j++ {
    if input.Infos.IsValid(j) {
        color := input.Samples.GetString(j, "color")
        x := input.Samples.GetInt(j, "x")
        y := input.Samples.GetInt(j, "y")
        shapesize := input.Samples.GetInt(j, "shapesize")

        log.Println("---Received Sample---")
        log.Printf("color: %s\n", color)
        log.Printf("x: %d\n", x)
        log.Printf("y: %d\n", y)
        log.Printf("shapesize: %d\n", shapesize)
}
```

#### accessing samples fields after a read/take
A `read()` or `take()` operation can return multiple samples. They are stored in an array. Every time you try to access a specific sample you have to specify an index (j in the example below).

You can access the date by getting a copy in a deserialized Go type object or you can access each field individually like the example above:

 * **Using a Go type object via JSON decoding**:

```go
// Define a Go type for shape data
// Add an annotation (e.g. json:"color") that indicates an corresponding field in a DDS type
// The Get function uses the built-in decoding function for JSON
type Shape struct {
        Color     string `json:"color"`
        X         int    `json:"x"`
        Y         int    `json:"y"`
        Shapesize int    `json:"shapesize"`
}

input.Take()
numOfSamples := input.Samples.GetLength()
for j := 0; j < numOfSamples; j++ {
    if input.Infos.IsValid(j) {
        var shape Shape
        err := input.Samples.Get(j, &shape)

        if err != nil {
            log.Println(err)
        }
        log.Println("---Received Sample---")
        log.Printf("color: %s\n", shape.Color)
        log.Printf("x: %d\n", shape.X)
        log.Printf("y: %d\n", shape.Y)
        log.Printf("shapesize: %d\n", shape.Shapesize)
}

```

