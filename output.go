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

/*******************
* Public Functions *
*******************/

// Write is a function to write a DDS data instance in an output
func (output *Output) Write() error {
	if output == nil {
		return errors.New("output is null")
	}

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
//
//	  identity={"writer_guid":[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15], "sequence_number":1},
//		 timestamp=1000000000)
func (output *Output) WriteWithParams(jsonStr string) error {
	if output == nil {
		return errors.New("output is null")
	}

	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := int(C.RTI_Connector_write(unsafe.Pointer(output.connector.native), output.nameCStr, jsonCStr))
	return checkRetcode(retcode)
}

// ClearMembers is a function to initialize a DDS data instance in an output
func (output *Output) ClearMembers() error {
	if output == nil {
		return errors.New("output is null")
	}

	retcode := int(C.RTI_Connector_clear(unsafe.Pointer(output.connector.native), output.nameCStr))
	return checkRetcode(retcode)
}

// Waits until the number of matched DDS subscription changes
// This method waits until new compatible subscriptions are discovered or
// existing subscriptions no longer match.

// Parameters:
//   timeout: The maximum time to wait in milliseconds. By default, infinite.

// Return: The change in the current number of matched outputs. If a positive number is returned, the input has matched with new publishers. If a negative number is returned the input has unmatched from an output. It is possible for multiple matches and/or unmatches to be returned (e.g., 0 could be returned, indicating that the input matched the same number of writers as it unmatched).
func (output *Output) WaitForSubscriptions(timeoutMs int) (int, error) {
	if output == nil {
		return -1, errors.New("output is null")
	}

	var currentCountChange C.int

	retcode := int(C.RTI_Connector_wait_for_matched_subscription(unsafe.Pointer(output.native), C.int(timeoutMs), &currentCountChange))
	return int(currentCountChange), checkRetcode(retcode)
}

// Returns information about the matched subscriptions

// This function returns a JSON string where each element is a dictionary with
// information about a subscription matched with this Output.

// Currently, the only key in the dictionaries is ``"name"``
// containing the subscription name. If a subscription doesn't have name,
// the value is ``None``.

// Note that Connector Inputs are automatically assigned a name from the
// *data_reader name* in the XML configuration.
func (output *Output) GetMatchedSubscriptions() (string, error) {
	if output == nil {
		return "", errors.New("output is null")
	}

	var jsonCStr *C.char

	retcode := int(C.RTI_Connector_get_matched_subscriptions(unsafe.Pointer(output.native), &jsonCStr))
	err := checkRetcode(retcode)
	if err != nil {
		return "", err
	}

	jsonGoStr := C.GoString(jsonCStr)
	C.RTI_Connector_free_string(jsonCStr)

	return jsonGoStr, nil
}
