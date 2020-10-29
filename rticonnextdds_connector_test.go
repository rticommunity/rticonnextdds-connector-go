package rti

import (
	"math"
	"path"
	"runtime"
	"testing"

	"github.com/rticommunity/rticonnextdds-connector-go/types"
	"github.com/stretchr/testify/assert"
)

const (
	participantProfile        = "MyParticipantLibrary::Zero"
	invalidParticipantProfile = "InvalidParticipantProfile"
)

// Helper functions
func newTestConnector() (*Connector, error) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")
	return NewConnector(participantProfile, xmlPath)
}

func newTestInput(connector *Connector) (*Input, error) {
	return connector.GetInput("MySubscriber::MyReader")
}

func newTestOutput(connector *Connector) (*Output, error) {
	return connector.GetOutput("MyPublisher::MyWriter")
}

// Connector test
func TestInvalidXMLPath(t *testing.T) {
	invalidXMLPath := "invalid/path/to/xml"

	connector, err := NewConnector(participantProfile, invalidXMLPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestInvalidParticipantProfile(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	connector, err := NewConnector(invalidParticipantProfile, xmlPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestMultipleConnectorCreation(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")
	var (
		connectors [5]*Connector
		err        error
	)
	for i := 0; i < 5; i++ {
		connectors[i], err = NewConnector(participantProfile, xmlPath)
		assert.Nil(t, err)
		assert.NotNil(t, connectors[i])
	}

	for i := 0; i < 5; i++ {
		assert.Nil(t, connectors[i].Delete())
	}
}

func TestConnectorDeletion(t *testing.T) {
	var nullConnector *Connector
	assert.NotNil(t, nullConnector.Delete())
}

// Input tests
func TestInvalidDR(t *testing.T) {
	invalidReaderName := "invalidDR"

	connector, err := newTestConnector()
	assert.Nil(t, err)
	assert.NotNil(t, connector)
	input, err := connector.GetInput(invalidReaderName)
	assert.Nil(t, input)
	assert.NotNil(t, err)
}

func TestCreateDR(t *testing.T) {
	readerName := "MySubscriber::MyReader"

	connector, err := newTestConnector()
	assert.Nil(t, err)
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
	assert.NotNil(t, nullConnector.Wait(-1))
}

// Output tests
func TestInvalidWriter(t *testing.T) {
	invalidWriterName := "invalidWriter"

	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()
	output, err := connector.GetOutput(invalidWriterName)
	assert.Nil(t, output)
	assert.NotNil(t, err)
}

func TestCreateWriter(t *testing.T) {
	writerName := "MyPublisher::MyWriter"

	connector, err := newTestConnector()
	assert.Nil(t, err)
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
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()
	input, err := newTestInput(connector)
	assert.Nil(t, err)
	output, err := newTestOutput(connector)
	assert.Nil(t, err)

	// Take any pre-existing samples from cache
	assert.Nil(t, input.Take())

	xs := int8(math.MaxInt8)
	s := int16(math.MaxInt16)
	us := uint16(math.MaxUint16)
	l := int32(math.MaxInt32)
	ll := int64(123)
	ul := uint32(math.MaxUint32)
	ull := uint64(456)
	f := float32(math.MaxFloat32)
	d := float64(math.MaxFloat64)

	c := byte('A')
	b := true
	st := "test"

	assert.Nil(t, output.Instance.SetInt8("xs", xs))
	assert.Nil(t, output.Instance.SetUint8("c", c))
	assert.Nil(t, output.Instance.SetByte("c", c))
	assert.Nil(t, output.Instance.SetString("st", st))
	assert.Nil(t, output.Instance.SetBoolean("b", b))
	assert.Nil(t, output.Instance.SetInt16("s", s))
	assert.Nil(t, output.Instance.SetUint16("us", us))
	assert.Nil(t, output.Instance.SetInt32("l", l))
	assert.Nil(t, output.Instance.SetUint32("ul", ul))
	assert.Nil(t, output.Instance.SetInt("l", int(l)))
	assert.Nil(t, output.Instance.SetInt64("ll", ll))
	assert.Nil(t, output.Instance.SetUint64("ull", ull))
	assert.Nil(t, output.Instance.SetUint("ul", uint(ul)))
	assert.Nil(t, output.Instance.SetRune("l", rune(l))) //nolint //this is because "l" rune is still an int32
	assert.Nil(t, output.Instance.SetFloat32("f", f))
	assert.Nil(t, output.Instance.SetFloat64("d", d))

	assert.Nil(t, output.Write())
	assert.Nil(t, connector.Wait(-1))
	assert.Nil(t, input.Take())

	sampleLength := input.Samples.GetLength()
	assert.Equal(t, sampleLength, 1)

	infoLength := input.Infos.GetLength()
	assert.Equal(t, infoLength, 1)

	valid := input.Infos.IsValid(0)
	assert.Equal(t, valid, true)

	assert.Equal(t, input.Samples.GetString(0, "st"), st)
	assert.Equal(t, input.Samples.GetBoolean(0, "b"), b)

	assert.Equal(t, input.Samples.GetUint8(0, "c"), c)
	assert.Equal(t, input.Samples.GetByte(0, "c"), c)
	assert.Equal(t, input.Samples.GetInt8(0, "xs"), xs)
	assert.Equal(t, input.Samples.GetInt16(0, "s"), s)
	assert.Equal(t, input.Samples.GetUint16(0, "us"), us)
	assert.Equal(t, input.Samples.GetInt32(0, "l"), l)
	assert.Equal(t, input.Samples.GetInt(0, "l"), int(l))
	assert.Equal(t, input.Samples.GetInt64(0, "ll"), ll)
	assert.Equal(t, input.Samples.GetUint(0, "ul"), uint(ul))
	assert.Equal(t, input.Samples.GetRune(0, "l"), rune(l)) //nolint //this is because "l" rune is still an int32
	assert.Equal(t, input.Samples.GetUint32(0, "ul"), ul)
	assert.Equal(t, input.Samples.GetUint64(0, "ull"), ull)
	assert.Equal(t, input.Samples.GetFloat32(0, "f"), f)
	assert.Equal(t, input.Samples.GetFloat64(0, "d"), d)

	assert.Nil(t, output.ClearMembers())

	// Testing Wait TimeOut
	assert.NotNil(t, connector.Wait(5))

	// Testing Read
	assert.Nil(t, output.Write())
	assert.Nil(t, connector.Wait(-1))
	assert.Nil(t, input.Read())

	id, err := input.Infos.GetIdentity(0)
	assert.Nil(t, err)
	assert.Equal(t, id.SequenceNumber, uint(2))
	// UUID can not be checked because it is unique to each run

	assert.NotEqual(t, input.Samples.GetString(0, "st"), st)
}

func TestJSON(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()
	input, err := newTestInput(connector)
	assert.Nil(t, err)
	output, err := newTestOutput(connector)
	assert.Nil(t, err)

	var outputTestData types.Test
	outputTestData.St = "output_test"
	assert.Nil(t, output.Instance.Set(&outputTestData))
	assert.Nil(t, output.Write())
	assert.Nil(t, connector.Wait(-1))
	assert.Nil(t, input.Take())

	var inputTestData types.Test
	assert.Nil(t, input.Samples.Get(0, &inputTestData))
	assert.Equal(t, inputTestData.St, outputTestData.St)
}
