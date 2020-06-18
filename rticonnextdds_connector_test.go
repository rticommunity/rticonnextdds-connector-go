package rti

import (
	"github.com/rticommunity/rticonnextdds-connector-go/types"
	"github.com/stretchr/testify/assert"
	"math"
	"path"
	"runtime"
	"testing"
)

// Helper functions
func newTestConnector() (connector *Connector) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")
	participantProfile := "MyParticipantLibrary::Zero"
	connector, _ = NewConnector(participantProfile, xmlPath)
	return connector
}

func newTestInput(connector *Connector) (input *Input) {
	input, _ = connector.GetInput("MySubscriber::MyReader")
	return input
}

func newTestOutput(connector *Connector) (output *Output) {
	output, _ = connector.GetOutput("MyPublisher::MyWriter")
	return output
}

// Connector test
func TestInvalidXMLPath(t *testing.T) {
	participantProfile := "MyParticipantLibrary::Zero"
	invalidXMLPath := "invalid/path/to/xml"

	connector, err := NewConnector(participantProfile, invalidXMLPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestInvalidParticipantProfile(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")
	invalidParticipantProfile := "InvalidParticipantProfile"

	connector, err := NewConnector(invalidParticipantProfile, xmlPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestMultipleConnectorCreation(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")
	participantProfile := "MyParticipantLibrary::Zero"
	var connectors [5]*Connector
	for i := 0; i < 5; i++ {
		connectors[i], _ = NewConnector(participantProfile, xmlPath)
		assert.NotNil(t, connectors[i])
	}

	for i := 0; i < 5; i++ {
		err := connectors[i].Delete()
		assert.Nil(t, err)
	}
}

func TestConnectorDeletion(t *testing.T) {
	var nullConnector *Connector
	err := nullConnector.Delete()
	assert.NotNil(t, err)
}

// Input tests
func TestInvalidDR(t *testing.T) {
	invalidReaderName := "invalidDR"

	connector := newTestConnector()
	input, err := connector.GetInput(invalidReaderName)
	assert.Nil(t, input)
	assert.NotNil(t, err)
}

func TestCreateDR(t *testing.T) {
	readerName := "MySubscriber::MyReader"

	connector := newTestConnector()
	defer connector.Delete()
	input, err := connector.GetInput(readerName)
	assert.NotNil(t, input)
	assert.NotNil(t, input.Samples)
	assert.NotNil(t, input.Infos)
	assert.Nil(t, err)

	var nullConnector *Connector
	input, err = nullConnector.GetInput(readerName)
	assert.Nil(t, input)
	assert.NotNil(t, err)
	err = nullConnector.Wait(-1)
	assert.NotNil(t, err)
}

// Output tests
func TestInvalidWriter(t *testing.T) {
	invalidWriterName := "invalidWriter"

	connector := newTestConnector()
	defer connector.Delete()
	output, err := connector.GetOutput(invalidWriterName)
	assert.Nil(t, output)
	assert.NotNil(t, err)
}

func TestCreateWriter(t *testing.T) {
	writerName := "MyPublisher::MyWriter"

	connector := newTestConnector()
	defer connector.Delete()
	output, err := connector.GetOutput(writerName)
	assert.NotNil(t, output)
	assert.NotNil(t, output.Instance)
	assert.Nil(t, err)

	var nullConnector *Connector
	output, err = nullConnector.GetOutput(writerName)
	assert.Nil(t, output)
	assert.NotNil(t, err)
}

// Data flow tests
func TestDataFlow(t *testing.T) {
	connector := newTestConnector()
	defer connector.Delete()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	// Take any pre-existing samples from cache
	input.Take()

	s := int16(math.MaxInt16)
	us := uint16(math.MaxUint16)
	l := int32(math.MaxInt32)
	ul := uint32(math.MaxUint32)
	// an integral value larger than 2^53 can only be retrieved or set with a dictionary/json string
	ll := int64(math.Pow(2, 52))
	ull := uint64(math.Pow(2, 52))
	f := float32(math.MaxFloat32)
	d := float64(math.MaxFloat64)

	c := byte('A')
	b := true
	st := "test"

	err := output.Instance.SetUint8("c", c)
	assert.Nil(t, err)

	err = output.Instance.SetByte("c", c)
	assert.Nil(t, err)

	err = output.Instance.SetString("st", st)
	assert.Nil(t, err)

	err = output.Instance.SetBoolean("b", b)
	assert.Nil(t, err)

	err = output.Instance.SetInt16("s", s)
	assert.Nil(t, err)

	err = output.Instance.SetUint16("us", us)
	assert.Nil(t, err)

	err = output.Instance.SetInt32("l", l)
	assert.Nil(t, err)

	err = output.Instance.SetUint32("ul", ul)
	assert.Nil(t, err)

	err = output.Instance.SetInt("l", int(l))
	assert.Nil(t, err)

	err = output.Instance.SetUint("ul", uint(ul))
	assert.Nil(t, err)

	err = output.Instance.SetRune("l", rune(l))
	assert.Nil(t, err)

	err = output.Instance.SetInt64("ll", ll)
	assert.Nil(t, err)

	err = output.Instance.SetUint64("ull", ull)
	assert.Nil(t, err)

	err = output.Instance.SetFloat32("f", f)
	assert.Nil(t, err)

	err = output.Instance.SetFloat64("d", d)
	assert.Nil(t, err)

	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	sampleLength, err := input.Samples.GetLength()
	assert.Nil(t, err)
	assert.Equal(t, sampleLength, 1)

	infoLength, err:= input.Infos.GetLength()
	assert.Nil(t, err)
	assert.Equal(t, infoLength, 1)

	valid, err := input.Infos.IsValid(0)
	assert.Nil(t, err)
	assert.Equal(t, valid, true)

	rst, err := input.Samples.GetString(0, "st")
	assert.Nil(t, err)
	assert.Equal(t, rst, st)

	rb, err := input.Samples.GetBoolean(0, "b")
	assert.Nil(t, err)
	assert.Equal(t, rb, b)

	rc, err := input.Samples.GetByte(0, "c")
	assert.Nil(t, err)
	assert.Equal(t, rc, c)

	rc, err = input.Samples.GetUint8(0, "c")
	assert.Nil(t, err)
	assert.Equal(t, rc, c)

	rs, err := input.Samples.GetInt16(0, "s")
	assert.Nil(t, err)
	assert.Equal(t, rs, s)

	rus, err := input.Samples.GetUint16(0, "us")
	assert.Nil(t, err)
	assert.Equal(t, rus, us)

	rl, err := input.Samples.GetInt32(0, "l")
	assert.Nil(t, err)
	assert.Equal(t, rl, l)

	rul, err := input.Samples.GetUint32(0, "ul")
	assert.Nil(t, err)
	assert.Equal(t, rul, ul)

	ri, err := input.Samples.GetInt(0, "l")
	assert.Nil(t, err)
	assert.Equal(t, ri, int(l))

	rui, err := input.Samples.GetUint(0, "ul")
	assert.Nil(t, err)
	assert.Equal(t, rui, uint(ul))

	rl, err = input.Samples.GetRune(0, "l")
	assert.Nil(t, err)
	assert.Equal(t, rl, l)

	rll, err := input.Samples.GetInt64(0, "ll")
	assert.Nil(t, err)
	assert.Equal(t, rll, ll)

	rull, err := input.Samples.GetUint64(0, "ull")
	assert.Nil(t, err)
	assert.Equal(t, rull, ull)

	rf, err := input.Samples.GetFloat32(0, "f")
	assert.Nil(t, err)
	assert.Equal(t, rf, f)

	rd, err := input.Samples.GetFloat64(0, "d")
	assert.Nil(t, err)
	assert.Equal(t, rd, d)

	output.ClearMembers()

	// Testing Wait TimeOut
	err = connector.Wait(5)
	assert.NotNil(t, err)

	// Testing Read
	output.Write()
	connector.Wait(-1)
	input.Read()
	rst, err = input.Samples.GetString(0, "st")
	assert.Nil(t, err)
	assert.Equal(t, rst, "")
}

func TestJSON(t *testing.T) {
	connector := newTestConnector()
	defer connector.Delete()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	var outputTestData types.Test
	outputTestData.St = "test"
	output.Instance.Set(&outputTestData)
	output.Write()

	err := connector.Wait(-1)
	assert.Nil(t, err)
	input.Take()

	var inputTestData types.Test
	input.Samples.Get(0, &inputTestData)

	assert.Equal(t, inputTestData.St, outputTestData.St)
}
