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
	"strconv"
	"unsafe"
)

/********
* Types *
*********/

// Infos is a sequence of info samples used by an input to read DDS meta data
type Infos struct {
	input *Input
}

// Identity is the structure for identifying
type Identity struct {
	WriterGUID     [16]byte `json:"writer_guid"`
	SequenceNumber int      `json:"sequence_number"`
}

/*******************
* Public Functions *
*******************/

// IsValid is a function to check validity of the element and return a boolean
func (infos *Infos) IsValid(index int) (bool, error) {
	if infos == nil || infos.input == nil || infos.input.connector == nil {
		return false, errors.New("infos, input, or connector is null")
	}
	if index < 0 {
		return false, errors.New("index cannot be negative")
	}

	memberNameCStr := C.CString("valid_data")
	defer C.free(unsafe.Pointer(memberNameCStr))
	var retVal C.int

	retcode := int(C.RTI_Connector_get_boolean_from_infos(unsafe.Pointer(infos.input.connector.native), &retVal, infos.input.nameCStr, C.int(index+1), memberNameCStr))
	err := checkRetcode(retcode)

	return (retVal != 0), err
}

// GetSourceTimestamp is a function to get the source timestamp of a sample
func (infos *Infos) GetSourceTimestamp(index int) (int64, error) {
	tsStr, err := infos.getJSONMember(index, "source_timestamp")
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return ts, nil
}

// GetReceptionTimestamp is a function to get the reception timestamp of a sample
func (infos *Infos) GetReceptionTimestamp(index int) (int64, error) {
	tsStr, err := infos.getJSONMember(index, "reception_timestamp")
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return ts, nil
}

// GetIdentity is a function to get the identity of a writer that sent the sample
func (infos *Infos) GetIdentity(index int) (Identity, error) {

	var writerID Identity

	identityStr, err := infos.getJSONMember(index, "sample_identity")
	if err != nil {
		return writerID, err
	}

	jsonByte := []byte(identityStr)
	err = json.Unmarshal(jsonByte, &writerID)
	if err != nil {
		return writerID, errors.New("JSON Unmarshal failed: " + err.Error())
	}

	return writerID, nil
}

// GetIdentityJSON is a function to get the identity of a writer in JSON
func (infos *Infos) GetIdentityJSON(index int) (string, error) {
	identityStr, err := infos.getJSONMember(index, "sample_identity")
	if err != nil {
		return "", err
	}

	return identityStr, nil
}

// GetRelatedIdentity is a function used for request-reply communications.
func (infos *Infos) GetRelatedIdentity(index int) (Identity, error) {

	var writerID Identity

	identityStr, err := infos.getJSONMember(index, "related_sample_identity")
	if err != nil {
		return writerID, err
	}

	jsonByte := []byte(identityStr)
	err = json.Unmarshal(jsonByte, &writerID)
	if err != nil {
		return writerID, errors.New("JSON Unmarshal failed: " + err.Error())
	}

	return writerID, nil
}

// GetRelatedIdentityJSON is a function used for get related identity in JSON.
func (infos *Infos) GetRelatedIdentityJSON(index int) (string, error) {
	identityStr, err := infos.getJSONMember(index, "related_sample_identity")
	if err != nil {
		return "", err
	}

	return identityStr, nil
}

// GetViewState is a function used to get a view state in string (either "NEW" or "NOT NEW").
func (infos *Infos) GetViewState(index int) (string, error) {
	viewStateStr, err := infos.getJSONMember(index, "view_state")
	if err != nil {
		return "", err
	}

	return viewStateStr, nil
}

// GetInstanceState is a function used to get a instance state in string (one of "ALIVE", "NOT_ALIVE_DISPOSED" or "NOT_ALIVE_NO_WRITERS").
func (infos *Infos) GetInstanceState(index int) (string, error) {
	instanceStateStr, err := infos.getJSONMember(index, "instance_state")
	if err != nil {
		return "", err
	}

	return instanceStateStr, nil
}

// GetSampleState is a function used to get a sample state in string (either "READ" or "NOT_READ").
func (infos *Infos) GetSampleState(index int) (string, error) {
	sampleStateStr, err := infos.getJSONMember(index, "sample_state")
	if err != nil {
		return "", err
	}

	return sampleStateStr, nil
}

// GetLength is a function to return the length of the
func (infos *Infos) GetLength() (int, error) {
	var retVal C.double
	retcode := int(C.RTI_Connector_get_sample_count(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, &retVal))
	err := checkRetcode(retcode)
	return int(retVal), err
}

func (infos *Infos) getJSONMember(index int, memberName string) (string, error) {
	memberNameCStr := C.CString(memberName)
	defer C.free(unsafe.Pointer(memberNameCStr))

	var retValCStr *C.char

	retcode := int(C.RTI_Connector_get_json_from_infos(unsafe.Pointer(infos.input.connector.native), infos.input.nameCStr, C.int(index+1), memberNameCStr, &retValCStr))
	err := checkRetcode(retcode)
	if err != nil {
		return "", err
	}

	retValGoStr := C.GoString(retValCStr)
	C.RTI_Connector_free_string(retValCStr)

	return retValGoStr, nil
}
