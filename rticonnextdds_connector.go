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
import (
	"encoding/json"
	"errors"
	"strconv"
	"unsafe"
)

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
	WriterGUID     [16]byte `json:"writer_guid"`
	SequenceNumber uint     `json:"sequence_number"`
}

// SampleHandler is an User defined function type that takes in pointers of
// Samples and Infos and will handle received samples.
type SampleHandler func(samples *Samples, infos *Infos)

const (
	// DDSRetCodeTimeout is a Return Code from CGO for timeout code
	DDSRetCodeTimeout = 10
	// DDSRetCodeOK is a Return Code from CGO for good state
	DDSRetCodeOK = 0
)

/********************
* Private Functions *
********************/

func newInstance(output *Output) *Instance {
	// Error checking for the output is skipped because it was already checked
	return &Instance{
		output: output,
	}
}

func newOutput(connector *Connector, outputName string) (*Output, error) {
	// Error checking for the connector is skipped because it was already checked

	output := new(Output)
	output.connector = connector

	output.nameCStr = C.CString(outputName)

	output.native = C.RTIDDSConnector_getWriter(unsafe.Pointer(connector.native), output.nameCStr)
	if output.native == nil {
		return nil, errors.New("invalid Publication::DataWriter name")
	}
	output.name = outputName
	output.Instance = newInstance(output)

	connector.Outputs = append(connector.Outputs, *output)

	return output, nil
}

func newInput(connector *Connector, inputName string) (*Input, error) {
	// Error checking for the connector is skipped because it was already checked

	input := new(Input)
	input.connector = connector

	input.nameCStr = C.CString(inputName)

	input.native = C.RTIDDSConnector_getReader(unsafe.Pointer(connector.native), input.nameCStr)
	if input.native == nil {
		return nil, errors.New("invalid Subscription::DataReader name")
	}
	input.name = inputName
	input.Samples = newSamples(input)
	input.Infos = newInfos(input)

	connector.Inputs = append(connector.Inputs, *input)

	return input, nil
}

func newSamples(input *Input) *Samples {
	// Error checking for the input is skipped because it was already checked
	return &Samples{
		input: input,
	}
}

func newInfos(input *Input) *Infos {
	// Error checking for the input is skipped because it was already checked
	return &Infos{
		input: input,
	}
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

	connector.native = C.RTIDDSConnector_new(configNameCStr, urlCStr, nil)
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

	C.RTIDDSConnector_delete(connector.native)
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
	switch int(C.RTIDDSConnector_wait(unsafe.Pointer(connector.native), C.int(timeoutMs))) {
	case DDSRetCodeOK:
		return nil
	case DDSRetCodeTimeout:
		return errors.New("timeout")
	default:
	}
	return errors.New("RTIDDSConnector_wait error")
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
func (instance *Instance) SetString(fieldName, value string) error {
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

	intValue := 0
	if value {
		intValue = 1
	}
	C.RTIDDSConnector_setBooleanIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, fieldNameCStr, C.int(intValue))
	return nil
}

// SetJSON is a function to set JSON string in the form of slice of bytes into Instance
func (instance *Instance) SetJSON(blob []byte) error {
	jsonCStr := C.CString(string(blob))
	defer C.free(unsafe.Pointer(jsonCStr))

	C.RTIDDSConnector_setJSONInstance(unsafe.Pointer(instance.output.connector.native), instance.output.nameCStr, jsonCStr)
	return nil
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

// Read is a function to read DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Read function
// does not remove DDS samples from the DDS DataReader's receive queue.
func (input *Input) Read() error {
	if input == nil {
		return errors.New("input is null")
	}

	// The C function does not return errors. In the future, we will update this when supported in the C layer
	C.RTIDDSConnector_read(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

// Take is a function to take DDS samples from the DDS DataReader
// and allow access them via the Connector Samples. The Take
// function removes DDS samples from the DDS DataReader's receive queue.
func (input *Input) Take() error {
	if input == nil {
		return errors.New("input is null")
	}
	// The C function does not return errors. In the future, we will update this when supported in the C layer
	C.RTIDDSConnector_take(unsafe.Pointer(input.connector.native), input.nameCStr)
	return nil
}

// GetLength is a function to get the number of samples
func (samples *Samples) GetLength() int {
	return int(C.RTIDDSConnector_getSamplesLength(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr))
}

// GetUint8 is a function to retrieve a value of type uint8 from the samples
func (samples *Samples) GetUint8(index int, fieldName string) uint8 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return uint8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetUint16 is a function to retrieve a value of type uint16 from the samples
func (samples *Samples) GetUint16(index int, fieldName string) uint16 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return uint16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetUint32 is a function to retrieve a value of type uint32 from the samples
func (samples *Samples) GetUint32(index int, fieldName string) uint32 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return uint32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetUint64 is a function to retrieve a value of type uint64 from the samples
func (samples *Samples) GetUint64(index int, fieldName string) uint64 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return uint64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetInt8 is a function to retrieve a value of type int8 from the samples
func (samples *Samples) GetInt8(index int, fieldName string) int8 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetInt16 is a function to retrieve a value of type int16 from the samples
func (samples *Samples) GetInt16(index int, fieldName string) int16 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetInt32 is a function to retrieve a value of type int32 from the samples
func (samples *Samples) GetInt32(index int, fieldName string) int32 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetInt64 is a function to retrieve a value of type int64 from the samples
func (samples *Samples) GetInt64(index int, fieldName string) int64 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetFloat32 is a function to retrieve a value of type float32 from the samples
func (samples *Samples) GetFloat32(index int, fieldName string) float32 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return float32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetFloat64 is a function to retrieve a value of type float64 from the samples
func (samples *Samples) GetFloat64(index int, fieldName string) float64 {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return float64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetInt is a function to retrieve a value of type int from the samples
func (samples *Samples) GetInt(index int, fieldName string) int {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetUint is a function to retrieve a value of type uint from the samples
func (samples *Samples) GetUint(index int, fieldName string) uint {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return uint(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetByte is a function to retrieve a value of type byte from the samples
func (samples *Samples) GetByte(index int, fieldName string) byte {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return byte(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetRune is a function to retrieve a value of type rune from the samples
func (samples *Samples) GetRune(index int, fieldName string) rune {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return rune(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr))
}

// GetBoolean is a function to retrieve a value of type boolean from the samples
func (samples *Samples) GetBoolean(index int, fieldName string) bool {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return int(C.RTIDDSConnector_getBooleanFromSamples(
		unsafe.Pointer(samples.input.connector.native),
		samples.input.nameCStr,
		C.int(index+1),
		fieldNameCStr),
	) != 0
}

// GetString is a function to retrieve a value of type string from the samples
func (samples *Samples) GetString(index int, fieldName string) string {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))
	return C.GoString((*C.char)(C.RTIDDSConnector_getStringFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), fieldNameCStr)))
}

// GetJSON is a function to retrieve a slice of bytes of a JSON string from the samples
func (samples *Samples) GetJSON(index int) ([]byte, error) {
	jsonCStr := C.RTIDDSConnector_getJSONSample(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1))
	defer C.RTIDDSConnector_freeString((*C.char)(jsonCStr))

	return []byte(C.GoString((*C.char)(jsonCStr))), nil
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

// IsValid is a function to check validity of the element and return a boolean
func (infos *Infos) IsValid(index int) bool {
	memberNameCStr := C.CString("valid_data")
	defer C.free(unsafe.Pointer(memberNameCStr))
	var retVal C.int

	C.RTI_Connector_get_boolean_from_infos(
		unsafe.Pointer(infos.input.connector.native),
		&retVal,
		infos.input.nameCStr,
		C.int(index+1), memberNameCStr,
	)
	return (retVal != 0)
}

// GetSourceTimestamp is a function to get the source timestamp of a sample
func (infos *Infos) GetSourceTimestamp(index int) (int64, error) {
	memberNameCStr := C.CString("source_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var jsonStr string
	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := C.RTI_Connector_get_json_from_infos(
		unsafe.Pointer(infos.input.connector.native),
		infos.input.nameCStr, C.int(index+1),
		memberNameCStr,
		&jsonCStr,
	)
	if retcode != DDSRetCodeOK {
		return 0, errors.New("RTI_Connector_get_json_from_infos failed")
	}

	ts, err := strconv.ParseInt(C.GoString(jsonCStr), 10, 64)
	if err != nil {
		return 0, errors.New("string conversion failed: " + err.Error())
	}

	return ts, nil
}

// GetReceptionTimestamp is a function to get the reception timestamp of a sample
func (infos *Infos) GetReceptionTimestamp(index int) (int64, error) {
	memberNameCStr := C.CString("reception_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var jsonStr string
	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	retcode := C.RTI_Connector_get_json_from_infos(
		unsafe.Pointer(infos.input.connector.native),
		infos.input.nameCStr,
		C.int(index+1),
		memberNameCStr,
		&jsonCStr,
	)
	if retcode != DDSRetCodeOK {
		return 0, errors.New("RTI_Connector_get_json_from_infos failed")
	}

	ts, err := strconv.ParseInt(C.GoString(jsonCStr), 10, 64)
	if err != nil {
		return 0, errors.New("string conversion failed: " + err.Error())
	}

	return ts, err
}

// GetIdentity is a function to get the identity of a writer that sent the sample
func (infos *Infos) GetIdentity(index int) (Identity, error) {
	memberNameCStr := C.CString("identity")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var jsonStr string
	jsonCStr := C.CString(jsonStr)
	defer C.free(unsafe.Pointer(jsonCStr))

	var writerID Identity
	retcode := C.RTI_Connector_get_json_from_infos(
		unsafe.Pointer(infos.input.connector.native),
		infos.input.nameCStr,
		C.int(index+1),
		memberNameCStr,
		&jsonCStr,
	)
	if retcode != DDSRetCodeOK {
		return writerID, errors.New("RTI_Connector_get_json_from_infos failed")
	}

	jsonByte := []byte(C.GoString(jsonCStr))
	if err := json.Unmarshal(jsonByte, &writerID); err != nil {
		return writerID, errors.New("JSON Unmarshal failed: " + err.Error())
	}

	return writerID, nil
}

// GetLength is a function to return the length of the
func (infos *Infos) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getInfosLength(
		unsafe.Pointer(infos.input.connector.native),
		infos.input.nameCStr),
	)
	return length
}
