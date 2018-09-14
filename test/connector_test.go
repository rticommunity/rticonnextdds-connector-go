package rti_test

import (
	"github.com/rticommunity/rticonnextdds-connector-go"
	"testing"
	"github.com/stretchr/testify/assert"
	"runtime"
	"path"
)

func TestInvalidXMLPath(t *testing.T) {
	participant_profile := "MyParticipantLibrary::Zero"
	invalid_xml_path := "invalid/path/to/xml"

	connector, err := rti.NewConnector(participant_profile, invalid_xml_path)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

func TestInvalidParticipantProfile(t *testing.T) {
        _, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./xml/ShapeExample.xml")
	invalid_participant_profile := "InvalidParticipantProfile"

        connector, err := rti.NewConnector(invalid_participant_profile, xml_path)
        assert.Nil(t, connector)
        assert.NotNil(t, err)
}

func TestMultipleConnectorCreation(t *testing.T) {
	_, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./xml/ShapeExample.xml")
	participant_profile := "MyParticipantLibrary::Zero"
	var connectors [5]*rti.Connector
	for i := 0; i < 5; i++ {
		connectors[i], _ = rti.NewConnector(participant_profile, xml_path)
		assert.NotNil(t, connectors[i])
	}
}
