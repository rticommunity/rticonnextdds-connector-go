package rti

import (
	"path"
	"runtime"
)

func ExampleNewConnector() {
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero"
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	connector.Delete()
	// Output:
}

func ExampleConnector_GetInput() {
	// Create Connector
	connector := newTestConnector()
	defer connector.Delete()

	input, err := connector.GetInput("MySubscriber::MyReader")
	if input == nil || err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleConnector_GetOutput() {
	// Create Connector
	connector := newTestConnector()
	defer connector.Delete()

	output, err := connector.GetInput("MyPublisher::MyWriter")
	if output == nil || err != nil {
		// Handle an error
	}
	// Output:
}
