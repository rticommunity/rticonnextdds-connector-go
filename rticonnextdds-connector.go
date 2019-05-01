/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

// Package rti implements functions of RTI Connector for Connext DDS in Go
package rti

// #cgo darwin CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -m64
// #cgo windows CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_WIN32 -DNDDS_DLL_VARIABLE
// #cgo linux,amd64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
// #cgo linux,arm CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX
// #cgo darwin LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Darwin16clang8.0 -lrtiddsconnector -ldl -lm -lpthread
// #cgo windows LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Win64VS2013 -lrtiddsconnector
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Linux2.6gcc4.4.5 -lrtiddsconnector -ldl -lnsl -lm -lpthread -lrt
// #cgo linux,arm LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/armv6vfphLinux3.xgcc4.7.2 -lrtiddsconnector -ldl -lnsl -lm -lpthread -lrt
// #include "rticonnextdds-connector.h"
// #include <stdlib.h>
import "C"
import "errors"
import "unsafe"
import "encoding/json"

/********
* Types *
*********/

// Connector is a container managing DDS inputs and outputs
type Connector struct {
	native  *C.struct_RTIDDSConnector
	Inputs  []Input
	Outputs []Output
}

// Output publishes DDS data
type Output struct {
	native    unsafe.Pointer // a pointer to a native DataWriter
	connector *Connector
	name      string // name of the native DataWriter
	nameCStr  *C.char
	Instance  *Instance
}

// Instance is used by an output to write DDS data
type Instance struct {
	output *Output
}

// Input subscribes to DDS data
type Input struct {
	native    unsafe.Pointer // a pointer to a native DataReader
	connector *Connector
	name      string // name of the native DataReader
	nameCStr  *C.char
	Samples   *Samples
	Infos     *Infos
}

// Samples is a sequence of data samples used by an input to read DDS data
type Samples struct {
	input *Input
}

// Infoss is a sequence of info samples used by an input to read DDS meta data
type Infos struct {
	input *Input
}

/********************
* Private Functions *
********************/

func newInstance(output *Output) (instance *Instance) {
	// Error checking for the output is skipped because it was already checked

	instance = new(Instance)
	instance.output = output

	return instance
}

func newOutput(connector *Connector, outputName string) (output *Output, err error) {
	// Error checking for the connector is skipped because it was already checked

	output = new(Output)
	output.connector = connector

	output.nameCStr = C.CString(outputName)

	output.native = C.RTIDDSConnector_getWriter(unsafe.Pointer(connector.native), output.nameCStr)
	if output.native == nil {
		err = errors.New("Invalid Publication::DataWriter name")
		return nil, err
	}
	output.name = outputName
	output.Instance = newInstance(output)

	connector.Outputs = append(connector.Outputs, *output)

	return output, nil
}

func newInput(connector *Connector, inputName string) (input *Input, err error) {
	// Error checking for the connector is skipped because it was already checked

	input = new(Input)
	input.connector = connector

	input.nameCStr = C.CString(inputName)

	input.native = C.RTIDDSConnector_getReader(unsafe.Pointer(connector.native), input.nameCStr)
	if input.native == nil {
		err = errors.New("Invalid Subscription::DataReader name")
		return nil, err
	}
	input.name = inputName
	input.Samples = newSamples(input)
	input.Infos = newInfos(input)

	connector.Inputs = append(connector.Inputs, *input)

	return input, nil
}

func newSamples(input *Input) (samples *Samples) {
	// Error checking for the input is skipped because it was already checked

	samples = new(Samples)
	samples.input = input
	return samples
}

func newInfos(input *Input) (infos *Infos) {
	// Error checking for the input is skipped because it was already checked

	infos = new(Infos)
	infos.input = input
	return infos
}

/*******************
* Public Functions *
*******************/

// NewConnector is a constructor of Connector
func NewConnector(configName string, fileName string) (connector *Connector, err error) {
	connector = new(Connector)

	configNameCStr := C.CString(configName)
	defer C.free(unsafe.Pointer(configNameCStr))
	fileNameCStr := C.CString(fileName)
	defer C.free(unsafe.Pointer(fileNameCStr))

	connector.native = C.RTIDDSConnector_new(configNameCStr, fileNameCStr, nil)
	if connector.native == nil {
		err = errors.New("Invalid participant profile, xml path or xml profile")
		return nil, err
	}

	return connector, nil
}

// Delete is a destructor of Connector
func (connector *Connector) Delete() (err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return err
	}

	// Delete memory allocated in C layer
	for _, input := range connector.Inputs {
		C.free(unsafe.Pointer(input.nameCStr))
	}
	for _, output := range connector.Outputs {
		C.free(unsafe.Pointer(output.nameCStr))
	}

	C.RTIDDSConnector_delete(connector.native)
	connector.native = nil

	return nil
}

// GetOutput returns an output object
func (connector *Connector) GetOutput(outputName string) (output *Output, err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return nil, err
	}

	output, err = newOutput(connector, outputName)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetInput returns an input object
func (connector *Connector) GetInput(inputName string) (input *Input, err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return nil, err
	}

	input, err = newInput(connector, inputName)
	if err != nil {
		return nil, err
	}
	return input, nil
}

// Wait is a function to block until data is available on an input
func (connector *Connector) Wait(timeout_ms int) (err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return err
	}

	retcode := int(C.RTIDDSConnector_wait(unsafe.Pointer(connector.native), (C.int)(timeout_ms)))
	if retcode == 10 /* DDS_RETCODE_TIMEOUT */ {
		err = errors.New("Timeout")
		return err
	} else if retcode != 0 /* DDS_RETCODE_OK */ {
		err = errors.New("RTIDDSConnector_wait error")
		return err
	}
	return nil
}

// Write is a function to write a DDS data instance in an output
func (output *Output) Write() error {
	// The C function does not return errors. In the futurue, we will check erros this when supported in the C layer
	// CON-24 (for more information)
	C.RTIDDSConnector_write(unsafe.Pointer(output.connector.native), output.nameCStr, nil)
	return nil
}

// ClearMembers is function to initialize a DDS data instance in an output
func (output *Output) ClearMembers() error {
	// The C function does not return errors. In the futurue, we will check erros when supported in C the C layer
	C.RTIDDSConnector_clear(unsafe.Pointer(output.connector.native), output.nameCStr)
	return nil
}

func (instance *Instance) SetUint8(fieldName string, value uint8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetUint16(fieldName string, value uint16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetUint32(fieldName string, value uint32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetUint64(fieldName string, value uint64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetInt8(fieldName string, value int8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetInt16(fieldName string, value int16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetInt32(fieldName string, value int32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetInt64(fieldName string, value int64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetUint(fieldName string, value uint) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetInt(fieldName string, value int) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetFloat32(fieldName string, value float32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetFloat64(fieldName string, value float64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetString(fieldName string, value string) error {

	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value_c_str := C.CString(value)
	defer C.free(unsafe.Pointer(value_c_str))

	C.RTIDDSConnector_setStringIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, value_c_str)

	return nil
}

func (instance *Instance) SetByte(fieldName string, value byte) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetRune(fieldName string, value rune) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

func (instance *Instance) SetBoolean(fieldName string, value bool) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	var int_value int
	if value == true {
		int_value = 1
	} else {
		int_value = 0
	}
	C.RTIDDSConnector_setBooleanIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.int(int_value))
	return nil
}

func (instance *Instance) SetJson(json []byte) error {
	json_c_str := C.CString(string(json))
	defer C.free(unsafe.Pointer(json_c_str))

	C.RTIDDSConnector_setJSONInstance(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, json_c_str)
	return nil
}

// TODO - think about a new name for this function (e.g. SetType, SetFromType, FromType)
func (instance *Instance) Set(v interface{}) (err error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = instance.SetJson(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (input *Input) Read() (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}

	// The C function does not return errors. In the futurue, we will update this when supported in the C layer
	C.RTIDDSConnector_read(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

func (input *Input) Take() (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}
	// The C function does not return errors. In the futurue, we will update this when supported in the C layer
	C.RTIDDSConnector_take(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

func (samples *Samples) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getSamplesLength(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr))
	return length
}

func (samples *Samples) GetUint8(index int, fieldName string) (value uint8) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetUint16(index int, fieldName string) (value uint16) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetUint32(index int, fieldName string) (value uint32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetUint64(index int, fieldName string) (value uint64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetInt8(index int, fieldName string) (value int8) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetInt16(index int, fieldName string) (value int16) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetInt32(index int, fieldName string) (value int32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetInt64(index int, fieldName string) (value int64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetFloat32(index int, fieldName string) (value float32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = float32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetFloat64(index int, fieldName string) (value float64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = float64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetInt(index int, fieldName string) (value int) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetUint(index int, fieldName string) (value uint) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetByte(index int, fieldName string) (value byte) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = byte(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetRune(index int, fieldName string) (value rune) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = rune(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

func (samples *Samples) GetBoolean(index int, fieldName string) bool {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value := int(C.RTIDDSConnector_getBooleanFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	if value != 0 {
		return true
	} else {
		return false
	}
}

func (samples *Samples) GetString(index int, fieldName string) (value string) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = C.GoString((*C.char)(C.RTIDDSConnector_getStringFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr)))
	return value
}

func (samples *Samples) GetJson(index int) (json []byte, e error) {
	jsonCStr := C.RTIDDSConnector_getJSONSample(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1))
	defer C.RTIDDSConnector_freeString((*C.char)(jsonCStr))

	json = []byte(C.GoString((*C.char)(jsonCStr)))

	return json, e
}

func (samples *Samples) Get(index int, v interface{}) (e error) {
	jsonData, e := samples.GetJson(index)
	if e != nil {
		return e
	}

	e = json.Unmarshal(jsonData, &v)
	if e != nil {
		return e
	}

	return e
}

func (infos *Infos) IsValid(index int) (valid bool) {
	memberNameCStr := C.CString("valid_data")
	defer C.free(unsafe.Pointer(memberNameCStr))

	if int(C.RTIDDSConnector_getBooleanFromInfos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr)) != 0 {
		valid = true
	} else {
		valid = false
	}
	return valid
}

func (infos *Infos) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getInfosLength(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr))
	return length
}
