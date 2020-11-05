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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rti "github.com/rticommunity/rticonnextdds-connector-go"
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
<subscriber name="MySubscriber"><data_reader name="MySquareReader" topic_ref="Square" /></subscriber>
</domain_participant></domain_participant_library></dds>"
`

func main() {
	// Create a channel to receive signals from OS
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Create a connector defined in the XML configuration
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", xmlString)
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer connector.Delete()

	// Get an input from the connector
	input, err := connector.GetInput("MySubscriber::MySquareReader")
	if err != nil {
		log.Panic(err)
	}

	run := true

	// Get values from a received sample and print them
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Received signal %v: terminating\n", sig)
			run = false
		default:
			input.Take()
			numOfSamples, _ := input.Samples.GetLength()
			for j := 0; j < numOfSamples; j++ {
				valid, _ := input.Infos.IsValid(j)
				if valid {
					color, _ := input.Samples.GetString(j, "color")
					x, _ := input.Samples.GetInt(j, "x")
					y, _ := input.Samples.GetInt(j, "y")
					shapesize, _ := input.Samples.GetInt(j, "shapesize")

					log.Println("---Received Sample---")
					log.Printf("color: %s\n", color)
					log.Printf("x: %d\n", x)
					log.Printf("y: %d\n", y)
					log.Printf("shapesize: %d\n", shapesize)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}
}
