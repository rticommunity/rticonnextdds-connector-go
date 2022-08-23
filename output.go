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

/*******************
* Public Functions *
*******************/

// Write is a function to write a DDS data instance in an output
func (output *Output) Write() error {
	retcode := int(C.RTI_Connector_write(unsafe.Pointer(output.connector.native), output.nameCStr, nil))
	return checkRetcode(retcode)
}

// WriteWithParams is a function to write a DDS data instance with parameters
// The supported parameters are:
// action: One of "write" (default), "dispose" or "unregister"
// source_timestamp: The source timestamp, an integer representing the total number of nanoseconds
// identity: A dictionary containing the keys "writer_guid" (a list of 16 bytes) and "sequence_number" (an integer) that uniquely identifies this sample.
// related_sample_identity: Used for request-reply communications. It has the same format as "identity"
// For example::
// output.Write(
//   identity={"writer_guid":[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15], "sequence_number":1},
//	 timestamp=1000000000)
func (output *Output) WriteWithParams(jsonStr string) error {
	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := int(C.RTI_Connector_write(unsafe.Pointer(output.connector.native), output.nameCStr, jsonCStr))
	return checkRetcode(retcode)
}

// ClearMembers is a function to initialize a DDS data instance in an output
func (output *Output) ClearMembers() error {
	retcode := int(C.RTI_Connector_clear(unsafe.Pointer(output.connector.native), output.nameCStr))
	return checkRetcode(retcode)
}
