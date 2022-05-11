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
	"encoding/json"
	"unsafe"
)

/********
* Types *
*********/

// Samples is a sequence of data samples used by an input to read DDS data
type Samples struct {
	input *Input
}

// getNumber is a function to return a number in double from a sample
func (samples *Samples) getNumber(index int, fieldName string, retVal *C.double) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_get_number_from_sample(unsafe.Pointer(samples.input.connector.native), retVal, samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return checkRetcode(retcode)
}

/*******************
* Public Functions *
*******************/

// GetLength is a function to get the number of samples
func (samples *Samples) GetLength() (int, error) {
	var retVal C.double
	retcode := int(C.RTI_Connector_get_sample_count(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, &retVal))
	err := checkRetcode(retcode)
	return int(retVal), err
}

// GetUint8 is a function to retrieve a value of type uint8 from the samples
func (samples *Samples) GetUint8(index int, fieldName string) (uint8, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return uint8(retVal), err
}

// GetUint16 is a function to retrieve a value of type uint16 from the samples
func (samples *Samples) GetUint16(index int, fieldName string) (uint16, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return uint16(retVal), err
}

// GetUint32 is a function to retrieve a value of type uint32 from the samples
func (samples *Samples) GetUint32(index int, fieldName string) (uint32, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return uint32(retVal), err
}

// GetUint64 is a function to retrieve a value of type uint64 from the samples
func (samples *Samples) GetUint64(index int, fieldName string) (uint64, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return uint64(retVal), err
}

// GetInt8 is a function to retrieve a value of type int8 from the samples
func (samples *Samples) GetInt8(index int, fieldName string) (int8, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return int8(retVal), err
}

// GetInt16 is a function to retrieve a value of type int16 from the samples
func (samples *Samples) GetInt16(index int, fieldName string) (int16, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return int16(retVal), err
}

// GetInt32 is a function to retrieve a value of type int32 from the samples
func (samples *Samples) GetInt32(index int, fieldName string) (int32, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return int32(retVal), err
}

// GetInt64 is a function to retrieve a value of type int64 from the samples
func (samples *Samples) GetInt64(index int, fieldName string) (int64, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return int64(retVal), err
}

// GetFloat32 is a function to retrieve a value of type float32 from the samples
func (samples *Samples) GetFloat32(index int, fieldName string) (float32, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return float32(retVal), err
}

// GetFloat64 is a function to retrieve a value of type float64 from the samples
func (samples *Samples) GetFloat64(index int, fieldName string) (float64, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return float64(retVal), err
}

// GetInt is a function to retrieve a value of type int from the samples
func (samples *Samples) GetInt(index int, fieldName string) (int, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return int(retVal), err
}

// GetUint is a function to retrieve a value of type uint from the samples
func (samples *Samples) GetUint(index int, fieldName string) (uint, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return uint(retVal), err
}

// GetByte is a function to retrieve a value of type byte from the samples
func (samples *Samples) GetByte(index int, fieldName string) (byte, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return byte(retVal), err
}

// GetRune is a function to retrieve a value of type rune from the samples
func (samples *Samples) GetRune(index int, fieldName string) (rune, error) {
	var retVal C.double
	err := samples.getNumber(index, fieldName, &retVal)
	return rune(retVal), err
}

// GetBoolean is a function to retrieve a value of type boolean from the samples
func (samples *Samples) GetBoolean(index int, fieldName string) (bool, error) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	var retVal C.int

	retcode := int(C.RTI_Connector_get_boolean_from_sample(unsafe.Pointer(samples.input.connector.native), &retVal, samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	err := checkRetcode(retcode)

	return (retVal != 0), err
}

// GetString is a function to retrieve a value of type string from the samples
func (samples *Samples) GetString(index int, fieldName string) (string, error) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	var retValCStr *C.char

	retcode := int(C.RTI_Connector_get_string_from_sample(unsafe.Pointer(samples.input.connector.native), &retValCStr, samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	err := checkRetcode(retcode)
	if err != nil {
		return "", err
	}

	retValGoStr := C.GoString(retValCStr)
	C.RTI_Connector_free_string(retValCStr)

	return retValGoStr, nil
}

// GetJSON is a function to retrieve a slice of bytes of a JSON string from the samples
func (samples *Samples) GetJSON(index int) ([]byte, error) {
	var retValCStr *C.char

	retcode := int(C.RTI_Connector_get_json_sample(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), &retValCStr))
	err := checkRetcode(retcode)
	if err != nil {
		return nil, err
	}

	retValGoStr := C.GoString(retValCStr)
	C.RTI_Connector_free_string(retValCStr)

	return []byte(retValGoStr), err
}

// Get is a function to retrieve all the information
// of the samples and put it into an interface
func (samples *Samples) Get(index int, v interface{}) error {
	jsonData, err := samples.GetJSON(index)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, &v)
}
