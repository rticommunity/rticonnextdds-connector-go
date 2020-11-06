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

/********************
* Private Functions *
********************/

func newInput(connector *Connector, inputName string) (*Input, error) {
	// Error checking for the connector is skipped because it was already checked

	input := new(Input)
	input.connector = connector

	input.nameCStr = C.CString(inputName)

	input.native = C.RTI_Connector_get_datareader(unsafe.Pointer(connector.native), input.nameCStr)
	if input.native == nil {
		return nil, errors.New("invalid Subscription::DataReader name")
	}
	input.name = inputName
	input.Samples = newSamples(input)
	input.Infos = newInfos(input)

	connector.Inputs = append(connector.Inputs, *input)

	return input, nil
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
