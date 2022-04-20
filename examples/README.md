RTI Connext Go Connector Examples
========

### Example Overview
#### Import the Connector library
If you want to use the Go Connector, you have to import the package.

```go
import "github.com/rticommunity/rticonnextdds-connector"
```

#### Instantiate a new connector
To create a new connector you have to pass a location of an XML configuration file and a configuration name in XML. For more information on the XML format check the [XML App Creation guide](https://community.rti.com/static/documentation/connext-dds/6.0.0/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf).

```go
connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
```
#### Delete a connector
To destroy all the DDS entities created by a connector, you should call the `Delete()`.

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
The content of an instance can be set using a Go type that matches the original DDS type, or field by field:

* **Using a Go type with JSON encoding**:

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

Nested fields can be accessed with the dot notation: `"x.y.z"`, and array or sequences with square brakets: `"x.y[1].z"`. 

#### Reading/taking samples
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

 * **Field by field**:
You can access each field individually like the example below. 
A `Read()` or `Take()` can return multiple samples. They are stored in an array. 
Every time you try to access a specific sample you have to specify an index (j in the example below).

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

 * **Using a Go type with JSON decoding**:
You can access sample data in a deserialized Go type object.  

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
