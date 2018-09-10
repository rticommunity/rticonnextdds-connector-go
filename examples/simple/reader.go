package main

import (
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
	"path"
	"runtime"
	"time"
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

	// Get an input from the connector
	input, err := connector.GetInput("MySubscriber::MySquareReader")
	if err != nil {
		log.Panic(err)
	}

	// Get values from a received sample and print them
	for i := 0; i < 500; i++ {
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
		}
		time.Sleep(time.Second * 1)
	}
}
