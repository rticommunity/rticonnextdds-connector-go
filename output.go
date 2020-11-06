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

// Output publishes DDS data
type Output struct {
	native    unsafe.Pointer // a pointer to a native DataWriter
	connector *Connector
	name      string // name of the native DataWriter
	nameCStr  *C.char
	Instance  *Instance
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
		return nil, errors.New("invalid Publication::DataWriter name")
	}
	output.name = outputName
	output.Instance = newInstance(output)

	connector.Outputs = append(connector.Outputs, *output)

	return output, nil
}

/*******************
* Public Functions *
*******************/

// Write is a function to write a DDS data instance in an output
func (output *Output) Write() error {
	retcode := int(C.RTI_Connector_write(unsafe.Pointer(output.connector.native), output.nameCStr, nil))
	return checkRetcode(retcode)
}

// ClearMembers is a function to initialize a DDS data instance in an output
func (output *Output) ClearMembers() error {
	retcode := int(C.RTI_Connector_clear(unsafe.Pointer(output.connector.native), output.nameCStr))
	return checkRetcode(retcode)
}
