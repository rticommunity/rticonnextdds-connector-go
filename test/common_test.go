package connector_test

import (
        "github.com/rticommunity/rticonnextdds-connector-go"
        "runtime"
        "path"
)


func newTestConnector()(connector *rti.Connector){
        _, cur_path, _, _ := runtime.Caller(0)
        xml_path := path.Join(path.Dir(cur_path), "./xml/ShapeExample.xml")
        participant_profile := "MyParticipantLibrary::Zero"
        connector, _ = rti.NewConnector(participant_profile, xml_path)
        return connector
}

func deleteTestConnector(connector *rti.Connector){
	connector.Delete()
}

