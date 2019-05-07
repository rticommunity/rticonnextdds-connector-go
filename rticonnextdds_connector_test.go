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
	_, cur_path, _, _ := runtime.Caller(0)
	xml_path := path.Join(path.Dir(cur_path), "./test/xml/Test.xml")
	participant_profile := "MyParticipantLibrary::Zero"
	connector, _ = NewConnector(participant_profile, xml_path)
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
	participant_profile := "MyParticipantLibrary::Zero"
	invalid_xml_path := "invalid/path/to/xml"

	connector, err := NewConnector(participant_profile, invalid_xml_path)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestInvalidParticipantProfile(t *testing.T) {
	_, cur_path, _, _ := runtime.Caller(0)
	xml_path := path.Join(path.Dir(cur_path), "./test/xml/Test.xml")
	invalid_participant_profile := "InvalidParticipantProfile"

	connector, err := NewConnector(invalid_participant_profile, xml_path)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestMultipleConnectorCreation(t *testing.T) {
	_, cur_path, _, _ := runtime.Caller(0)
	xml_path := path.Join(path.Dir(cur_path), "./test/xml/Test.xml")
	participant_profile := "MyParticipantLibrary::Zero"
	var connectors [5]*Connector
	for i := 0; i < 5; i++ {
		connectors[i], _ = NewConnector(participant_profile, xml_path)
		assert.NotNil(t, connectors[i])
	}

	for i := 0; i < 5; i++ {
		err := connectors[i].Delete()
		assert.Nil(t, err)
	}
}

func TestConnectorDeletion(t *testing.T) {
	var null_connector *Connector
	err := null_connector.Delete()
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

	var null_connector *Connector
	input, err = null_connector.GetInput(readerName)
	assert.Nil(t, input)
	assert.NotNil(t, err)
	err = null_connector.Wait(-1)
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

	var null_connector *Connector
	output, err = null_connector.GetOutput(writerName)
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
	//ll := int64(math.MaxInt64)
	//ull := uint64(math.MaxUint64)
	f := float32(math.MaxFloat32)
	d := float64(math.MaxFloat64)

	c := byte('A')
	b := true
	st := "test"

	output.Instance.SetUint8("c", c)
	output.Instance.SetByte("c", c)
	output.Instance.SetString("st", st)
	output.Instance.SetBoolean("b", b)
	output.Instance.SetInt16("s", s)
	output.Instance.SetUint16("us", us)
	output.Instance.SetInt32("l", l)
	output.Instance.SetUint32("ul", ul)
	output.Instance.SetInt("l", int(l))
	output.Instance.SetUint("ul", uint(ul))
	output.Instance.SetRune("l", rune(l))
	//output.Instance.SetInt64("ll", ll)
	//output.Instance.SetUint64("ull", ull)
	output.Instance.SetFloat32("f", f)
	output.Instance.SetFloat64("d", d)

	output.Write()

	err := connector.Wait(-1)
	assert.Nil(t, err)
	input.Take()

	sample_length := input.Samples.GetLength()
	assert.Equal(t, sample_length, 1)

	info_length := input.Infos.GetLength()
	assert.Equal(t, info_length, 1)

	valid := input.Infos.IsValid(0)
	assert.Equal(t, valid, true)

	assert.Equal(t, input.Samples.GetString(0, "st"), st)
	assert.Equal(t, input.Samples.GetBoolean(0, "b"), b)

	assert.Equal(t, input.Samples.GetUint8(0, "c"), c)
	assert.Equal(t, input.Samples.GetByte(0, "c"), c)
	assert.Equal(t, input.Samples.GetInt16(0, "s"), s)
	assert.Equal(t, input.Samples.GetUint16(0, "us"), us)
	assert.Equal(t, input.Samples.GetInt32(0, "l"), l)
	assert.Equal(t, input.Samples.GetInt(0, "l"), int(l))
	assert.Equal(t, input.Samples.GetUint(0, "ul"), uint(ul))
	assert.Equal(t, input.Samples.GetRune(0, "l"), rune(l))
	assert.Equal(t, input.Samples.GetUint32(0, "ul"), ul)
	//assert.Equal(t, input.Samples.GetInt64(0, "ll"), ll)
	//assert.Equal(t, input.Samples.GetUint64(0, "ull"), ull)
	assert.Equal(t, input.Samples.GetFloat32(0, "f"), f)
	assert.Equal(t, input.Samples.GetFloat64(0, "d"), d)

	output.ClearMembers()
	output.Write()
	connector.Wait(-1)
	input.Read()
	assert.NotEqual(t, input.Samples.GetString(0, "st"), st)

	// Testing Wait TimeOut
	err = connector.Wait(5)
	assert.NotNil(t, err)
}

func TestJSON(t *testing.T) {
	connector := newTestConnector()
	defer connector.Delete()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	var output_test_data types.Test
	output_test_data.St = "test"
	output.Instance.Set(&output_test_data)
	output.Write()

	err := connector.Wait(-1)
	assert.Nil(t, err)
	input.Take()

	var input_test_data types.Test
	input.Samples.Get(0, &input_test_data)

	assert.Equal(t, input_test_data.St, output_test_data.St)

}
