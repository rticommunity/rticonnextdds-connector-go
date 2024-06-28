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

// Input subscribes to DDS data
type Input struct {
	native    unsafe.Pointer // a pointer to a native DataReader
	connector *Connector
	name      string // name of the native DataReader
	nameCStr  *C.char
	Samples   *Samples
	Infos     *Infos
}

/*******************
* Public Functions *
*******************/

// Read is a function to read DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Read function
// does not remove DDS samples from the DDS DataReader's receive queue.
func (input *Input) Read() error {
	if input == nil {
		return errors.New("input is null")
	}

	retcode := int(C.RTI_Connector_read(unsafe.Pointer(input.connector.native), input.nameCStr))
	return checkRetcode(retcode)
}

// Take is a function to take DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Take
// function removes DDS samples from the DDS DataReader's receive queue.
func (input *Input) Take() error {
	if input == nil {
		return errors.New("input is null")
	}

	retcode := int(C.RTI_Connector_take(unsafe.Pointer(input.connector.native), input.nameCStr))
	return checkRetcode(retcode)
}

// Waits until this input matches or unmatches a compatible DDS subscription.
// If the operation times out, it will raise :class:`TimeoutError`.
// Parameters:
//   timeout: The maximum time to wait in milliseconds. Set -1 if you want infinite.
// Return: The change in the current number of matched outputs. If a positive number is returned, the input has matched with new publishers. If a negative number is returned the input has unmatched from an output. It is possible for multiple matches and/or unmatches to be returned (e.g., 0 could be returned, indicating that the input matched the same number of writers as it unmatched).
func (input *Input) WaitForPublications(timeoutMs int) (int, error) {
	if input == nil {
		return -1, errors.New("input is null")
	}

	var currentCountChange C.int

	retcode := int(C.RTI_Connector_wait_for_matched_publication(unsafe.Pointer(input.native), C.int(timeoutMs), &currentCountChange))
	return int(currentCountChange), checkRetcode(retcode)
}

// Returns information about the matched publications
// This function returns a JSON string where each element is a dictionary with
// information about a publication matched with this Input.

// Currently, the only key in the dictionaries is ``"name"``,
// containing the publication name. If a publication doesn't have name,
// the value for the key ``name`` is ``None``.

// Note that Connector Outputs are automatically assigned a name from the
// *data_writer name* in the XML configuration.
func (input *Input) GetMatchedPublications() (string, error) {
        if input == nil {
                return "", errors.New("input is null")
        }

	var jsonCStr *C.char

	retcode := int(C.RTI_Connector_get_matched_publications(unsafe.Pointer(input.native), &jsonCStr))
	err := checkRetcode(retcode)
	if err != nil {
		return "", err
	}

	jsonGoStr := C.GoString(jsonCStr)
	C.RTI_Connector_free_string(jsonCStr)

	return jsonGoStr, nil
}
