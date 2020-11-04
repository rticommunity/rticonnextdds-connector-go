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
	native  *C.RTI_Connector
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
	// DDSRetCodeNoData is a Return Code from CGO for no data return
	DDSRetCodeNoData = 11
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

	output.native = C.RTI_Connector_get_datawriter(unsafe.Pointer(connector.native), output.nameCStr)
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

// checkRetcode is a function to check return code
func checkRetcode(retcode int) error {
	switch retcode {
	case DDSRetCodeOK:
	case DDSRetCodeNoData:
		return errors.New("DDS Exceptrion: No Data")
	case DDSRetCodeTimeout:
		return errors.New("DDS Exception: Timeout")
	default:
		return errors.New("DDS Exception: " + C.GoString((*C.char)(C.RTI_Connector_get_last_error_message)))
	}
	return nil
}

// getNumber is a function to return a number in double from a sample
func (samples *Samples) getNumber(index int, fieldName string, retVal *C.double) error {
	fieldNameCStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldNameCStr))

	retcode := int(C.RTI_Connector_get_number_from_sample(unsafe.Pointer(samples.input.connector.native), retVal, samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	return checkRetcode(retcode)
}

// GetUint8 is a function to retrieve a value of type uint8 from the samples

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

	connector.native = C.RTI_Connector_new(configNameCStr, urlCStr, nil)
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

	C.RTI_Connector_delete(connector.native)
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

	retcode := int(C.RTI_Connector_wait_for_data(unsafe.Pointer(connector.native), C.int(timeoutMs)))
	return checkRetcode(retcode)
}

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

	var retVal *C.char

	retcode := int(C.RTI_Connector_get_string_from_sample(unsafe.Pointer(samples.input.connector.native), &retVal, samples.input.nameCStr, C.int(index+1), fieldNameCStr))
	err := checkRetcode(retcode)

	return C.GoString(retVal), err
}

// GetJSON is a function to retrieve a slice of bytes of a JSON string from the samples
func (samples *Samples) GetJSON(index int) ([]byte, error) {
	var retVal *C.char

	retcode := int(C.RTI_Connector_get_json_sample(unsafe.Pointer(samples.input.connector.native), samples.input.nameCStr, C.int(index+1), &retVal))
	err := checkRetcode(retcode)

	return []byte(C.GoString(retVal)), err
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
func (infos *Infos) IsValid(index int) (bool, error) {
	memberNameCStr := C.CString("valid_data")
	defer C.free(unsafe.Pointer(memberNameCStr))
	var retVal C.int

	retcode := int(C.RTI_Connector_get_boolean_from_infos(unsafe.Pointer(infos.input.connector.native), &retVal, infos.input.nameCStr, C.int(index+1), memberNameCStr))
	err := checkRetcode(retcode)

	return (retVal != 0), err
}

// GetSourceTimestamp is a function to get the source timestamp of a sample
func (infos *Infos) GetSourceTimestamp(index int) (int64, error) {
	memberNameCStr := C.CString("source_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var retVal *C.char

	retcode := int(C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &retVal))
	err := checkRetcode(retcode)
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(C.GoString(retVal), 10, 64)
	if err != nil {
		return 0, err
	}

	return ts, nil
}

// GetReceptionTimestamp is a function to get the reception timestamp of a sample
func (infos *Infos) GetReceptionTimestamp(index int) (int64, error) {
	memberNameCStr := C.CString("reception_timestamp")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var retVal *C.char

	retcode := int(C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &retVal))
	err := checkRetcode(retcode)
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(C.GoString(retVal), 10, 64)
	if err != nil {
		return 0, err
	}

	return ts, err
}

// GetIdentity is a function to get the identity of a writer that sent the sample
func (infos *Infos) GetIdentity(index int) (Identity, error) {
	memberNameCStr := C.CString("identity")
	defer C.free(unsafe.Pointer(memberNameCStr))

	var retVal *C.char
	var writerID Identity

	retcode := int(C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &retVal))
	err := checkRetcode(retcode)
	if err != nil {
		return writerID, err
	}

	jsonByte := []byte(C.GoString(retVal))
	err = json.Unmarshal(jsonByte, &writerID)
	if err != nil {
		return writerID, errors.New("JSON Unmarshal failed: " + err.Error())
	}

	return writerID, nil
}

// GetLength is a function to return the length of the
func (infos *Infos) GetLength() (int, error) {
	var retVal C.double
	retcode := int(C.RTI_Connector_get_sample_count(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, &retVal))
	err := checkRetcode(retcode)
	return int(retVal), err
}
