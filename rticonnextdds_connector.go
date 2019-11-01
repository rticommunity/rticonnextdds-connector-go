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

// #cgo windows CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_WIN32 -DNDDS_DLL_VARIABLE
// #cgo linux,arm CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX
// #cgo windows LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Win64VS2013 -lrtiddsconnector
// #cgo linux,arm LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/armv6vfphLinux3.xgcc4.7.2 -lrtiddsconnector -ldl -lnsl -lm -lpthread -lrt
// #include "rticonnextdds-connector.h"
// #include <stdlib.h>
import "C"
import "errors"
import "unsafe"
import "encoding/json"
import "strconv"
//import "fmt"


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

// Infos is a sequence of info samples used by an input to read DDS meta data
type Infos struct {
	input *Input
}

// Identity is the structure for identifying
type Identity struct {
	WriterGuid     [16]byte `json:"writer_guid"`
	SequenceNumber uint     `json:"sequence_number"`
}

// SampleHandler is an User defined function type that takes in pointers of
// Samples and Infos and will handle received samples.
type SampleHandler func(samples *Samples, infos *Infos)

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

// NewConnector is a constructor of Connector.
//
// url is the location of XML documents in URL format. For example:
//  File specification: file:///usr/local/default_dds.xml
//  String specification: str://"<dds><qos_library>…</qos_library></dds>"
// If you omit the URL schema name, Connector will assume a file name. For example:
//  File Specification: /usr/local/default_dds.xml
func NewConnector(configName string, url string) (connector *Connector, err error) {
	connector = new(Connector)

	configNameCStr := C.CString(configName)
	defer C.free(unsafe.Pointer(configNameCStr))
	urlCStr := C.CString(url)
	defer C.free(unsafe.Pointer(urlCStr))

	connector.native = C.RTIDDSConnector_new(configNameCStr, urlCStr, nil)
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
func (connector *Connector) Wait(timeoutMs int) (err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return err
	}

	retcode := int(C.RTIDDSConnector_wait(unsafe.Pointer(connector.native), (C.int)(timeoutMs)))
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
	// The C function does not return errors. In the future, we will check errors when supported in the C layer
	// CON-24 (for more information)
	C.RTIDDSConnector_write(unsafe.Pointer(output.connector.native), output.nameCStr, nil)
	return nil
}

// ClearMembers is a function to initialize a DDS data instance in an output
func (output *Output) ClearMembers() error {
	// The C function does not return errors. In the future, we will check errors when supported in the C layer
	C.RTIDDSConnector_clear(unsafe.Pointer(output.connector.native), output.nameCStr)
	return nil
}

// SetUint8 is a function to set a value of type uint8 into samples
func (instance *Instance) SetUint8(fieldName string, value uint8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetUint16 is a function to set a value of type uint16 into samples
func (instance *Instance) SetUint16(fieldName string, value uint16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetUint32 is a function to set a value of type uint32 into samples
func (instance *Instance) SetUint32(fieldName string, value uint32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetUint64 is a function to set a value of type uint64 into samples
func (instance *Instance) SetUint64(fieldName string, value uint64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetInt8 is a function to set a value of type int8 into samples
func (instance *Instance) SetInt8(fieldName string, value int8) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetInt16 is a function to set a value of type int16 into samples
func (instance *Instance) SetInt16(fieldName string, value int16) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetInt32 is a function to set a value of type int32 into samples
func (instance *Instance) SetInt32(fieldName string, value int32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetInt64 is a function to set a value of type int64 into samples
func (instance *Instance) SetInt64(fieldName string, value int64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetUint is a function to set a value of type uint into samples
func (instance *Instance) SetUint(fieldName string, value uint) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetInt is a function to set a value of type int into samples
func (instance *Instance) SetInt(fieldName string, value int) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetFloat32 is a function to set a value of type float32 into samples
func (instance *Instance) SetFloat32(fieldName string, value float32) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetFloat64 is a function to set a value of type float64 into samples
func (instance *Instance) SetFloat64(fieldName string, value float64) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetString is a function that set a string to a fieldname of the samples
func (instance *Instance) SetString(fieldName string, value string) error {

	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	valueCStr := C.CString(value)
	defer C.free(unsafe.Pointer(valueCStr))

	C.RTIDDSConnector_setStringIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, valueCStr)

	return nil
}

// SetByte is a function to set a byte to a fieldname of the samples
func (instance *Instance) SetByte(fieldName string, value byte) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetRune is a function to set rune to a fieldname of the samples
func (instance *Instance) SetRune(fieldName string, value rune) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.double(value))
	return nil
}

// SetBoolean is a function to set boolean to a fieldname of the samples
func (instance *Instance) SetBoolean(fieldName string, value bool) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	var intValue int
	if value {
		intValue = 1
	} else {
		intValue = 0
	}
	C.RTIDDSConnector_setBooleanIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.int(intValue))
	return nil
}

// SetJSON is a function to set JSON string in the form of slice of bytes into Instance
func (instance *Instance) SetJSON(json []byte) error {
	jsonCStr := C.CString(string(json))
	defer C.free(unsafe.Pointer(jsonCStr))

	C.RTIDDSConnector_setJSONInstance(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, jsonCStr)
	return nil
}

// Set is a function that consumes an interface
// of multiple samples with different types and value
// TODO - think about a new name for this a function (e.g. SetType, SetFromType, FromType)
func (instance *Instance) Set(v interface{}) (err error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = instance.SetJSON(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Read is a function to read DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Read function
// does not remove DDS samples from the DDS DataReader's receive queue.
func (input *Input) Read() (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}

	// The C function does not return errors. In the future, we will update this when supported in the C layer
	C.RTIDDSConnector_read(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

// Take is a function to take DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Take
// function removes DDS samples from the DDS DataReader's receive queue.
func (input *Input) Take() (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}
	// The C function does not return errors. In the future, we will update this when supported in the C layer
	C.RTIDDSConnector_take(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

/*
// AsyncSubscribe is a function to subscribe DDS samples in an asynchronous way.
// Internllay, it takes DDS samples from the DDS DataReader when they arrive.
// Then, it invokes the callback function (cb SampleHandler) that will handle received samples.
func (input *Input) AsyncSubscribe(cb SampleHandler) (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}
	//input.mu.Lock()
	//defer input.mu.Unlock()
	go func() {
		for {
			input.connector.Wait(-1)
			input.Take()
			cb(input.Samples, input.Infos)
		}
	}()
	return nil
}

// ChannelSubscribe is a function to subscribe DDS samples with a Go channel.
// Internally, it taks DDS samples from the DDS DataReader when they arrive.
// Then, it sends arrived DDS samples to the channel (samples chan *Samples).
func (input *Input) ChannelSubscribe(samples chan *Samples) (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}
	//input.mu.Lock()
	//defer input.mu.Unlock()
	go func() {
		for {
			input.connector.Wait(-1)
			input.Take()
			samples <- input.Samples
		}
	}()
	return nil
}
*/

// GetLength is a function to get the number of samples
func (samples *Samples) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getSamplesLength(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr))
	return length
}

// GetUint8 is a function to retrieve a value of type uint8 from the samples
func (samples *Samples) GetUint8(index int, fieldName string) (value uint8) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetUint16 is a function to retrieve a value of type uint16 from the samples
func (samples *Samples) GetUint16(index int, fieldName string) (value uint16) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetUint32 is a function to retrieve a value of type uint32 from the samples
func (samples *Samples) GetUint32(index int, fieldName string) (value uint32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetUint64 is a function to retrieve a value of type uint64 from the samples
func (samples *Samples) GetUint64(index int, fieldName string) (value uint64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetInt8 is a function to retrieve a value of type int8 from the samples
func (samples *Samples) GetInt8(index int, fieldName string) (value int8) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetInt16 is a function to retrieve a value of type int16 from the samples
func (samples *Samples) GetInt16(index int, fieldName string) (value int16) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetInt32 is a function to retrieve a value of type int32 from the samples
func (samples *Samples) GetInt32(index int, fieldName string) (value int32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetInt64 is a function to retrieve a value of type int64 from the samples
func (samples *Samples) GetInt64(index int, fieldName string) (value int64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetFloat32 is a function to retrieve a value of type float32 from the samples
func (samples *Samples) GetFloat32(index int, fieldName string) (value float32) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = float32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetFloat64 is a function to retrieve a value of type float64 from the samples
func (samples *Samples) GetFloat64(index int, fieldName string) (value float64) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = float64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetInt is a function to retrieve a value of type int from the samples
func (samples *Samples) GetInt(index int, fieldName string) (value int) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = int(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetUint is a function to retrieve a value of type uint from the samples
func (samples *Samples) GetUint(index int, fieldName string) (value uint) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = uint(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetByte is a function to retrieve a value of type byte from the samples
func (samples *Samples) GetByte(index int, fieldName string) (value byte) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = byte(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetRune is a function to retrieve a value of type rune from the samples
func (samples *Samples) GetRune(index int, fieldName string) (value rune) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = rune(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value
}

// GetBoolean is a function to retrieve a value of type boolean from the samples
func (samples *Samples) GetBoolean(index int, fieldName string) bool {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value := int(C.RTIDDSConnector_getBooleanFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return value != 0
}

// GetString is a function to retrieve a value of type string from the samples
func (samples *Samples) GetString(index int, fieldName string) (value string) {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	value = C.GoString((*C.char)(C.RTIDDSConnector_getStringFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr)))
	return value
}

// GetJSON is a function to retrieve a slice of bytes of a JSON string from the samples
func (samples *Samples) GetJSON(index int) (json []byte, e error) {
	jsonCStr := C.RTIDDSConnector_getJSONSample(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1))
	defer C.RTIDDSConnector_freeString((*C.char)(jsonCStr))

	json = []byte(C.GoString((*C.char)(jsonCStr)))

	return json, e
}

// Get is a function to retrieve all the information
// of the samples and put it into an interface
func (samples *Samples) Get(index int, v interface{}) (e error) {
	jsonData, e := samples.GetJSON(index)
	if e != nil {
		return e
	}

	e = json.Unmarshal(jsonData, &v)
	if e != nil {
		return e
	}

	return e
}

// IsValid is a function to check validity of the element and return a boolean
func (infos *Infos) IsValid(index int) bool {
	memberNameCStr := C.CString("valid_data")
	defer C.free(unsafe.Pointer(memberNameCStr))
	var retVal C.int

	C.RTI_Connector_get_boolean_from_infos(unsafe.Pointer(infos.input.connector.native), &retVal, infos.input.nameCStr, C.int(index+1), memberNameCStr)

	if retVal != 0 {
		return true
	}
	return false
}

// GetSourceTimestamp is a function to get the source timestamp of a sample
func (infos *Infos) GetSourceTimestamp(index int) (ts int, err error) {
	memberNameCStr := C.CString("source_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

    var jsonStr string
    jsonCStr := C.CString(jsonStr)
    defer C.free(unsafe.Pointer(jsonCStr))

	retcode := C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &jsonCStr)
    if retcode != 0 /* DDS_RETCODE_OK */ {
        err = errors.New("RTI_Connector_get_json_from_infos failed")
        return ts, err
    }

	ts, err = strconv.Atoi(C.GoString((*C.char)(jsonCStr)))
	if err != nil {
		err = errors.New("String conversion failed")
		return ts, err
	}

	return ts, err
}

// GetReceptionTimestamp is a function to get the reception timestamp of a sample
func (infos *Infos) GetReceptionTimestamp(index int) (ts int, err error) {
	memberNameCStr := C.CString("reception_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

    var jsonStr string
    jsonCStr := C.CString(jsonStr)
    defer C.free(unsafe.Pointer(jsonCStr))

	retcode := C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &jsonCStr)
    if retcode != 0 /* DDS_RETCODE_OK */ {
        err = errors.New("RTI_Connector_get_json_from_infos failed")
        return ts, err
    }

	ts, err = strconv.Atoi(C.GoString((*C.char)(jsonCStr)))
	if err != nil {
		err = errors.New("String conversion failed")
		return ts, err
	}

	return ts, err
}

// GetIdentity is a function to get the identity of a writer that sent the sample
func (infos *Infos) GetIdentity(index int) (writerId Identity, err error) {
	memberNameCStr := C.CString("identity")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var jsonStr string
	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &jsonCStr)
	if retcode != 0 /* DDS_RETCODE_OK */ {
		err = errors.New("RTI_Connector_get_json_from_infos failed")
		return writerId, err
	}

	jsonByte := []byte(C.GoString((*C.char)(jsonCStr)))
	err = json.Unmarshal(jsonByte, &writerId)
	if err != nil {
		err = errors.New("JSON Unmarshal failed")
		return writerId, err
	}

	return writerId, nil
}

// GetLength is a function to return the length of the
func (infos *Infos) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getInfosLength(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr))
	return length
}
