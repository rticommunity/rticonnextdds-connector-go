package main

import (
	"github.com/kyoungho/rticonnextdds-connector"
	"path"
	"runtime"
	"time"
	"log"
)

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../ShapeExample.xml")

	// Create a connector defined in the XML configuration
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
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
