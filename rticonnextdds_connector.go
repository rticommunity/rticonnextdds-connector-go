/*****************************************************************************
*   (c) 2020 Copyright, Real-Time Innovations.  All rights reserved.         *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

// Package rti implements functions of RTI Connector for Connext DDS in Go
package rti

// #cgo windows CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_WIN32 -DNDDS_DLL_VARIABLE
// #cgo linux,arm CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX
// #cgo windows LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Win64VS2013 -lrtiddsconnector
// #cgo linux,arm LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/armv6vfphLinux3.xgcc4.7.2 -lrtiddsconnector -ldl -lnsl -lm -lpthread -lrt
// #include "rticonnextdds-connector.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
)

/********
* Types *
*********/

// Connector is a container managing DDS inputs and outputs
type Connector struct {
	native  *C.RTI_Connector
	Inputs  []Input
	Outputs []Output
}

// SampleHandler is an User defined function type that takes in pointers of
// Samples and Infos and will handle received samples.
type SampleHandler func(samples *Samples, infos *Infos)

const (
	// DDSRetCodeNoData is a Return Code from CGO for no data return
	DDSRetCodeNoData = 11
	// DDSRetCodeTimeout is a Return Code from CGO for timeout code
	DDSRetCodeTimeout = 10
	// DDSRetCodeOK is a Return Code from CGO for good state
	DDSRetCodeOK = 0
)

/********************
* Private Functions *
********************/

// checkRetcode is a function to check return code
func checkRetcode(retcode int) error {
	switch retcode {
	case DDSRetCodeOK:
	case DDSRetCodeNoData:
		return errors.New("DDS Exceptrion: No Data")
	case DDSRetCodeTimeout:
		return errors.New("DDS Exception: Timeout")
	default:
		return errors.New("DDS Exception: " + C.GoString((*C.char)(C.RTI_Connector_get_last_error_message)))
	}
	return nil
}

/*******************
* Public Functions *
*******************/

// NewConnector is a constructor of Connector.
//
// url is the location of XML documents in URL format. For example:
//  File specification: file:///usr/local/default_dds.xml
//  String specification: str://"<dds><qos_library>…</qos_library></dds>"
// If you omit the URL schema name, Connector will assume a file name. For example:
//  File Specification: /usr/local/default_dds.xml
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

	// Delete memory allocated in C layer
	for _, input := range connector.Inputs {
		C.free(unsafe.Pointer(input.nameCStr))
	}
	for _, output := range connector.Outputs {
		C.free(unsafe.Pointer(output.nameCStr))
	}

	C.RTI_Connector_delete(connector.native)
	connector.native = nil

	return nil
}

// GetOutput returns an output object
func (connector *Connector) GetOutput(outputName string) (*Output, error) {
	if connector == nil {
		return nil, errors.New("connector is null")
	}

	return newOutput(connector, outputName)
}

// GetInput returns an input object
func (connector *Connector) GetInput(inputName string) (*Input, error) {
	if connector == nil {
		return nil, errors.New("connector is null")
	}

	return newInput(connector, inputName)
}

// Wait is a function to block until data is available on an input
func (connector *Connector) Wait(timeoutMs int) error {
	if connector == nil {
		return errors.New("connector is null")
	}

	retcode := int(C.RTI_Connector_wait_for_data(unsafe.Pointer(connector.native), C.int(timeoutMs)))
	return checkRetcode(retcode)
}
