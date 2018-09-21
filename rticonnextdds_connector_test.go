package rti

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"runtime"
	"path"
)

// Helper functions
func newTestConnector()(connector *Connector) {
        _, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./test/xml/ShapeExample.xml")
        participant_profile := "MyParticipantLibrary::Zero"
        connector, _ = NewConnector(participant_profile, xml_path)
        return connector
}

func deleteTestConnector(connector *Connector){
        connector.Delete()
}

func newTestInput(connector *Connector)(input *Input) {
	input, _ = connector.GetInput("MySubscriber::MySquareReader")
	return input
}

func newTestOutput(connector *Connector)(output *Output) {
	output, _ = connector.GetOutput("MyPublisher::MySquareWriter")
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
        xml_path := path.Join(path.Dir(cur_path), "./test/xml/ShapeExample.xml")
	invalid_participant_profile := "InvalidParticipantProfile"

        connector, err := NewConnector(invalid_participant_profile, xml_path)
        assert.Nil(t, connector)
        assert.NotNil(t, err)
}

func TestMultipleConnectorCreation(t *testing.T) {
	_, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./test/xml/ShapeExample.xml")
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
        readerName := "MySubscriber::MySquareReader"

        connector := newTestConnector()
        input, err := connector.GetInput(readerName)
        assert.NotNil(t, input)
        assert.NotNil(t, input.Samples)
        assert.NotNil(t, input.Infos)
        assert.Nil(t, err)
        deleteTestConnector(connector)

	var null_connector *Connector
	input, err = null_connector.GetInput(readerName)
	assert.Nil(t, input)
	assert.NotNil(t, err)
}

// Output tests
func TestInvalidWriter(t *testing.T) {
        invalidWriterName := "invalidWriter"

        connector := newTestConnector()
        output, err := connector.GetOutput(invalidWriterName)
        assert.Nil(t, output)
        assert.NotNil(t, err)
}

func TestCreateWriter(t *testing.T) {
        writerName := "MyPublisher::MySquareWriter"

        connector := newTestConnector()
        output, err := connector.GetOutput(writerName)
        assert.NotNil(t, output)
        assert.NotNil(t, output.Instance)
        assert.Nil(t, err)
        deleteTestConnector(connector)

        var null_connector *Connector
        output, err = null_connector.GetOutput(writerName)
        assert.Nil(t, output)
        assert.NotNil(t, err)
}

// Data flow tests
func TestDataFlow(t *testing.T) {
	connector := newTestConnector()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	// Take any pre-existing samples from cache
	input.Take()

        output.Instance.SetInt("x", 1)
        output.Instance.SetInt("y", 1)
        output.Instance.SetBoolean("z", true)
        output.Instance.SetInt("shapesize", 5)
        output.Instance.SetString("color", "BLUE")
	output.Write()

	connector.Wait(10)
        input.Take()

	sample_length := input.Samples.GetLength()
	assert.Equal(t, sample_length, 1)

	info_length := input.Infos.GetLength()
	assert.Equal(t, info_length, 1)

	valid := input.Infos.IsValid(0)
	assert.Equal(t, valid, true)

        color := input.Samples.GetString(0, "color")
	assert.Equal(t, color, "BLUE")
        x := input.Samples.GetInt(0, "x")
	assert.Equal(t, x, 1)
        y := input.Samples.GetInt(0, "y")
	assert.Equal(t, y, 1)
        z := input.Samples.GetBoolean(0, "y")
	assert.Equal(t, z, true)
        shapesize := input.Samples.GetInt(0, "shapesize")
	assert.Equal(t, shapesize, 5)
}
