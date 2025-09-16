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
	"encoding/json"
	"errors"
	"unsafe"
)

/********
* Types *
*********/

// Instance is used by an output to write DDS data
type Instance struct {
	output *Output
}

/*******************
* Public Functions *
*******************/

// SetUint8 is a function to set a value of type uint8 into samples
func (instance *Instance) SetUint8(fieldName string, value uint8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetUint16 is a function to set a value of type uint16 into samples
func (instance *Instance) SetUint16(fieldName string, value uint16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetUint32 is a function to set a value of type uint32 into samples
func (instance *Instance) SetUint32(fieldName string, value uint32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetUint64 is a function to set a value of type uint64 into samples
func (instance *Instance) SetUint64(fieldName string, value uint64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetInt8 is a function to set a value of type int8 into samples
func (instance *Instance) SetInt8(fieldName string, value int8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetInt16 is a function to set a value of type int16 into samples
func (instance *Instance) SetInt16(fieldName string, value int16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetInt32 is a function to set a value of type int32 into samples
func (instance *Instance) SetInt32(fieldName string, value int32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetInt64 is a function to set a value of type int64 into samples
func (instance *Instance) SetInt64(fieldName string, value int64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetUint is a function to set a value of type uint into samples
func (instance *Instance) SetUint(fieldName string, value uint) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetInt is a function to set a value of type int into samples
func (instance *Instance) SetInt(fieldName string, value int) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetFloat32 is a function to set a value of type float32 into samples
func (instance *Instance) SetFloat32(fieldName string, value float32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetFloat64 is a function to set a value of type float64 into samples
func (instance *Instance) SetFloat64(fieldName string, value float64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetString is a function that set a string to a fieldname of the samples
func (instance *Instance) SetString(fieldName, value string) error {
	if instance == nil || instance.output == nil || instance.output.connector == nil {
		return errors.New("instance, output, or connector is null")
	}
	if fieldName == "" {
		return errors.New("fieldName cannot be empty")
	}

	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	valueCStr := C.CString(value)
	defer C.free(unsafe.Pointer(valueCStr))

	retcode := int(C.RTI_Connector_set_string_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, valueCStr))
	return checkRetcode(retcode)
}

// SetByte is a function to set a byte to a fieldname of the samples
func (instance *Instance) SetByte(fieldName string, value byte) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetRune is a function to set rune to a fieldname of the samples
func (instance *Instance) SetRune(fieldName string, value rune) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_set_number_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value)))
	return checkRetcode(retcode)
}

// SetBoolean is a function to set boolean to a fieldname of the samples
func (instance *Instance) SetBoolean(fieldName string, value bool) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	intValue := 0
	if value {
		intValue = 1
	}
	retcode := int(C.RTI_Connector_set_boolean_into_samples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.int(intValue)))
	return checkRetcode(retcode)
}

// SetJSON is a function to set JSON string in the form of slice of bytes into Instance
func (instance *Instance) SetJSON(blob []byte) error {
	jsonCStr := C.CString(string(blob))
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := int(C.RTI_Connector_set_json_instance(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, jsonCStr))
	return checkRetcode(retcode)
}

// Set is a function that consumes an interface
// of multiple samples with different types and value
// TODO - think about a new name for this a function (e.g. SetType, SetFromType, FromType)
func (instance *Instance) Set(v interface{}) error {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return instance.SetJSON(jsonData)
}
