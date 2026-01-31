/*****************************************************************************
*   (c) 2020 Copyright, Real-Time Innovations.  All rights reserved.         *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

// Package rti implements functions of RTI Connector for Connext DDS in Go.
//
// RTI Connector provides a lightweight, easy-to-use API for RTI Connext DDS
// that enables rapid development of distributed applications. It is built on
// XML-Based Application Creation and Dynamic Data, allowing you to define
// data types and QoS policies in XML configuration files.
//
// Quick Start Example:
//
//	connector, err := rti.NewConnector("MyParticipant", "config.xml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer connector.Delete()
//
//	// Write data
//	output, _ := connector.GetOutput("MyWriter")
//	output.Instance.SetString("color", "RED")
//	output.Write()
//
//	// Read data
//	input, _ := connector.GetInput("MyReader")
//	input.Take()
//	length, _ := input.Samples.GetLength()
//	for i := 0; i < length; i++ {
//	    color, _ := input.Samples.GetString(i, "color")
//	    fmt.Printf("Received: %s\n", color)
//	}
//
// For complete examples, see: https://github.com/rticommunity/rticonnextdds-connector-go/tree/master/examples
package rti

// #include "rticonnextdds-connector.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

/********
* Errors *
*********/

// ErrNoData is returned when there is no data available in the DDS layer
var ErrNoData = errors.New("DDS Exception: No Data")

// ErrTimeout is returned when there is a timeout in the DDS layer
var ErrTimeout = errors.New("DDS Exception: Timeout")

/********
* Types *
*********/

// Connector is a container managing DDS inputs and outputs.
//
// It represents the main entry point for DDS communication, encapsulating
// a DDS DomainParticipant and providing access to DataReaders (inputs) and
// DataWriters (outputs) defined in XML configuration files.
//
// The Connector is not thread-safe. You must provide your own synchronization
// when using it from multiple goroutines.
type Connector struct {
	native  *C.RTI_Connector
	Inputs  []Input  // Collection of available DataReaders
	Outputs []Output // Collection of available DataWriters
}

// SampleHandler is a user-defined function type for processing received DDS samples.
//
// It takes pointers to Samples (containing the actual data) and Infos (containing
// metadata like timestamps, sample states, etc.) and processes the received samples.
//
// Example usage:
//
//	handler := func(samples *rti.Samples, infos *rti.Infos) {
//	    length, _ := samples.GetLength()
//	    for i := 0; i < length; i++ {
//	        if valid, _ := infos.IsValid(i); valid {
//	            color, _ := samples.GetString(i, "color")
//	            fmt.Printf("Received: %s\n", color)
//	        }
//	    }
//	}
type SampleHandler func(samples *Samples, infos *Infos)

const (
	// DDSRetCodeNoData is a Return Code from CGO for no data return
	DDSRetCodeNoData = 11
	// DDSRetCodeTimeout is a Return Code from CGO for timeout code
	DDSRetCodeTimeout = 10
	// DDSRetCodeOK is a Return Code from CGO for good state
	DDSRetCodeOK = 0
)

/*******************
* Public Functions *
*******************/

// NewConnector creates a new Connector instance from XML configuration.
//
// Parameters:
//   - configName: The name of the DomainParticipant configuration from the XML file
//     (format: "ParticipantLibrary::ParticipantName")
//   - url: The location of XML configuration documents
//
// URL formats supported:
//   - File path: "file:///usr/local/config.xml" or "/usr/local/config.xml"
//   - Inline XML: "str://<dds><qos_library>...</qos_library></dds>"
//
// Returns:
//   - *Connector: A new Connector instance ready for use
//   - error: Non-nil if the configuration is invalid or file cannot be loaded
//
// Example:
//
//	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", "./ShapeExample.xml")
//	if err != nil {
//	    log.Fatal("Failed to create connector:", err)
//	}
//	defer connector.Delete()
func NewConnector(configName, url string) (*Connector, error) {
	connector := new(Connector)

	configNameCStr := C.CString(configName)
	defer C.free(unsafe.Pointer(configNameCStr))
	urlCStr := C.CString(url)
	defer C.free(unsafe.Pointer(urlCStr))

	connector.native = C.RTI_Connector_new(configNameCStr, urlCStr, nil)
	if connector.native == nil {
		return nil, errors.New("invalid participant profile, xml path or xml profile")
	}

	return connector, nil
}

// Delete is a destructor of Connector
func (connector *Connector) Delete() error {
	if connector == nil {
		return errors.New("connector is null")
	}

	// Check if already deleted
	if connector.native == nil {
		return nil
	}

	// Delete memory allocated in C layer
	for _, input := range connector.Inputs {
		C.free(unsafe.Pointer(input.nameCStr))
	}
	for _, output := range connector.Outputs {
		C.free(unsafe.Pointer(output.nameCStr))
	}

	C.RTI_Connector_delete(connector.native)
	connector.native = nil

	// Clear the slices to prevent double-free if called again
	connector.Inputs = nil
	connector.Outputs = nil

	return nil
}

// isValid checks if the connector is valid and not deleted
func (connector *Connector) isValid() error {
	if connector == nil {
		return errors.New("connector is null")
	}
	if connector.native == nil {
		return errors.New("connector has been deleted")
	}
	return nil
}

// GetOutput returns an output object
func (connector *Connector) GetOutput(outputName string) (*Output, error) {
	if err := connector.isValid(); err != nil {
		return nil, err
	}

	return newOutput(connector, outputName)
}

// GetInput returns an input object
func (connector *Connector) GetInput(inputName string) (*Input, error) {
	if err := connector.isValid(); err != nil {
		return nil, err
	}

	return newInput(connector, inputName)
}

// Wait is a function to block until data is available on an input
func (connector *Connector) Wait(timeoutMs int) error {
	if err := connector.isValid(); err != nil {
		return err
	}

	retcode := int(C.RTI_Connector_wait_for_data(unsafe.Pointer(connector.native), C.int(timeoutMs)))
	return checkRetcode(retcode)
}

/********************
* Private Functions *
********************/

func newOutput(connector *Connector, outputName string) (*Output, error) {
	// Error checking for the connector is skipped because it was already checked

	output := new(Output)
	output.connector = connector

	output.nameCStr = C.CString(outputName)

	output.native = C.RTI_Connector_get_datawriter(unsafe.Pointer(connector.native), output.nameCStr)
	if output.native == nil {
		// Free the allocated C string before returning error
		C.free(unsafe.Pointer(output.nameCStr))
		return nil, errors.New("invalid Publication::DataWriter name")
	}
	output.name = outputName
	output.Instance = newInstance(output)

	connector.Outputs = append(connector.Outputs, *output)

	return output, nil
}

func newInput(connector *Connector, inputName string) (*Input, error) {
	// Error checking for the connector is skipped because it was already checked

	input := new(Input)
	input.connector = connector

	input.nameCStr = C.CString(inputName)

	input.native = C.RTI_Connector_get_datareader(unsafe.Pointer(connector.native), input.nameCStr)
	if input.native == nil {
		// Free the allocated C string before returning error
		C.free(unsafe.Pointer(input.nameCStr))
		return nil, errors.New("invalid Subscription::DataReader name")
	}
	input.name = inputName
	input.Samples = newSamples(input)
	input.Infos = newInfos(input)

	connector.Inputs = append(connector.Inputs, *input)

	return input, nil
}

func newInstance(output *Output) *Instance {
	// Error checking for the output is skipped because it was already checked
	return &Instance{
		output: output,
	}
}

func newSamples(input *Input) *Samples {
	// Error checking for the input is skipped because it was already checked
	return &Samples{
		input: input,
	}
}

func newInfos(input *Input) *Infos {
	// Error checking for the input is skipped because it was already checked
	return &Infos{
		input: input,
	}
}

// checkRetcode is a function to check return code
func checkRetcode(retcode int) error {
	switch retcode {
	case DDSRetCodeOK:
	case DDSRetCodeNoData:
		return ErrNoData
	case DDSRetCodeTimeout:
		return ErrTimeout
	default:
		// Try to get detailed error message from C library
		errMsg := C.GoString((*C.char)(C.RTI_Connector_get_last_error_message()))
		if errMsg == "" {
			// If no detailed message available, provide context based on return code
			return fmt.Errorf("DDS Exception: error code %d (no detailed message available from RTI Connector)", retcode)
		}
		return fmt.Errorf("DDS Exception: %s (error code %d)", errMsg, retcode)
	}
	return nil
}
