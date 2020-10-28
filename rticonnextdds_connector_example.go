package rti

import (
	"fmt"
	"path"
	"runtime"

	"github.com/rticommunity/rticonnextdds-connector-go/types"
)

func ExampleNewConnector() {
	// Create Connector
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero" //nolint // Example
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	if err := connector.Delete(); err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleConnector_GetInput() {
	// Create Connector
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero" //nolint // Example
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	defer func() {
		if err := connector.Delete(); err != nil {
			// Handle an error
		}
	}()

	input, err := connector.GetInput("MySubscriber::MyReader")
	if input == nil || err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleConnector_GetOutput() {
	// Create Connector
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero" //nolint // Example
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	defer func() {
		if err := connector.Delete(); err != nil {
			// Handle an error
		}
	}()

	output, err := connector.GetInput("MyPublisher::MyWriter")
	if output == nil || err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleInfos_GetIdentity() {
	// Create Connector
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero" //nolint // Example
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	defer func() {
		if err := connector.Delete(); err != nil {
			// Handle an error
		}
	}()

	input, err := connector.GetInput("MySubscriber::MyReader")
	if err != nil {
		// Handle an error
	}
	output, err := connector.GetOutput("MyPublisher::MyWriter")
	if err != nil {
		// Handle an error
	}

	var outputTestData types.Test
	outputTestData.St = "test" //nolint // Example
	if err := output.Instance.Set(&outputTestData); err != nil {
		// Handle an error
	}
	if err := output.Write(); err != nil {
		// Handle an error
	}

	if err := connector.Wait(-1); err != nil {
		// Handle an error
	}
	if err := input.Take(); err != nil {
		// Handle an error
	}
	writerID, err := input.Infos.GetIdentity(0)
	if err != nil {
		// Handle an error
	}
	fmt.Printf("wrtier_guid: %x\n", writerID.WriterGUID)
	fmt.Printf("seuqnece_number: %d\n", writerID.SequenceNumber)
	// Output:
	// seuqnece_number: 1
}
