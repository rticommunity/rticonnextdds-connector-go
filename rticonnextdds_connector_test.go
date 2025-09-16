package rti

import (
	"math"
	"path"
	"runtime"
	"testing"
	"time"

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

// This test function ensures that an error is raised if an incorrect xml path is passed to the Connector constructor.
func TestInvalidXMLPath(t *testing.T) {
	invalidXMLPath := "invalid/path/to/xml"

	connector, err := NewConnector(participantProfile, invalidXMLPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

// This test function ensures that an error is raised if an invalid participant profile name is passed to the Connector constructor.
func TestInvalidParticipantProfile(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	connector, err := NewConnector(invalidParticipantProfile, xmlPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

// This test function ensures that an error is raised if an invalid xml file is passed to the Connector constructor.
func TestInvalidXMLProfile(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/InvalidXml.xml")

	connector, err := NewConnector(participantProfile, xmlPath)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

// This function tests the correct instantiation of Connector object.
func TestConnectorCreation(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	assert.NotNil(t, connector)
}

// This function tests the correct instantiation of multiple Connector objects in succession.
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

// Tests that it is possible to load two xml files using the url group syntax
func TestLoadMultipleFiles(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath1 := path.Join(path.Dir(curPath), "./test/xml/TestConnector1.xml")
	xmlPath2 := path.Join(path.Dir(curPath), "./test/xml/TestConnector2.xml")

	connector, err := NewConnector("MyParticipantLibrary2::MyParticipant2", xmlPath1+";"+xmlPath2)
	assert.Nil(t, err)
	assert.NotNil(t, connector)

	output, err := connector.GetOutput("MyPublisher2::MySquareWriter2")
	assert.Nil(t, err)
	assert.NotNil(t, output)
}

// Tests that a domain_participant defined in XML alonside participant_qos can be used to create a Connector object.
func TestConnectorCreationWithParticipantQos(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/TestConnector1.xml")

	connector, err := NewConnector("MyParticipantLibrary::ConnectorWithParticipantQos", xmlPath)
	assert.Nil(t, err)
	assert.NotNil(t, connector)
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

	xs := int8(math.MaxInt8)
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

	assert.Nil(t, output.Instance.SetUint("ul", uint(ul)))

	assert.Nil(t, output.Instance.SetRune("l", l))

	assert.Nil(t, output.Instance.SetInt64("ll", ll))

	assert.Nil(t, output.Instance.SetUint64("ull", ull))

	assert.Nil(t, output.Instance.SetFloat32("f", f))

	assert.Nil(t, output.Instance.SetFloat64("d", d))

	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	sampleLength, err := input.Samples.GetLength()
	assert.Nil(t, err)
	assert.Equal(t, sampleLength, 1)

	infoLength, err := input.Infos.GetLength()
	assert.Nil(t, err)
	assert.Equal(t, infoLength, 1)

	valid, err := input.Infos.IsValid(0)
	assert.Nil(t, err)
	assert.Equal(t, valid, true)

	viewState, err := input.Infos.GetViewState(0)
	assert.Nil(t, err)
	assert.Equal(t, viewState, "NEW")

	instanceState, err := input.Infos.GetInstanceState(0)
	assert.Nil(t, err)
	assert.Equal(t, instanceState, "ALIVE")

	sampleState, err := input.Infos.GetSampleState(0)
	assert.Nil(t, err)
	assert.Equal(t, sampleState, "NOT_READ")

	rst, err := input.Samples.GetString(0, "st")
	assert.Nil(t, err)
	assert.Equal(t, rst, st)

	rb, err := input.Samples.GetBoolean(0, "b")
	assert.Nil(t, err)
	assert.Equal(t, rb, b)

	rc, err := input.Samples.GetByte(0, "c")
	assert.Nil(t, err)
	assert.Equal(t, rc, c)

	rxs, err := input.Samples.GetInt8(0, "xs")
	assert.Nil(t, err)
	assert.Equal(t, rxs, xs)

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

	assert.Nil(t, output.ClearMembers())

	// Testing Wait TimeOut
	err = connector.Wait(5)
	t.Log(err)
	assert.NotNil(t, err)

	// Testing Read
	err = output.Write()
	assert.Nil(t, err)
	err = connector.Wait(-1)
	assert.Nil(t, err)
	err = input.Read()
	assert.Nil(t, err)
	rst, err = input.Samples.GetString(0, "st")
	assert.Nil(t, err)
	assert.Equal(t, rst, "")

	id, err := input.Infos.GetIdentity(0)
	assert.Nil(t, err)
	assert.Equal(t, id.SequenceNumber, int(2))
	// UUID can not be checked because it is unique to each run

	ts, err := input.Infos.GetReceptionTimestamp(0)
	assert.Nil(t, err)
	assert.NotNil(t, ts)               // Unique time per each run
	assert.NotNil(t, time.Unix(0, ts)) // Unique time per each run
	assert.NotEqual(t, ts, 0)          // Unique time per each run

	gt, err := input.Infos.GetSourceTimestamp(0)
	assert.Nil(t, err)
	assert.NotNil(t, gt)               // Unique time per each run
	assert.NotNil(t, time.Unix(0, gt)) // Unique time per each run
	assert.NotEqual(t, gt, 0)          // Unique time per each run
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

func TestSimpleMatching(t *testing.T) {
	connector, err := newTestConnector()
	defer connector.Delete()
	input, err := newTestInput(connector)
	output, err := newTestOutput(connector)

	change, err := input.WaitForPublications(2000)
	assert.Nil(t, err)
	assert.Equal(t, change, 1)

	matches, err := input.GetMatchedPublications()
	assert.Nil(t, err)
	assert.Equal(t, matches, "[{\"name\":\"MyWriter\"}]")

	change, err = output.WaitForSubscriptions(2000)
	assert.Nil(t, err)
	assert.Equal(t, change, 1)

	matches, err = output.GetMatchedSubscriptions()
	assert.Nil(t, err)
	assert.Equal(t, matches, "[{\"name\":\"MyReader\"}]")
}
