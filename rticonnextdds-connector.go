/*****************************************************************************
*    (c) 2005-2015 Copyright, Real-Time Innovations, All rights reserved.    *
*                                                                            *
*  RTI grants Licensee a license to use, modify, compile, and create         *
*  derivative works of the Software.  Licensee has the right to distribute   *
*  object form only for use with RTI products. The Software is provided      *
*  "as is", with no warranty of any type, including any warranty for fitness *
*  for any purpose. RTI is under no obligation to maintain or support the    *
*  Software.  RTI shall not be liable for any incidental or consequential    *
*  damages arising out of the use or inability to use the software.          *
*                                                                            *
*****************************************************************************/

// Package rti implements functions of RTI Connector for Connext DDS in Go
package rti

// #cgo darwin CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -m64
// #cgo linux,amd64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
// #cgo linux,arm CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX
// #cgo darwin LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Darwin16clang8.0 -lrtiddsconnector -ldl -lm -lpthread
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

type Connector struct {
	native  *C.struct_RTIDDSConnector
	Inputs  []Input
	Outputs []Output
}

type Output struct {
	native     unsafe.Pointer // a pointer to a native DataWriter
	connector  *Connector
	name       string // name of the native DataWriter
	name_c_str *C.char
	Instance   *Instance
}

type Instance struct {
	output *Output
}

type Input struct {
	native     unsafe.Pointer // a pointer to a native DataReader
	connector  *Connector
	name       string // name of the native DataReader
	name_c_str *C.char
	Samples    *Samples
	Infos      *Infos
}

type Samples struct {
	input *Input
}

type Infos struct {
	input *Input
}

/********************
* Private Functions *
********************/

func newInstance(output *Output) (instance *Instance, err error) {
	// Error checking for the output is skipped because it was already checked

	instance = new(Instance)
	instance.output = output

	return instance, nil
}

func newOutput(connector *Connector, output_name string) (output *Output, err error) {
	// Error checking for the connector is skipped because it was already checked

	output = new(Output)
	output.connector = connector

	output.name_c_str = C.CString(output_name)

	output.native = C.RTIDDSConnector_getWriter(unsafe.Pointer(connector.native), output.name_c_str)
	if output.native == nil {
		err = errors.New("Invalid Publication::DataWriter name")
		return nil, err
	}
	output.name = output_name
	output.Instance, err = newInstance(output)
	if err != nil {
		err = errors.New("newInstance error")
		return nil, err
	}

	connector.Outputs = append(connector.Outputs, *output)

	return output, nil
}

func newInput(connector *Connector, input_name string) (input *Input, err error) {
	// Error checking for the connector is skipped because it was already checked

	input = new(Input)
	input.connector = connector

	input.name_c_str = C.CString(input_name)

	input.native = C.RTIDDSConnector_getReader(unsafe.Pointer(connector.native), input.name_c_str)
	if input.native == nil {
		err = errors.New("Invalid Subscription::DataReader name")
		return nil, err
	}
	input.name = input_name
	input.Samples, err = newSamples(input)
	if err != nil {
		err = errors.New("newSamples error")
		return nil, err
	}
	input.Infos, err = newInfos(input)
	if err != nil {
		err = errors.New("newInfos error")
		return nil, err
	}

	connector.Inputs = append(connector.Inputs, *input)

	return input, nil
}

func newSamples(input *Input) (samples *Samples, err error) {
	// Error checking for the input is skipped because it was already checked

	// TODO - need to check "new" could fail. if so, I need to check after calling this
	samples = new(Samples)
	samples.input = input
	return samples, nil
}

func newInfos(input *Input) (infos *Infos, err error) {
	// Error checking for the input is skipped because it was already checked

	infos = new(Infos)
	infos.input = input
	return infos, nil
}

/*******************
* Public Functions *
*******************/

func NewConnector(config_name string, file_name string) (connector *Connector, err error) {
	connector = new(Connector)

	config_name_c_str := C.CString(config_name)
	defer C.free(unsafe.Pointer(config_name_c_str))
	file_name_c_str := C.CString(file_name)
	defer C.free(unsafe.Pointer(file_name_c_str))

	connector.native = C.RTIDDSConnector_new(config_name_c_str, file_name_c_str, nil)
	if connector.native == nil {
		err = errors.New("Invalid participant profile, xml path or xml profile")
		return nil, err
	}

	return connector, nil
}

func (connector *Connector) Delete() (err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return err
	}

	// Delete memory allocated in C layer
	for _, input := range connector.Inputs {
		C.free(unsafe.Pointer(input.name_c_str))
	}
	for _, output := range connector.Outputs {
		C.free(unsafe.Pointer(output.name_c_str))
	}

	C.RTIDDSConnector_delete(connector.native)
	connector.native = nil

	return nil
}

func (connector *Connector) GetOutput(output_name string) (output *Output, err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return nil, err
	}

	output, err = newOutput(connector, output_name)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (connector *Connector) GetInput(input_name string) (input *Input, err error) {
	if connector == nil {
		err = errors.New("Connector is null")
		return nil, err
	}

	input, err = newInput(connector, input_name)
	if err != nil {
		return nil, err
	}
	return input, nil
}

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
	}
	return nil
}

func (output *Output) Write() error {
	// The C function does not return errors. In the futurue, we will check erros this when supported in the C layer
	// CON-24 (for more information)
	C.RTIDDSConnector_write(unsafe.Pointer(output.connector.native), output.name_c_str, nil)
	return nil
}

func (output *Output) ClearMembers() error {
	// The C function does not return errors. In the futurue, we will check erros when supported in C the C layer
	C.RTIDDSConnector_clear(unsafe.Pointer(output.connector.native), output.name_c_str)
	return nil
}

func (instance *Instance) SetUint8(field_name string, value uint8) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetUint16(field_name string, value uint16) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetUint32(field_name string, value uint32) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetUint64(field_name string, value uint64) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetInt8(field_name string, value int8) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetInt16(field_name string, value int16) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetInt32(field_name string, value int32) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetInt64(field_name string, value int64) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetUint(field_name string, value uint) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetInt(field_name string, value int) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetFloat32(field_name string, value float32) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetFloat64(field_name string, value float64) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetString(field_name string, value string) error {

	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value_c_str := C.CString(value)
	defer C.free(unsafe.Pointer(value_c_str))

	C.RTIDDSConnector_setStringIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, value_c_str)

	return nil
}

func (instance *Instance) SetByte(field_name string, value byte) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetRune(field_name string, value rune) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	C.RTIDDSConnector_setNumberIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.double(value))
	return nil
}

func (instance *Instance) SetBoolean(field_name string, value bool) error {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	var int_value int
	if value == true {
		int_value = 1
	} else {
		int_value = 0
	}
	C.RTIDDSConnector_setBooleanIntoSamples(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, field_name_c_str, C.int(int_value))
	return nil
}

func (instance *Instance) SetJson(json []byte) error {
	json_c_str := C.CString(string(json))
	defer C.free(unsafe.Pointer(json_c_str))

	C.RTIDDSConnector_setJSONInstance(unsafe.Pointer(instance.output.connector.native), instance.output.name_c_str, json_c_str)
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
	C.RTIDDSConnector_read(unsafe.Pointer(input.connector.native), input.name_c_str)
	return nil
}

func (input *Input) Take() (err error) {
	if input == nil {
		err = errors.New("Input is null")
		return err
	}
	// The C function does not return errors. In the futurue, we will update this when supported in the C layer
	C.RTIDDSConnector_take(unsafe.Pointer(input.connector.native), input.name_c_str)
	return nil
}

func (samples *Samples) GetLength() (length int) {
	length = int(C.RTIDDSConnector_getSamplesLength(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str))
	return length
}

func (samples *Samples) GetUint8(index int, field_name string) (value uint8) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = uint8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetUint16(index int, field_name string) (value uint16) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = uint16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetUint32(index int, field_name string) (value uint32) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = uint32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetUint64(index int, field_name string) (value uint64) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = uint64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetInt8(index int, field_name string) (value int8) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = int8(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetInt16(index int, field_name string) (value int16) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = int16(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetInt32(index int, field_name string) (value int32) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = int32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetInt64(index int, field_name string) (value int64) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = int64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetFloat32(index int, field_name string) (value float32) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = float32(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetFloat64(index int, field_name string) (value float64) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = float64(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetInt(index int, field_name string) (value int) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = int(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetUint(index int, field_name string) (value uint) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = uint(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetByte(index int, field_name string) (value byte) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = byte(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetRune(index int, field_name string) (value rune) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = rune(C.RTIDDSConnector_getNumberFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	return value
}

func (samples *Samples) GetBoolean(index int, field_name string) bool {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value := int(C.RTIDDSConnector_getBooleanFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str))
	if value != 0 {
		return true
	} else {
		return false
	}
}

func (samples *Samples) GetString(index int, field_name string) (value string) {
	field_name_c_str := C.CString(field_name)
	defer C.free(unsafe.Pointer(field_name_c_str))

	value = C.GoString((*C.char)(C.RTIDDSConnector_getStringFromSamples(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1), field_name_c_str)))
	return value
}

func (samples *Samples) GetJson(index int) (json []byte, e error) {
	json_c_str := C.RTIDDSConnector_getJSONSample(unsafe.Pointer(samples.input.connector.native), samples.input.name_c_str, C.int(index+1))
	defer C.RTIDDSConnector_freeString((*C.char)(json_c_str))

	json = []byte(C.GoString((*C.char)(json_c_str)))

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
	input_name_c_str := C.CString(infos.input.name)
	defer C.free(unsafe.Pointer(input_name_c_str))
	member_name_c_str := C.CString("valid_data")
	defer C.free(unsafe.Pointer(member_name_c_str))

	if int(C.RTIDDSConnector_getBooleanFromInfos(unsafe.Pointer(infos.input.connector.native), input_name_c_str, C.int(index+1), member_name_c_str)) != 0 {
		valid = true
	} else {
		valid = false
	}
	return valid
}

func (infos *Infos) GetLength() (length int) {
	input_name_c_str := C.CString(infos.input.name)
	defer C.free(unsafe.Pointer(input_name_c_str))

	length = int(C.RTIDDSConnector_getInfosLength(unsafe.Pointer(infos.input.connector.native), input_name_c_str))
	return length
}
