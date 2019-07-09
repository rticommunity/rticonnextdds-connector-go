/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package main

import (
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
	"time"
)

const xmlString = `str://"<dds><qos_library name="QosLibrary"><qos_profile name="def" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true"/></qos_library>
<types>
<struct name="ShapeType" extensibility="extensible">
<member name="color" stringMaxLength="128" id="0" type="string" key="true"/>
<member name="x" id="1" type="long"/>
<member name="y" id="2" type="long"/>
<member name="shapesize" id="3" type="long"/>
</struct>
</types>
<domain_library name="MyDomainLibrary">
<domain name="MyDomain" domain_id="0"><register_type name="ShapeType" type_ref="ShapeType"/><topic name="Square" register_type_ref="ShapeType"/></domain>
</domain_library>
<domain_participant_library name="MyParticipantLibrary">
<domain_participant name="Zero" domain_ref="MyDomainLibrary::MyDomain">
<publisher name="MyPublisher"><data_writer name="MySquareWriter" topic_ref="Square" /></publisher>
</domain_participant></domain_participant_library></dds>"
`

func main() {
	// Create a connector defined in the XML configuration
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", xmlString)
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer connector.Delete()

	// Get an output from the connector
	output, err := connector.GetOutput("MyPublisher::MySquareWriter")
	if err != nil {
		log.Panic(err)
	}

	// Set values to the instance and write the instance
	for i := 0; i < 10; i++ {
		output.Instance.SetInt("x", i)
		output.Instance.SetInt("y", i*2)
		output.Instance.SetInt("shapesize", 30)
		output.Instance.SetString("color", "BLUE")
		output.Write()
		log.Println("Writing...")
		time.Sleep(time.Second * 1)
	}
}
