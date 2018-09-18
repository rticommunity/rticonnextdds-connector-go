package rti

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"runtime"
	"path"
)

// Helper functions
func newTestConnector()(connector *Connector){
        _, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./test/xml/ShapeExample.xml")
        participant_profile := "MyParticipantLibrary::Zero"
        connector, _ = NewConnector(participant_profile, xml_path)
        return connector
}

func deleteTestConnector(connector *Connector){
        connector.Delete()
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
}
